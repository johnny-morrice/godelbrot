package libgodelbrot

import (
    "math/big"
    "strconv"
    "encoding/json"
    "io"
    "fmt"
    "github.com/johnny-morrice/godelbrot/config"
)

type NativeInfo struct {
    UserRequest config.Request
    // Describe the render strategy in use
    RenderStrategy config.RenderMode
    // Describe the numerics system in use
    NumericsStrategy config.NumericsMode
    PaletteType PaletteKind
    Precision uint
}

type PaletteKind uint8

const (
    Grayscale = PaletteKind(iota)
    Redscale
    Pretty
)

type BigInfo struct {
    RealMin big.Float
    RealMax big.Float
    ImagMin big.Float
    ImagMax big.Float
}

type SerialBigInfo struct {
    RealMin string
    RealMax string
    ImagMin string
    ImagMax string
}

// Info completely describes the render process
type Info struct {
    NativeInfo
    BigInfo
}

func (info *Info) bignums() []*big.Float {
    return []*big.Float{
        &info.RealMin,
        &info.RealMax,
        &info.ImagMin,
        &info.ImagMax,
    }
}

// IsAccurate returns True if the bignums used internally by info are all accurate.
func (info *Info) IsAccurate() bool {
    for _, x := range info.bignums() {
        if x.Acc() == big.Below {
            return false
        }
    }
    return true
}

// AddPrec increases the precision of all Infos internal bignums by delta bits.
func (info *Info) AddPrec(delta int) {
    nextPrec := int(info.Precision) + delta
    if nextPrec <= 0 {
        msg := fmt.Sprintf("Invalid precision: %v", nextPrec)
        panic(msg)
    }
    prec := uint(nextPrec)
    info.Precision = prec
    for _, x := range info.bignums() {
        x.SetPrec(prec)
    }
}

// Generate a user request that corresponts to the info numerics
func (info *Info) GenRequest() config.Request {
    req := info.UserRequest
    req.RealMin = emitBig(&info.RealMin)
    req.RealMax = emitBig(&info.RealMax)
    req.ImagMin = emitBig(&info.ImagMin)
    req.ImagMax = emitBig(&info.ImagMax)

    return req
}

// UserInfo is a variant of Info that can be easily serialized
type UserInfo struct {
    NativeInfo
    SerialBigInfo
}

func Friendly(desc *Info) *UserInfo {
    userDesc := &UserInfo{}
    userDesc.NativeInfo = desc.NativeInfo
    userDesc.RealMin = emitBig(&desc.RealMin)
    userDesc.RealMax = emitBig(&desc.RealMax)
    userDesc.ImagMin = emitBig(&desc.ImagMin)
    userDesc.ImagMax = emitBig(&desc.ImagMax)
    return userDesc
}

func Unfriendly(userDesc *UserInfo) (*Info, error) {
    desc := &Info{}
    desc.NativeInfo = userDesc.NativeInfo

    ubnds := []string{
        userDesc.RealMin,
        userDesc.RealMax,
        userDesc.ImagMin,
        userDesc.ImagMax,
    }
    bnds := make([]*big.Float, len(ubnds))

    for i, u := range ubnds {
        b, err := parseBig(u)
        if err != nil {
            return nil, fmt.Errorf("Error parsing bound: %v", err)
        }
        b.SetPrec(userDesc.Precision)
        bnds[i] = b
    }

    desc.RealMin = *bnds[0]
    desc.RealMax = *bnds[1]
    desc.ImagMin = *bnds[2]
    desc.ImagMax = *bnds[3]

    return desc, nil
}

func ToJSON(desc *Info) ([]byte, error) {
    userDesc := Friendly(desc)
    return json.MarshalIndent(userDesc, "", "    ")
}

func FromJSON(format []byte) (*Info, error) {
    userDesc := new(UserInfo)
    err := json.Unmarshal(format, userDesc)
    if err == nil {
        desc, converr := Unfriendly(userDesc)
        return desc, converr
    } else {
        return nil, err
    }
}

func WriteInfo(w io.Writer, desc *Info) error {
    text, jerr := ToJSON(desc)
    if jerr != nil {
        return jerr
    }
    _, werr := w.Write(text)
    return werr
}

func ReadInfo(r io.Reader) (*Info, error) {
    ui := &UserInfo{}
    dec := json.NewDecoder(r)
    err := dec.Decode(ui)
    if err != nil {
        return nil, err
    }
    return Unfriendly(ui)
}

type InfoPkt struct {
    Info *Info
    Err error
}

type uipkt struct {
    ui *UserInfo
    err error
}

func ReadInfoStream(r io.Reader) <-chan InfoPkt {
    uich := make(chan uipkt)
    go func() {
        dec := json.NewDecoder(r)
        for i := 0; dec.More(); i++ {
            ui := &UserInfo{}
            readerr := dec.Decode(ui)
            if readerr != nil {
                message := fmt.Errorf("Error after %v JSON objects: %v", i, readerr)
                uich<- uipkt{err: message,}
                continue
            }
            uich<- uipkt{ui: ui,}
        }
        close(uich)
    }()

    infoch := make(chan InfoPkt)
    go func() {
        for uipkt := range uich {
            if uipkt.err != nil {
                infoch<- InfoPkt{Err: uipkt.err,}
                continue
            }
            inf, err := Unfriendly(uipkt.ui)
            if err != nil {
                infoch<- InfoPkt{Err: uipkt.err,}
                continue
            }
            infoch<- InfoPkt{Info: inf,}
        }
        close(infoch)
    }()

    return infoch
}

func DefaultRequest() *config.Request {
    return &config.Request {
        IterateLimit:   DefaultIterations,
        DivergeLimit:   DefaultDivergeLimit,
        RegionCollapse: DefaultCollapse,
        RegionSamples:  DefaultRegionSamples,
        RealMin:        float2str(real(MandelbrotMin)),
        ImagMin:        float2str(imag(MandelbrotMin)),
        RealMax:        float2str(real(MandelbrotMax)),
        ImagMax:        float2str(imag(MandelbrotMax)),
        ImageHeight:    DefaultImageHeight,
        ImageWidth:     DefaultImageWidth,
        FixAspect:      config.Shrink,
        PaletteCode:    "grayscale",
        Jobs: 1,
    }
}

func float2str(num float64) string {
    return strconv.FormatFloat(num, 'f', -1, 64)
}