package libgodelbrot

import (
    "encoding/json"
    "errors"
    "io"
    "math/big"
    "github.com/johnny-morrice/godelbrot/internal/bigbase"
)

type ZoomTarget struct {
    Xmin uint
    Xmax uint
    Ymin uint
    Ymax uint
    // Reconsider numerical system and render modes as appropriate.
    Reconfigure bool
    // Increase precision.  With Reconfigure, this should automatically engage arbitrary
    // precision mode.
    UpPrec bool
    // Number of frames for zoom
    Frames uint
}

func (zt *ZoomTarget) Validate() error {
    if zt.Xmin >= zt.Xmax || zt.Ymin >= zt.Ymax {
        return errors.New("Min and max zoom boundaries are invalid.")
    }
    return nil
}

// Zoom into a portion of the previous image.
type Zoom struct {
    Prev Info
    ZoomTarget
}

type distort struct {
    prev big.Float
    next big.Float
}

func (d distort) para(time *big.Float) *big.Float {
    prec := d.next.Prec()
    delta := bigbase.MakeBigFloat(0.0, prec)
    delta.Sub(&d.prev, &d.next)
    delta.Abs(&delta)
    delta.Mul(&delta, time)
    if d.prev.Cmp(&d.next) > 0 {
        return delta.Sub(&d.prev, &delta)
    } else {
        return delta.Add(&d.prev, &delta)
    }
}

// Frame zooms towards the target coordinates.  Degree = 1 is a complete zoom.
func (z *Zoom) rescope(degree float64) (*Info, error) {
    info := z.lens(degree)

    if z.Reconfigure {
        return Configure(&info.UserRequest)
    } else {
        return info, nil
    }
}

func (z *Zoom) lens(degree float64) *Info {
    baseapp := makeBaseFacade(&z.Prev)
    app := makeBigBaseFacade(&z.Prev, baseapp)
    num := bigbase.Make(app)

    time := bigbase.MakeBigFloat(degree, num.Precision)

    min := num.PixelToPlane(int(z.Xmin), int(z.Ymin))
    max := num.PixelToPlane(int(z.Xmax), int(z.Ymax))

    target := []big.Float{
        min.R,
        min.I,
        max.R,
        max.I,
    }

    bounds := []big.Float{
        z.Prev.RealMin,
        z.Prev.ImagMin,
        z.Prev.RealMax,
        z.Prev.ImagMax,
    }
    zoom := make([]*big.Float, len(bounds))
    for i, b := range bounds {
        d := distort{
            prev: b,
            next: target[i],
        }

        zoom[i] = d.para(&time)
    }

    info := new(Info)
    *info = z.Prev
    info.RealMin = *zoom[0]
    info.ImagMin = *zoom[1]
    info.RealMax = *zoom[2]
    info.ImagMax = *zoom[3]
    info.UserRequest = info.GenRequest()

    return info
}

func (z *Zoom) Magnify(degree float64) (*Info, error) {
    info := z.lens(degree)

    if z.UpPrec {
        for !info.IsAccurate() {
            z.Prev.AddPrec(1)
            z.Prev.UserRequest = z.Prev.GenRequest()
            info = z.lens(degree)
        }
    }

    if z.Reconfigure {
        return Configure(&info.UserRequest)
    } else {
        return info, nil
    }
}


// Movie is a parametric expansion of frames.
func (z *Zoom) Movie() ([]*Info, error) {
    cnt := z.Frames
    if cnt == 0 {
        return []*Info{}, nil
    }

    // Compute last frame first to encourage numerical stability
    const fullZoom = 1.0
    last, lerr := z.Magnify(fullZoom)
    if lerr != nil {
        return nil, lerr
    }

    // Compute intervening frames
    interval := 1.0 / float64(cnt)
    time := 0.0
    frames := make([]*Info, cnt)
    for i := uint(0); i < cnt - 1; i++ {
        time += interval
        info, err := z.rescope(time)
        if err != nil {
            return nil, err
        }
        frames[i] = info
    }

    frames[cnt - 1] = last
    return frames, nil
}

type UserZoom struct {
    Prev UserInfo
    ZoomTarget
}

func ReadZoom(r io.Reader) (*Zoom, error) {
    uz := &UserZoom{}
    dec := json.NewDecoder(r)
    decerr := dec.Decode(uz)
    if decerr != nil {
        return nil, decerr
    }


    info, inferr := Unfriendly(&uz.Prev)
    if inferr != nil {
        return nil, inferr
    }

    z := &Zoom{}
    z.Prev = *info
    z.ZoomTarget = uz.ZoomTarget
    return z, nil
}

func WriteZoom(w io.Writer, z *Zoom) error {
    uz := &UserZoom{}
    uz.Prev = *Friendly(&z.Prev)
    uz.ZoomTarget = z.ZoomTarget

    enc := json.NewEncoder(w)
    err := enc.Encode(uz)
    if err != nil {
        return err
    }
    return nil
}