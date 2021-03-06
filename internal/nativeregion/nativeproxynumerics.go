package nativeregion

import (
	"github.com/johnny-morrice/godelbrot/internal/nativesequence"
	"github.com/johnny-morrice/godelbrot/internal/region"
)

type NativeRegionProxy struct {
	*NativeRegionNumerics
	LocalRegion nativeRegion
}

// Check we implement the interface
var _ region.RegionNumerics = NativeRegionProxy{}

func (proxy NativeRegionProxy) ClaimExtrinsics() {
	proxy.NativeRegionNumerics.Region = proxy.LocalRegion
}

func (proxy NativeRegionProxy) Extrinsically(f func()) {
	old := proxy.NativeRegionNumerics.Region
	proxy.ClaimExtrinsics()
	f()
	proxy.NativeRegionNumerics.Region = old
}

type NativeSequenceProxy struct {
	*nativesequence.NativeSequenceNumerics
	LocalRegion nativeRegion
}

func (proxy NativeSequenceProxy) ClaimExtrinsics() {
	base := proxy.NativeSequenceNumerics.NativeBaseNumerics
	rectangle := proxy.LocalRegion.rect(&base)
	proxy.NativeSequenceNumerics.SubImage(rectangle)
}

func (proxy NativeSequenceProxy) Extrinsically(f func()) {
	cmin := complex(proxy.RealMin, proxy.ImagMin)
	cmax := complex(proxy.RealMax, proxy.ImagMax)

	proxy.ClaimExtrinsics()
	f()
	proxy.RealMin = real(cmin)
	proxy.ImagMin = imag(cmin)
	proxy.RealMax = real(cmax)
	proxy.ImagMax = imag(cmax)

	// Should we cache this somewhere?  New object?
	proxy.RestorePicBounds()
}
