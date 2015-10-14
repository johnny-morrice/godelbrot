package libgodelbrot

import (
    "testing"
)

type bnMockRenderApp mockRenderApplication

func (mock bnMockRenderApp) okay() bool {
    return mock.tLimits && mock.tPictureWidth && mock.tPictureHeight
}


func TestCreateBaseNumerics(t *testing.T) {
    mock := bnMockRenderApp(mockRenderApplication{})
    mock.pictureWidth = 10
    mock.pictureHeight = 20
    mock.iterLimit = 200

    numerics := CreateBaseNumerics(mock)

    if !mock.okay() {
        t.Error("Did not call mock.Limit()")
    }

    okay := numerics.picXMin == 0
    okay &&= numerics.picXMax == 10
    okay &&= numerics.picYMin == 0
    okay &&= numerics.picYMax == 20
    okay &&= numerics.iterLimit == 200

    if !okay {
        t.Error("numerics had unexpected value", numerics)
    }
}