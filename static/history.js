// History of all fractals rendered
function RenderHistory() {
    this.hist = [];
}

// Resize and render the fractal on screen
RenderHistory.prototype.resizeRender = function() {
    this.last().resize().render();
}

// Render the last fractal in the history
RenderHistory.prototype.fractalBack = function() {
    var last = this.hist.pop();
    if (!this.hist.length) {
        this.hist.push(last);
    }
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
    return this.hist[this.hist.length - 1];
}

// Singleton namespace for fractal history
var History = {};

History.getRenderHistory = function() {
    if (!History.renderHistory) {
        History.renderHistory = new RenderHistory();
    }
    return History.renderHistory;
}