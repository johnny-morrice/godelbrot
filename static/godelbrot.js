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

// Open or close or zoom into the zoom box
function clickZoomBox(event) {
    return clickCallback(function () {
        switch (event.button) {
            case 0:
                Zoom.getZoomBox().zoomStep(event.clientX, event.clientY);
            break;
            case 1:
                Zoom.getZoomBox().cancel();
            break;
        }
    })
}

// Update status bar with fractal location and redraw the zoombox
function fractalSelect(event) {
    var mouseX = event.clientX;
    var mouseY = event.clientY;
    var c = imageToPlane(mouseX, mouseY);
    var status = document.getElementById("status");
    status.textContent = "r: " + c.real + " i: " + c.imag;
    Zoom.getZoomBox().move(mouseX, mouseY);
}

// Close the zoombox on escape press
function keyZoomBox(event) {
    if (event.char = "q") {
        Zoom.getZoomBox().cancel();
    }
}

// Intialize the app
function initialize(event) {
    document.addEventListener("keypress", keyZoomBox);
    restart();
}