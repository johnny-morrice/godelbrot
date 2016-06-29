package godelbrot

import (
    "math/big"
    "testing"
    "github.com/johnny-morrice/godelbrot/internal/bigbase"
)

func TestMagnifyHalf(t *testing.T) {
        target := ZoomTarget{}
    target.Xmin = 10
    target.Xmax = 30
    target.Ymin = 20
    target.Ymax = 50

    req := DefaultRequest()
    req.RealMin = "0.5"
    req.RealMax = "0.6"
    req.ImagMin = "0.3"
    req.ImagMax = "0.4"
    req.ImageWidth = 100
    req.ImageHeight = 100
    req.Precision = 53

    expect := []*big.Float{
        big.NewFloat(0.505),
        big.NewFloat(0.565),
        big.NewFloat(0.325),
        big.NewFloat(0.39),
    }

    z := Zoom{ZoomTarget: target,}
    prev, conferr := Configure(req)
    if conferr != nil {
        t.Fatal(conferr)
    }
    z.Prev = *prev

    mag, magerr := z.Magnify(0.5)
    if magerr != nil {
        t.Fatal(magerr)
    }
    actual := []*big.Float{
        &mag.RealMin,
        &mag.RealMax,
        &mag.ImagMin,
        &mag.ImagMax,
    }
    for i, ex := range expect {
        ac := actual[i]
        margin := big.NewFloat(0.0)
        margin.Sub(ex, ac)
        margin.Abs(margin)
        if margin.Cmp(big.NewFloat(0.001)) > 0  {
            t.Error("Fail at", i,
                "expected", bigbase.DbgF(*ex), "but received", bigbase.DbgF(*ac))
        }
    }
}

func TestMagnifyWhole(t *testing.T) {
    target := ZoomTarget{}
    target.Xmin = 10
    target.Xmax = 30
    target.Ymin = 20
    target.Ymax = 50

    req := DefaultRequest()
    req.RealMin = "0.5"
    req.RealMax = "0.6"
    req.ImagMin = "0.3"
    req.ImagMax = "0.4"
    req.ImageWidth = 100
    req.ImageHeight = 100
    req.Precision = 53

    expect := []*big.Float{
        big.NewFloat(0.51),
        big.NewFloat(0.53),
        big.NewFloat(0.35),
        big.NewFloat(0.38),
    }

    z := Zoom{ZoomTarget: target,}
    prev, conferr := Configure(req)
    if conferr != nil {
        t.Fatal(conferr)
    }
    z.Prev = *prev

    mag, magerr := z.Magnify(1.0)
    if magerr != nil {
        t.Fatal(magerr)
    }
    actual := []*big.Float{
        &mag.RealMin,
        &mag.RealMax,
        &mag.ImagMin,
        &mag.ImagMax,
    }
    for i, ex := range expect {
        ac := actual[i]
        margin := big.NewFloat(0.0)
        margin.Sub(ex, ac)
        margin.Abs(margin)
        if margin.Cmp(big.NewFloat(0.001)) > 0  {
            t.Error("Fail at", i,
                "expected", bigbase.DbgF(*ex), "but received", bigbase.DbgF(*ac))
        }
    }
}

func TestMovie(t *testing.T) {
    const framecnt = 2

    target := ZoomTarget{}
    target.Xmin = 10
    target.Xmax = 30
    target.Ymin = 20
    target.Ymax = 50
    target.Frames = framecnt

    req := DefaultRequest()
    req.RealMin = "0.5"
    req.RealMax = "0.6"
    req.ImagMin = "0.3"
    req.ImagMax = "0.4"
    req.ImageWidth = 100
    req.ImageHeight = 100
    req.Precision = 53

    expectations := make([][]*big.Float, framecnt)
    expectations[0] = []*big.Float{
        big.NewFloat(0.505),
        big.NewFloat(0.565),
        big.NewFloat(0.325),
        big.NewFloat(0.39),
    }
    expectations[1] = []*big.Float{
        big.NewFloat(0.51),
        big.NewFloat(0.53),
        big.NewFloat(0.35),
        big.NewFloat(0.38),
    }

    z := Zoom{ZoomTarget: target,}
    prev, conferr := Configure(req)
    if conferr != nil {
        t.Fatal(conferr)
    }
    z.Prev = *prev

    frames, magerr := z.Movie()
    if magerr != nil {
        t.Fatal(magerr)
    }
    actualities := make([][]*big.Float, framecnt)
    for i, fr := range frames {
        fract := []*big.Float{
            &fr.RealMin,
            &fr.RealMax,
            &fr.ImagMin,
            &fr.ImagMax,
        }
        actualities[i] = fract
    }
    for i, expectSet := range expectations {
        actualSet := actualities[i]
        for j, ex := range expectSet {
            ac := actualSet[j]
            margin := big.NewFloat(0.0)
            margin.Sub(ex, ac)
            margin.Abs(margin)
            if margin.Cmp(big.NewFloat(0.001)) > 0  {
                t.Error("Fail at frame", i, "bound", j,
                    "expected", bigbase.DbgF(*ex), "but received", bigbase.DbgF(*ac))
            }
        }
    }
}