function ZoomBox() {
    this.state = ZoomState.hidden();
}

ZoomBox.prototype.zoomStep = function(mouseX, mouseY) {
    switch (this.state) {
        case ZoomState.hidden():
            this.state = ZoomState.select();
            this.addBox(mouseX, mouseY);
        break;
        case ZoomState.select():
            var cDims = {};
            var boxDims = this.boxDims;
            console.log(boxDims);
            var cMin = imageToPlane(boxDims.minX, boxDims.minY);
            var cMax = imageToPlane(boxDims.maxX, boxDims.maxY);
            cDims.realMin = cMin.real;
            cDims.realMax = cMax.real;
            cDims.imagMin = cMax.imag; // Note on plane y is upside down
            cDims.imagMax = cMin.imag; 
            console.log(cDims);
            var mandelbrot = new Mandelbrot(cDims, imageDimensions());
            History.getRenderHistory().render(mandelbrot);
            this.cancel();
        break;
    }
}

ZoomBox.prototype.cancel = function() {
    this.state = ZoomState.hidden();
    this.removeBox();
}

ZoomBox.prototype.move = function(mouseX, mouseY) {
    if (this.state == ZoomState.select()) {
        this.moveBox(mouseX, mouseY);
    }
}

ZoomBox.prototype.moveBox = function(mouseX, mouseY) {
    var minX = Math.min(mouseX, this.boxAnchor.x);
    var minY = Math.min(mouseY, this.boxAnchor.y);

    var maxX = Math.max(mouseX, this.boxAnchor.x);
    var maxY = Math.max(mouseY, this.boxAnchor.y);

    var imgAspect = imageAspect();

    var boxWidth = maxX - minX;
    var boxHeight = maxY - minY;

    var boxAspect = boxWidth / boxHeight;

    if (imgAspect < boxAspect) {
        // Too fat, make thinner
        boxWidth = boxHeight * imgAspect;
    } else if (imgAspect > boxAspect) {
        // Too tall, make shorter
        boxHeight = boxWidth / imgAspect;
    }

    if (minX == this.boxAnchor.x) {
        maxX = minX + boxWidth;
    } else {
        minX = maxX - boxWidth;
    }

    if (minY == this.boxAnchor.y) {
        maxY = minY + boxHeight;
    } else {
        minY = maxY - boxHeight;
    }

    this.boxDims = {}
    this.boxDims.minX = minX;
    this.boxDims.maxX = maxX;
    this.boxDims.minY = minY;
    this.boxDims.maxY = maxY;
    this.box.style.left = minX + "px";
    this.box.style.top = minY + "px";
    this.box.style.maxWidth = boxWidth + "px";
    this.box.style.maxHeight = boxHeight + "px";
}

ZoomBox.prototype.addBox = function(x, y) {
    this.boxAnchor = {};
    this.boxAnchor.x = x;
    this.boxAnchor.y = y;
    this.box = document.createElement("div");
    this.box.style.position = "absolute";
    this.box.id = "ZoomBox";
    this.box.addEventListener("mousemove", function(self){
        return function(event) {
            self.move(event.clientX, event.clientY);
        }
    }(this));
    this.box.addEventListener("mouseclick", function(self){
        return function(event) {
            self.zoomStep(event.clientX, event.clientY);
        }
    }(this));
    document.body.appendChild(this.box);
}

ZoomBox.prototype.removeBox = function() {
    this.box.parentNode.removeChild(this.box);
    this.boxAnchor = null;
    this.box = null;
    this.boxDims = null;
}

// Singleston namespace for zoom states
var ZoomState = {}
ZoomState.hidden = function () { return 0; }
ZoomState.select = function () { return 1; }

// Singleton namespace to manage zoombox

var Zoom = {}

Zoom.getZoomBox = function() {
    if (!Zoom.zoomBox) {
        Zoom.zoomBox = new ZoomBox();
    }

    return Zoom.zoomBox;
}