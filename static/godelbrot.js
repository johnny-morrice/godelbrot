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

// Complex plane utilities
function defaultMandelbrotDimensions() {
    var elemDim = imageDimensions();
    var dimensions = {};
    var realMin = -2.01;
    var realMax = 0.59;
    var imagMin = -1.89;
    var imagMax = 1.58;

    var pWidth = realMax - realMin;
    var pHeight = imagMax - imagMin;

    var pAspect = pWidth / pHeight;
    var imageAspect = elemDim.imageWidth / elemDim.imageHeight;

    var expectWidth;
    if (pAspect > imageAspect) {
        // Expecting excess at bottom of image
        var taller = pWidth / imageAspect;
        // Add excess to top and bottom, in order to center image
        var resize = taller - pHeight;
        var centerResize = resize / 2;
        imagMin -= centerResize;
        imagMax += centerResize;
    } else if (pAspect < imageAspect) {
        // Expecting excess at right of image
        var fatter = pHeight * imageAspect;
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

function elemHeight(element) {
    var rect = element.getBoundingClientRect();
    return rect.height;
}