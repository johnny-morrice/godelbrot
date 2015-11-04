package nativesharedregion

import (
	"sync"
	"testing"
	"functorama.com/demo/sharedregion"
	"functorama.com/demo/nativesequence"
	"functorama.com/demo/nativeregion"
)

func TestCreateNativeSharedRegion(t *testing.T) {
	const jobCount = 2

	// Pointer to non-zero region
	region := createRegion()
	region.SqrtDivergeLimit = 3.0

	shared := CreateNativeSharedRegion(region, jobCount)

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

func TestRegionGrabThreadPrototypeEdge(t *testing.T) {
	const jobCount = 1
	region := createRegion()
	shared := CreateNativeSharedRegion(region, jobCount)

	testMutantEdge(t, shared)
}


func TestRegionGrabThreadPrototypeParallel(t *testing.T) {
	// We are testing that one thread can mutate its own state without touching that of others
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	const jobCount = 3

	// Pointer to non-zero region
	region := createRegion()

	shared := CreateNativeSharedRegion(region, jobCount)

	testMutantParallel(t, shared, jobCount)
}



func TestSharedChildren(t *testing.T) {
	const jobCount = 1
	const expectCount = 4

	region := createRegion()
	shared := CreateNativeSharedRegion(region, jobCount)

	shared.Split()

	children := shared.SharedChildren()
	actualCount := len(children)

	if actualCount != expectCount {
		t.Error("Expected", expectCount, "children",
			"but received", actualCount)
	}
}

func TestSequenceGrabThreadPrototypeParallel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	const jobCount = 3
	// Pointer to non-zero region
	region := createRegion()

	shared := CreateNativeSharedRegion(region, jobCount)
	sequence := shared.NativeSharedSequence()

	testMutantParallel(t, sequence, jobCount)
}

func TestSequenceGrabThreadPrototypeEdge(t *testing.T) {
	const jobCount = 1
	region := createRegion()
	shared := CreateNativeSharedRegion(region, jobCount)
	sequence := shared.NativeSharedSequence()

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
	sharedregion.OpaqueThreadPrototype
	mutate()
	isMutant() bool
	id() uint
}

func testMutantParallel(t *testing.T, shared mutator, jobCount uint) {
	const runs = 10000
	const mutantId = 0

	hold := sync.WaitGroup{}

	checkDiverge := func(threadId uint, local mutator) {
		hold.Add(1)
		defer hold.Done()
		for i := 0; i < runs; i++ {
			local.GrabThreadPrototype(threadId)
			if threadId == mutantId {
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

	for threadId := uint(0); threadId < jobCount; threadId++ {
		go checkDiverge(threadId, shared)
	}

	hold.Wait()
}

func testMutantEdge(t *testing.T, shared mutator) {
	const successId uint = 0
	const badId uint = 1000

	shared.GrabThreadPrototype(successId)
	if shared.id() != successId {
		t.Error("Expected threadId", successId,
			"but received", shared.id())
	}

	expectPanic(t, func() { shared.GrabThreadPrototype(badId) })
}

const mutateDiverge = 4.0

func (region NativeSharedRegion) id() uint {
	return region.threadId
}

func (region NativeSharedRegion) mutate() {
	region.SqrtDivergeLimit = mutateDiverge
}

func (region NativeSharedRegion) isMutant() bool {
	return region.SqrtDivergeLimit == mutateDiverge
}

func (sequence NativeSharedSequence) id() uint {
	return sequence.threadId
}

func (sequence NativeSharedSequence) mutate() {
	sequence.SqrtDivergeLimit = mutateDiverge
}

func (sequence NativeSharedSequence) isMutant() bool {
	return sequence.SqrtDivergeLimit == mutateDiverge
}

func createRegion() *nativeregion.NativeRegionNumerics {
	return &nativeregion.NativeRegionNumerics{
		SequenceNumerics: &nativesequence.NativeSequenceNumerics{},
	}
}
