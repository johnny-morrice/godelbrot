package libgodelbrot

import (
    "testing"
)

func TestConfigure(t *testing.T) {
    config := easyConfig()
    expectedHUnit := 0.02
    expectedVUnit := 0.02
    if config.HorizUnit != expectedHUnit {
        t.Error("Expected ", expectedHUnit, " but received ", config.HorizUnit)
    }
    if config.VerticalUnit != expectedVUnit {
        t.Error("Expected ", expectedVUnit, " but received ", config.VerticalUnit)
    }
}

func TestPlaneBottomRight(t *testing.T) {
    config := easyConfig()
    expected := 1 - 1i
    actual := config.PlaneBottomRight()
    if expected != actual {
        t.Error("Expected ", expected, " but was ", actual)
    }
}

func TestPlaneToPixel(t *testing.T) {
    config := easyConfig()

    qA := complex(0.1, 0.1)
    qB := complex(0.1, -0.1)
    qC := complex(-0.1, -0.1)
    qD := complex(-0.1, 0.1)
    origin := complex(0.0, 0.0)
    offset := complex(-1.0, 1.0)

    var expectPixAx uint = 55
    var expectPixAY uint = 45

    var expectPixBx uint = 55
    var expectPixBy uint = 55

    var expectPixCx uint = 45
    var expectPixCy uint = 55

    var expectPixDx uint = 45
    var expectPixDy uint = 45

    var expectOx uint = 50
    var expectOy uint = 50

    var expectOffsetX uint = 0
    var expectOffsetY uint = 0

    points := []complex128{qA, qB, qC, qD, origin, offset}
    expectedXs := []uint{
        expectPixAx, 
        expectPixBx, 
        expectPixCx, 
        expectPixDx, 
        expectOx, 
        expectOffsetX,
    }
    expectedYs := []uint{
        expectPixAY, 
        expectPixBy, 
        expectPixCy, 
        expectPixDy, 
        expectOy, 
        expectOffsetY,
    }

    for i, point := range points {
        expectedX := expectedXs[i]
        expectedY := expectedYs[i]
        actualX, actualY := config.PlaneToPixel(point)
        if actualX != expectedX || actualY != expectedY {
            t.Error("Error on point", i, ":", point, 
                " expected (", expectedX, ",", expectedY, ") but was",
                "(", actualX, ",", actualY, ")")
        }
    }
}

func easyConfig() *RenderConfig {
    params := RenderParameters{
        IterateLimit: 100,
        DivergeLimit: 1.0,
        Width: 100,
        Height: 100,
        TopLeft: complex(-1.0, 1.0),
        BottomRight: complex(1.0, -1.0),
        Frame: CornerFrame,
        RegionCollapse: 10,
    }
    return params.Configure()
}