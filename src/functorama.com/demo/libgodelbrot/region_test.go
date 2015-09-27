package libgodelbrot

import (
    "testing"
)

type regionSameness struct {
    a complex128
    b complex128
    n int
    same bool
}

func sameRegion(a Region, b Region) regionSameness {
    aPoints := a.Points()
    bPoints := b.Points()

    for i, ap := range aPoints {
        bp := bPoints[i]
        if (ap.c != bp.c) {
            return regionSameness{
                a: ap.c,
                b: bp.c,
                n: i,
                same: false,
            }
        }
    }
    return regionSameness{ same: true }
}

func TestRegionSplitPos(t *testing.T) {
    left := 1.0
    right := 3.0
    top := 3.0
    bottom := 1.0

    topLeft := complex(left, top)
    bottomRight := complex(right, bottom)
    midPoint := complex(2.0, 2.0)

    leftSideMid := complex(1.0, 2.0)
    rightSideMid := complex(3.0, 2.0)
    topSideMid := complex(2.0, 3.0)
    bottomSideMid := complex(2.0, 1.0)

    subjectRegion := NewRegion(topLeft, bottomRight)

    expected := []Region{
        NewRegion(topLeft, midPoint),
        NewRegion(topSideMid, rightSideMid),
        NewRegion(leftSideMid, bottomSideMid),
        NewRegion(midPoint, bottomRight),
    }

    actual := subjectRegion.Split()

    for i, ex := range expected {
        similarity := sameRegion(ex, actual.children[i])
        if (!similarity.same) {
            t.Error(
                "Unexpected child region ", i, 
                ", expected point ", similarity.n, 
                "to be ", similarity.a, 
                " but was ", similarity.b,
            )
        }
    }
}

func TestRegionSplitNeg(t *testing.T) {
    left := -100.0
    right := -24.0
    top := -10.0
    bottom := -340.0

    topLeft := complex(left, top)
    bottomRight := complex(right, bottom)
    midPoint := complex(-62.0, -175.0)

    leftSideMid := complex(-100.0, -175.0)
    rightSideMid := complex(-24.0, -175.0)
    topSideMid := complex(-62.0, -10.0)
    bottomSideMid := complex(-62.0, -340.0)

    subjectRegion := NewRegion(topLeft, bottomRight)

    expected := []Region{
        NewRegion(topLeft, midPoint),
        NewRegion(topSideMid, rightSideMid),
        NewRegion(leftSideMid, bottomSideMid),
        NewRegion(midPoint, bottomRight),
    }

    actual := subjectRegion.Split()

    for i, ex := range expected {
        similarity := sameRegion(ex, actual.children[i])
        if (!similarity.same) {
            t.Error(
                "Unexpected child region ", i, 
                ", expected point ", similarity.n, 
                "to be ", similarity.a, 
                " but was ", similarity.b,
            )
        }
    }
}

func TestRegionSplitPosAndNeg(t *testing.T) {
    left := -100.0
    right := 24.0
    top := 10.0
    bottom := -340.0

    topLeft := complex(left, top)
    bottomRight := complex(right, bottom)
    midPoint := complex(-38.0, -165.0)

    leftSideMid := complex(-100.0, -165.0)
    rightSideMid := complex(24.0, -165.0)
    topSideMid := complex(-38.0, 10.0)
    bottomSideMid := complex(-38.0, -340.0)

    subjectRegion := NewRegion(topLeft, bottomRight)

    expected := []Region{
        NewRegion(topLeft, midPoint),
        NewRegion(topSideMid, rightSideMid),
        NewRegion(leftSideMid, bottomSideMid),
        NewRegion(midPoint, bottomRight),
    }

    actual := subjectRegion.Split()

    for i, ex := range expected {
        similarity := sameRegion(ex, actual.children[i])
        if (!similarity.same) {
            t.Error(
                "Unexpected child region ", i, 
                ", expected point ", similarity.n, 
                "to be ", similarity.a, 
                " but was ", similarity.b,
            )
        }
    }
}

