// Complex plane utilities
function imageToPlane(x, y) {
    var mandelbrot = History.getRenderHistory().last();
    var planeUnits = mandelbrot.planeUnits();
    var absR = x * planeUnits.real;
    var absI = y * planeUnits.imag;
    var r = mandelbrot.realMin + absR;
    var i = mandelbrot.imagMax - absI;
    var c = {};
    c.real = r;
    c.imag = i;
    return c;
}

function defaultMandelbrotDimensions() {
    var dimensions = {};
    var realMin = -2.01;
    var realMax = 0.59;
    var imagMin = -1.89;
    var imagMax = 1.58;

    var pWidth = realMax - realMin;
    var pHeight = imagMax - imagMin;

    var pAspect = pWidth / pHeight;
    var iAspect = imageAspect();

    var expectWidth;
    if (pAspect > iAspect) {
        // Expecting excess at bottom of image
        var taller = pWidth / iAspect;
        // Add excess to top and bottom, in order to center image
        var resize = taller - pHeight;
        var centerResize = resize / 2;
        imagMin -= centerResize;
        imagMax += centerResize;
    } else if (pAspect < iAspect) {
        // Expecting excess at right of image
        var fatter = pHeight * iAspect;
        var resize = fatter - pWidth;
        var centerResize = resize / 2;
        realMin -= centerResize;
        realMax += centerResize;
    }

    dimensions.realMax = realMax;
    dimensions.realMin = realMin;
    dimensions.imagMax = imagMax;
    dimensions.imagMin = imagMin;

    return dimensions;
}