package libgodelbrot

import (
    "math/big"
    "testing"
    "github.com/johnny-morrice/godelbrot/internal/bigbase"
)

func TestZoom(t *testing.T) {
    target := ZoomTarget{}
    target.Xmin = 10
    target.Xmax = 30
    target.Ymax = 20
    target.Ymin = 50

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