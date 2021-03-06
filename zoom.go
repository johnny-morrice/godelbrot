package godelbrot

import (
	"encoding/json"
	"github.com/johnny-morrice/godelbrot/config"
	"github.com/johnny-morrice/godelbrot/internal/bigbase"
	"io"
	"log"
	"math/big"
)

type ZoomTarget struct {
	config.ZoomTarget
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

	delta.Sub(&d.next, &d.prev)
	delta.Mul(&delta, time)

	extra := bigbase.MakeBigFloat(0.0, prec)

	return extra.Add(&d.prev, &delta)
}

// Frame zooms towards the target coordinates.  Degree = 1 is a complete zoom.
func (z *Zoom) rescope(degree float64) (*Info, error) {
	info := z.lens(degree)

	if z.Reconfigure {
		log.Println("Reconfigure zoom Info")
		return Configure(&info.UserRequest)
	} else {
		return info, nil
	}
}

func (z *Zoom) lens(degree float64) *Info {
	appinfo := new(Info)
	*appinfo = z.Prev
	appinfo.UserRequest.FixAspect = config.Stretch
	baseapp := makeBaseFacade(appinfo)
	app := makeBigBaseFacade(appinfo, baseapp)
	num := bigbase.Make(app)

	time := bigbase.MakeBigFloat(degree, num.Precision)

	// Y min and max reversed as pixels grow downward...
	min := num.PixelToPlane(int(z.Xmin), int(z.Ymax))
	max := num.PixelToPlane(int(z.Xmax), int(z.Ymin))

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

		res := d.para(&time)
		zoom[i] = res
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
		dumpacc := false
		for !info.IsAccurate() {
			if __DEBUG {
				dumpacc = true
			}
			z.Prev.AddPrec(1)
			z.Prev.UserRequest = z.Prev.GenRequest()
			info = z.lens(degree)
		}

		if dumpacc {
			log.Printf("Accuracy increased to %v", z.Prev.RealMin.Prec())
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
	if __DEBUG {
		log.Printf("Rendering %v frames", z.Frames)
	}
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
	for i := uint(0); i < cnt-1; i++ {
		time += interval
		info, err := z.rescope(time)
		if err != nil {
			return nil, err
		}
		frames[i] = info
	}

	frames[cnt-1] = last
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
