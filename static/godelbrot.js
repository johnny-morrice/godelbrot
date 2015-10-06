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
    var dimensions = {};
    dimensions.realMin = -2.01;
    dimensions.realMax = 0.59;
    dimensions.imagMin = -1.89;
    dimensions.imagMax = 1.60;
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