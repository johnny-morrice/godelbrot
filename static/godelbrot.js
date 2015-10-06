// History of all fractals rendered
function RenderHistory() {
    this.hist = [];
}

// Resize and render the fractal on screen
RenderHistory.prototype.resizeRender = function() {
    var next = this.hist[0].resize();
    next.render();
}

// Render the last fractal in the history
RenderHistory.prototype.fractalBack = function() {
    this.hist.pop();
    this.resizeRender();
}

// Render a new mandelbrot
RenderHistory.prototype.render = function(mandelbrot) {
    mandelbrot.render();
    this.hist.push(mandelbrot);
}

// Clear all history
RenderHistory.prototype.clear = function() {
    this.hist = [];
}

// Last Mandelbrot render command
RenderHistory.prototype.last = function() {
    return this.hist[0];
}

// Singleton namespace for fractal history
var History = {};

History.getRenderHistory = function() {
    if (!History.renderHistory) {
        History.renderHistory = new RenderHistory();
    }
    return History.renderHistory;
}

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

// Restart the application
function restart() {
    var renderHistory = History.getRenderHistory();
    renderHistory.clear();

    var mandelbrot = new Mandelbrot(defaultMandelbrotDimensions(), 
        imageDimensions());

    renderHistory.render(mandelbrot);
}

// Click-suitable callback for restarting the app
function restartClick() {
    return clickCallback(function() { 
        restart() 
    });
}

// Click suitable callback for going back one step
function fractalBackClick() {
    return clickCallback(function () { 
        History.getRenderHistory().fractalBack(); 
    });
}

// Disable processing of an anchor's href attribute
function clickCallback(callback) {
    callback();
    return false;
}

// Resize the application
function resize() {
    History.getRenderHistory().resizeRender();
}

// Open or close and zoom into the zoom box
function fractalZoomBox(event) {
    return clickCallback(function () {

    })
}

// Update status bar with fractal location and redraw the zoombox
function fractalSelect(event) {
    var c = imageToPlane(event.clientX, event.clientY);
    var status = document.getElementById("status");
    status.textContent = "r: " + c.real + " i: " + c.imag;
}

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

// Screen utilities

function imageDimensions() {
    var toolbar = document.getElementById('toolbar');
    // Assume standards compliance for now
    var dimensions = {};
    dimensions.imageWidth = window.innerWidth;
    dimensions.imageHeight = window.innerHeight - elemHeight(toolbar);
    return dimensions;
}

// Aspect ratio of image
function imageAspect() {
    var elemDim = imageDimensions();
    return elemDim.imageWidth / elemDim.imageHeight;
}

function elemHeight(element) {
    var rect = element.getBoundingClientRect();
    return rect.height;
}