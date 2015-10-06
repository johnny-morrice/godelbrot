// A single rendering of the mandelbrot fractal
function Mandelbrot(fractalDimensions, docDimensions) {
    this.realMin = fractalDimensions.realMin;
    this.realMax = fractalDimensions.realMax;
    this.imagMin = fractalDimensions.imagMin;
    this.imagMax = fractalDimensions.imagMax;
    this.imageWidth = docDimensions.imageWidth;
    this.imageHeight = docDimensions.imageHeight;
}

// Render the Mandelbrot set
Mandelbrot.prototype.render = function() {
    var packet = {
        "Command": "render",
        "Render": {
            "ImageWidth": this.imageWidth,
            "ImageHeight": this.imageHeight,
            "RealMin": this.realMin,
            "RealMax": this.realMax,
            "ImagMin": this.imagMin,
            "ImagMax": this.imagMax
        }
    };
    var json = JSON.stringify(packet);
    var encoded = encodeURIComponent(json);
    var service = "/service?godelbrotPacket=" + encoded;
    var fractal = document.getElementById("fractal");
    fractal.src = service
}

// Create new Mandelbrot viewing same area of set but at different image size
Mandelbrot.prototype.resize = function() {
    var fractalDimensions = {};
    fractalDimensions.realMin = this.realMin;
    fractalDimensions.realMax = this.realMax;
    fractalDimensions.imagMin = this.imagMin;
    fractalDimensions.imagMax = this.imagMax;
    return new Mandelbrot(fractalDimensions, imageDimensions());
}

// Conversion unit from image to plane
Mandelbrot.prototype.planeUnits = function() {
    var width = this.realMax - this.realMin;
    var height = this.imagMax - this.imagMin;
    var elemDim = imageDimensions(); 
    var units = {};
    units.real = width / elemDim.imageWidth;
    units.imag = height / elemDim.imageHeight;
    return units;
}