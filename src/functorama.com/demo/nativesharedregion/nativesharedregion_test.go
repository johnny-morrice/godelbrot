package nativesharedregion

import (
	"testing"
	"functorama.com/demo/nativeregion"
)

func TestCreateNativeSharedRegion(t *testing.T) {
	const jobCount = 2

	// Pointer to non-zero region 
	region := &nativeregion.NativeRegionNumerics{}
	region.SqrtDivergeLimit = 3.0

	shared := CreateNativeSharedRegion(numerics, jobCount)

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

func TestRegionGrabThreadPrototype(t *testing.T) {
	const jobCount
	region := &nativeregion.NativeRegionNumerics{}
	shared := CreateNativeSharedRegion(region, jobCount)

	const successId = 1
	shared.GrabThreadPrototype(successId)
	if shared.threadId != successId {
		t.Error("Expected threadId", successId,
			"but received", shared.threadId)
	}

	const badId = -1
	expectPanic(t, func() { shared.GrabThreadPrototype(badId) })
}

func TestSharedChildren(t *testing.T) {
	region := &nativeregion.NativeRegionNumerics{
	}
}

func expectPanic(t *testing.T, broken func()) {
	smoothSailing := true

	inner := func() {
		defer func() {
			failure := recover()
			smoothSailing = failure == nil
		}
		broken()
	}

	inner()

	if smoothSailing {
		t.Error("Expected panic")
	}
}