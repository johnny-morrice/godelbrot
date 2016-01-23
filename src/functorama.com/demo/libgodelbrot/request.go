package libgodelbrot

import (
    "math/big"
    "runtime"
    "strconv"
    "encoding/json"
    "io"
    "bytes"
    "fmt"
)

type NativeInfo struct {
    UserRequest Request
    // Describe the render strategy in use
    RenderStrategy RenderMode
    // Describe the numerics system in use
    NumericsStrategy NumericsMode
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
    buff := bytes.Buffer{}
    _, rerr := buff.ReadFrom(r)

    if rerr != nil {
        return nil, rerr
    }

    return FromJSON(buff.Bytes())
}

// Available render algorithms
type RenderMode uint

const (
    AutoDetectRenderMode       = RenderMode(iota)
    RegionRenderMode
    SequenceRenderMode
    SharedRegionRenderMode
)

// Available numeric systems
type NumericsMode uint

const (
    // Functions should auto-detect the correct system for rendering
    AutoDetectNumericsMode = NumericsMode(iota)
    // Use the native CPU arithmetic operations
    NativeNumericsMode
    // Use arithmetic based around the standard library big.Float type
    BigFloatNumericsMode
)

// Request is a user description of the render to be accomplished
type Request struct {
    IterateLimit uint8
    DivergeLimit float64
    RealMin      string
    RealMax      string
    ImagMin      string
    ImagMax      string
    ImageWidth   uint
    ImageHeight  uint
    PaletteCode      string
    FixAspect        bool
    // Render algorithm
    Renderer RenderMode
    // Number of render threads
    Jobs           uint16
    RegionCollapse uint
    // Numerical system
    Numerics NumericsMode
    // Number of samples taken when detecting region render glitches
    RegionSamples uint
    // Number of bits for big.Float rendering
    Precision uint
}

func DefaultRequest() *Request {
    jobs := runtime.NumCPU() + 1
    return &Request{
        IterateLimit:   DefaultIterations,
        DivergeLimit:   DefaultDivergeLimit,
        RegionCollapse: DefaultCollapse,
        RegionSamples:  DefaultRegionSamples,
        Jobs:           uint16(jobs), // If this overflows, please send money to enable support
                                      // your ridiculous SPARC machine
        RealMin:        float2str(real(MandelbrotMin)),
        ImagMin:        float2str(imag(MandelbrotMin)),
        RealMax:        float2str(real(MandelbrotMax)),
        ImagMax:        float2str(imag(MandelbrotMax)),
        ImageHeight:    DefaultImageHeight,
        ImageWidth:     DefaultImageWidth,
        PaletteCode:    "grayscale",
    }
}

func float2str(num float64) string {
    return strconv.FormatFloat(num, 'f', -1, 64)
}