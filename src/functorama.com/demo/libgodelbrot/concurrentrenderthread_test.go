package libgodelbrot

import (
    "testing"
)

type mockSharedRegionNumerics struct {
    mockRegionNumerics

    tGrabThreadPrototype bool
}

func (mock mockSharedRegionNumerics) GrabThreadPrototype(threadId uint) {
    mock.tGrabThreadPrototype = true
}

type threadOutputExpect struct {
    members int
    children int
    uniform int
}

func TestRenderThreadFactory(t *testing.T) {
    mock := mockRenderApplication{}
    factory := newRenderThreadFactory(mock)

    threads := []renderThreadfactory{factory(nil, nil), factory(nil, nil)}

    if !mock.tConcurrentConfig {
        t.Error("Mock did not receive expected method call")
    }

    for i, th := range threads {
        if th.threadId != i {
            t.Error("Thread", i, "had incorrect threadId: ", th)
        }
    }
}

func TestThreadRun(t *testing.T) {
    th := createThread()
    commandCount := 2
    th.inputChan := make(chan renderInput, commandCount)
    th.outputChan := make(chan renderOutput)
    th.inputChan <- renderInput{
        command: render,
        regions: uniformer(),
    },
    th.inputChan <- renderInput{command: stop}

    th.run()
    out := th.outputChan

    const context = "TestThreadRun"
    threadOutputCheck(t, out, threadOutputExpect{1, 0, 0}, context)
}

func TestThreadPass(t *testing.T) {
    // Key: Uniformer Count, Subdivider Count, Member Count
    
    // 0 0 0
    zero := threadPassOutput([]mockSharedRegionNumerics{})

    // 1 0 0
    oneUniform := threadPassOutput([]mockSharedRegionNumerics{uniformer()})
    // 2 0 0
    twoUniform := threadPassOutput([]mockSharedRegionNumerics{uniformer(), uniformer()})

    // 0 1C 0
    oneChild := threadPassOutput([]mockSharedRegionNumerics{subdivider()})
    // 0 2C 0
    twoChild := threadPassOutput([]mockSharedRegionNumerics{subdivider(), subdivider()})

    // 0 0 1C.M
    oneMember := threadPassOutput([]mockSharedRegionNumerics{collapser())
    // 0 0 2C.M.
    twoMember := threadPassOutput([]mockSharedRegionNumerics{collapser(), collapser()})

    // 1 1C 0
    oneUniOneChild := threadPassOutput([]mockSharedRegionNumerics{uniformer(), subdivider()})
    // 1 0 1C.M.
    oneUniOneMember := threadPassOutput([]mockSharedRegionNumerics{uniformer(), subdivider()})
    // 0 1C 1C.M
    oneChildOneMember := threadPassOutput([]mockSharedRegionNumerics{subdivider(), collapser()})

    // 1 1C 1C.M
    all := threadPassOutput([]mockSharedRegionNumerics{uniformer(), subdivider(), collapser()})

    const context = "TestThreadPass"

    threadOutputCheck(t, zero, threadOutputExpect{}, context)

    threadOutputCheck(t, oneUniform, threadOutputExpect{1, 0, 0}, context)
    threadOutputCheck(t, twoUniform, threadOutputExpect{2, 0, 0}, context)

    threadOutputCheck(t, oneChild, threadOutputExpect{0, children, 0}, context)
    threadOutputCheck(t, twoChild, threadOutputExpect{0, 2 * children, 0}, context)

    threadOutputCheck(t, oneMember, threadOutputExpect{collapseMembers, 0, 0}, context)
    threadOutputCheck(t, twoMember, threadOutputExpect{2 * collapseMembers, 0, 0}, context)

    threadOutputCheck(t, oneUniOneChild, threadOutputExpect{1, children, 0}, context)
    threadOutputCheck(t, oneUniOneMember, threadOutputExpect{1, 0, collapseMembers}, context)
    threadOutputCheck(t, oneChildOneMember, threadOutputExpect{0, children, collapseMembers}, context)

    threadOutputCheck(t, all, threadOutputExpect{1, children, collapseMembers}, context)
}

func TestThreadStep(t *testing.T) {
    coll := collapser()
    subd := subdivider()
    uni := uniformer()
    collapsed := threadStepOutput(coll)
    subdivided := threadStepOutput(subd)
    uniformed := threadStepOutput(uni)

    for _, mock := range []mockSharedRegionNumerics{coll, subd, uni} {
        if !mockStepOkay(mock) {
            t.Error("General case methods not called on region:", mock)
        }
    }

    if !(collapsed.tRegionalSequenceNumerics && collapsed.tMandelbrotMembers) {
        t.Error("Expected methods not called on collapse region:", collapsed)
    }

    if !(subdivided.tSplit && subdivided.tChildren && subdivided.tEvaluateAllPoints) {
        t.Error("Expected methods not called on subdivided region:", subdivided)
    }

    if !uniformed.tEvaluateAllPoints {
        t.Error("Expected methods not called on uniform region:", uniformed)
    }

    const context = "TestThreadStep"
    threadOutputCheck(t, collapsed, threadOutputExpect{0, 0, collapseMembers}, context)
    threadOutputCheck(t, subdivided, threadOutputExpect{0, children, 0}, context)
    threadOutputCheck(t, uniformed, threadOutputExpect{1, 0, 0}, context)
}

func collapser(length int) mockSharedRegionNumerics {
    return mockSharedRegionNumerics{
        path: collapse,
        length: length,
    }
}

func mockStepOkayGeneral(mock mockSharedRegionNumerics) bool {
    okay := mock.tCollapseSize
    okay := okay && mock.Rect
    okay := okay && mock.tGrabThreadPrototype
    okay := okay && mock.tClaimExtrinsics
    okay := okay && mock.tUniform
    return okay
}

func subdivider(length int) mockSharedRegionNumerics {
    return mockSharedRegionNumerics{
        path: subdivide,
        length: length,
    }
}

func uniformer(length int) mockSharedRegionNumerics {
    return mockSharedRegionNumerics{
        path: uniform,
        length: length,
    }
}

func threadOutputCheck(t *testing.T, expect threadOutputExpect, context string) {
    actualUniformCount := len(out.uniform)
    actualChildCount := len(out.children)
    actualMemberCount := len(out.members)

    if actualMemberCount != expect.members || actualChildCount != expect.children || actualUniformCount != expect.uniform {
        t.Error("In context ", context, 
            ", expected output counts ", expect, " but received (", 
            actualUniformCount, actualChildCount, actualMemberCount, ")")
    }
}

func threadPassOutput(helpers []mockSharedRegionNumerics) renderOutput {
    numerics := make([]SharedRegionNumerics, len(helpers))
    for i, h := range helpers {
        numerics[i] = mockSharedRegionNumerics
    }

    th := createThread()
    return th.pass(numerics)
}

func threadStepOutput(helper mockSharedRegionNumerics) renderOutput {
    output := renderOutput{}
    input := renderInput{
        command: render,
    }

    input.regions = []SharedRegionNumerics{mockSharedRegionNumerics}
    th := createThread()
    th.Step(input, &output)

    return output
}

func createThread() renderThread {
    return renderThread{
        config: renderParameters{
            RegionCollapseSize: collapseSize,
        },
    }
}