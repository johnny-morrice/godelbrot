package sharedregion

import (
	"testing"
	"functorama.com/demo/base"
	"functorama.com/demo/region"
)

const children = 4
const collapseCount = 20
const collapseSize = 10

type threadOutputExpect struct {
    uniform  int
    children int
    members  int
}

func TestRenderThreadFactory(t *testing.T) {
	mock := &MockRenderApplication{}
	factory := NewRenderThreadFactory(mock)

	threads := []RenderThread{factory.Build(nil, nil), factory.Build(nil, nil)}

	if !mock.TSharedRegionConfig {
		t.Error("Mock did not receive expected method call")
	}

	for i, th := range threads {
		if th.ThreadId != uint(i) {
			t.Error("Thread", i, "had incorrect ThreadId: ", th)
		}
	}
}

func TestThreadRun(t *testing.T) {
	th := createThread()
	const commandCount = 2
	const uniformLength = 1
	iChan := make(chan RenderInput, commandCount)
	oChan := make(chan RenderOutput)
	th.InputChan = iChan
	th.OutputChan = oChan

	iChan <- RenderInput{
		Command: ThreadRender,
		Regions: []SharedRegionNumerics{uniformer()},
	}
	iChan <- RenderInput{Command: ThreadStop}

	go th.Run()
	out := <- oChan

	const context = "TestThreadRun"
	threadOutputCheck(t, out, threadOutputExpect{1, 0, 0}, context)
}

func TestThreadPass(t *testing.T) {
	// Key: Uniformer Count, Subdivider Count, Member Count

	// 0 0 0
	zero := threadPassOutput([]SharedRegionNumerics{})

	// 1 0 0
	oneUniform := threadPassOutput([]SharedRegionNumerics{uniformer()})
	// 2 0 0
	twoUniform := threadPassOutput([]SharedRegionNumerics{uniformer(), uniformer()})

	// 0 1C 0
	oneChild := threadPassOutput([]SharedRegionNumerics{subdivider()})
	// 0 2C 0
	twoChild := threadPassOutput([]SharedRegionNumerics{subdivider(), subdivider()})

	// 0 0 1C.M
	oneMember := threadPassOutput([]SharedRegionNumerics{collapser()})
	// 0 0 2C.M.
	twoMember := threadPassOutput([]SharedRegionNumerics{collapser(), collapser()})

	// 1 1C 0
	oneUniOneChild := threadPassOutput([]SharedRegionNumerics{uniformer(), subdivider()})
	// 1 0 1C.M.
	oneUniOneMember := threadPassOutput([]SharedRegionNumerics{uniformer(), collapser()})
	// 0 1C 1C.M
	oneChildOneMember := threadPassOutput([]SharedRegionNumerics{subdivider(), collapser()})

	// 1 1C 1C.M
	all := threadPassOutput([]SharedRegionNumerics{uniformer(), subdivider(), collapser()})

	const context = "TestThreadPass"

	threadOutputCheck(t, zero, threadOutputExpect{}, context)

	threadOutputCheck(t, oneUniform, threadOutputExpect{1, 0, 0}, context)
	threadOutputCheck(t, twoUniform, threadOutputExpect{2, 0, 0}, context)

	threadOutputCheck(t, oneChild, threadOutputExpect{0, children, 0}, context)
	threadOutputCheck(t, twoChild, threadOutputExpect{0, 2 * children, 0}, context)

	threadOutputCheck(t, oneMember, threadOutputExpect{0, 0, collapseCount}, context)
	threadOutputCheck(t, twoMember, threadOutputExpect{0, 0, 2 * collapseCount}, context)

	threadOutputCheck(t, oneUniOneChild, threadOutputExpect{1, children, 0}, context)
	threadOutputCheck(t, oneUniOneMember, threadOutputExpect{1, 0, collapseCount}, context)
	threadOutputCheck(t, oneChildOneMember, threadOutputExpect{0, children, collapseCount}, context)

	threadOutputCheck(t, all, threadOutputExpect{1, children, collapseCount}, context)
}

func TestThreadStep(t *testing.T) {
	coll := collapser()
	subd := subdivider()
	uni := uniformer()
	collapsed := threadStepOutput(coll)
	subdivided := threadStepOutput(subd)
	uniformed := threadStepOutput(uni)

	for i, mock := range []*MockNumerics{coll, subd, uni} {
		if !stepOkayGeneral(mock) {
			t.Error("General case methods not called on region", i, mock)
		}
	}

	if !(coll.TSharedRegionSequence) {
		t.Error("Expected methods not called on collapse region:", coll)
	}

	if !(subd.TSplit && subd.TSharedChildren && subd.TEvaluateAllPoints) {
		t.Error("Expected methods not called on subdivided region:", subd)
	}

	if !(uni.TOnGlitchCurve && uni.TEvaluateAllPoints)  {
		t.Error("Expected methods not called on uniform region:", uni)
	}

	const context = "TestThreadStep"
	threadOutputCheck(t, collapsed, threadOutputExpect{0, 0, collapseCount}, context)
	threadOutputCheck(t, subdivided, threadOutputExpect{0, children, 0}, context)
	threadOutputCheck(t, uniformed, threadOutputExpect{1, 0, 0}, context)
}

func stepOkayGeneral(mock *MockNumerics) bool {
	okay := mock.TRect && mock.TGrabThreadPrototype
	okay = okay && mock.TClaimExtrinsics
	return okay
}

func collapser() *MockNumerics {
	mock := mocker(region.CollapsePath)
	mock.SharedMockSequence.Captured = make([]base.PixelMember, collapseCount)
	return mock
}

func subdivider() *MockNumerics {
	mock := mocker(region.SubdividePath)
	children := make([]*MockNumerics, children)

	for i := 0; i < len(children); i++ {
		children[i] = uniformer()
	}

	mock.SharedMockChildren = children
	return mock
}

func uniformer() *MockNumerics {
	return mocker(region.UniformPath)
}

func mocker(path region.RegionType) *MockNumerics {
	mock := &MockNumerics{}
	mock.Path = path
	mockSequence := &MockSequence{}
	mock.SharedMockSequence = mockSequence
	mock.AppCollapseSize = collapseSize
	return mock
}

func threadOutputCheck(t *testing.T, actual RenderOutput, expect threadOutputExpect, context string) {
	actualUniformCount := len(actual.UniformRegions)
	actualChildCount := len(actual.Children)
	actualMemberCount := len(actual.Members)

	okay := actualMemberCount == expect.members
	okay = okay && actualChildCount == expect.children
	okay = okay && actualUniformCount == expect.uniform
	if !okay {
		t.Error("In context", context,
			"expected output counts", expect, "but received (",
			actualUniformCount, actualChildCount, actualMemberCount, ")")
	}
}

func threadPassOutput(numerics []SharedRegionNumerics) RenderOutput {
	th := createThread()
	return th.Pass(numerics)
}

func threadStepOutput(numerics SharedRegionNumerics) RenderOutput {
	output := RenderOutput{}
	th := createThread()
	th.Step(numerics, &output)

	return output
}

func createThread() RenderThread {
	return RenderThread{
		RegionConfig: region.RegionConfig{
			CollapseSize: collapseSize,
		},
	}
}
