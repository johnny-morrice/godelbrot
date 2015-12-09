package bigsharedregion

import (
	"sync"
	"testing"
	"functorama.com/demo/bigbase"
	"functorama.com/demo/sharedregion"
)

const prec = 53

func TestMake(t *testing.T) {
	const jobCount = 2

	app := makeApp(jobCount)

	shared := Make(app)

	actualProtoCount := len(shared.prototypes)
	actualSeqCount := len(shared.sequencePrototypes)

	if actualProtoCount != jobCount {
		t.Error("Expected", jobCount, "region prototypes",
			"but received", actualProtoCount)
	}

	if actualSeqCount != jobCount {
		t.Error("Expected", jobCount, "sequence prototypes",
			"but received", actualSeqCount)
	}
}

func TestRegionGrabWorkerPrototypeEdge(t *testing.T) {
	const jobCount = 1

	app := makeApp(jobCount)
	shared := Make(app)

	testMutantEdge(t, shared)
}


func TestRegionGrabWorkerPrototypeParallel(t *testing.T) {
	// We are testing that one thread can mutate its own state without touching that of others
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	const jobCount = 3

	app := makeApp(jobCount)
	shared := Make(app)

	testMutantParallel(t, shared, jobCount)
}

func TestSharedChildren(t *testing.T) {
	const jobCount = 1
	const expectCount = 4

	app := makeApp(jobCount)
	shared := Make(app)

	shared.Split()

	children := shared.SharedChildren()
	actualCount := len(children)

	if actualCount != expectCount {
		t.Error("Expected", expectCount, "children",
			"but received", actualCount)
	}
}

func TestSequenceGrabWorkerPrototypeParallel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	const jobCount = 3
	// Pointer to non-zero region
	app := makeApp(jobCount)

	shared := Make(app)
	sequence := shared.SharedSequence()

	testMutantParallel(t, sequence, jobCount)
}

func TestSequenceGrabWorkerPrototypeEdge(t *testing.T) {
	const jobCount = 1
	app := makeApp(jobCount)
	shared := Make(app)
	sequence := shared.SharedSequence()

	testMutantEdge(t, sequence)
}

func expectPanic(t *testing.T, broken func()) {
	smoothSailing := true

	inner := func() {
		defer func() {
			failure := recover()
			smoothSailing = failure == nil
		}()
		broken()
	}

	inner()

	if smoothSailing {
		t.Error("Expected panic")
	}
}

type mutator interface {
	sharedregion.OpaqueWorkerPrototype
	mutate()
	isMutant() bool
	id() uint16
}

func testMutantParallel(t *testing.T, shared mutator, jobCount uint16) {
	const runs = 10000
	const mutantId = 0

	hold := sync.WaitGroup{}

	checkDiverge := func(workerId uint16, local mutator) {
		hold.Add(1)
		defer hold.Done()
		for i := 0; i < runs; i++ {
			local.GrabWorkerPrototype(workerId)
			if workerId == mutantId {
				if i == 0 {
					local.mutate()
				} else {
					if !local.isMutant() {
						t.Fatal("Expected mutation")
					}
				}
			} else {
				if local.isMutant() {
					t.Fatal("Expected non-mutation")
				}
			}
		}
	}

	for workerId := uint16(0); workerId < jobCount; workerId++ {
		go checkDiverge(workerId, shared)
	}

	hold.Wait()
}

func testMutantEdge(t *testing.T, shared mutator) {
	const successId uint16 = 0
	const badId uint16 = 1000

	shared.GrabWorkerPrototype(successId)
	if shared.id() != successId {
		t.Error("Expected workerId", successId,
			"but received", shared.id())
	}

	expectPanic(t, func() { shared.GrabWorkerPrototype(badId) })
}

var mutateDiverge = bigbase.MakeBigFloat(4.0, prec)

func (region BigSharedRegion) id() uint16 {
	return region.workerId
}

func (region BigSharedRegion) mutate() {
	region.SqrtDivergeLimit = mutateDiverge
}

func (region BigSharedRegion) isMutant() bool {
	return bigbase.BigEq(&region.SqrtDivergeLimit, &mutateDiverge)
}

func (sequence BigSharedSequence) id() uint16 {
	return sequence.workerId
}

func (sequence BigSharedSequence) mutate() {
	sequence.SqrtDivergeLimit = mutateDiverge
}

func (sequence BigSharedSequence) isMutant() bool {
	return bigbase.BigEq(&sequence.SqrtDivergeLimit, &mutateDiverge)
}

func makeApp(jobCount uint16) *MockRenderApplication {
	app := &MockRenderApplication{}
	app.SharedConfig.Jobs = jobCount
	app.PictureWidth = 1
	app.PictureHeight = 1
	app.UserMin = bigbase.MakeBigComplex(1.0, 1.0, prec)
	app.UserMax = bigbase.MakeBigComplex(2.0, 2.0, prec)
	app.Prec = prec
	return app
}