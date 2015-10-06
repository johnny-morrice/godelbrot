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
            // For now, just disable the zoom
            this.cancel();
    }
}

ZoomBox.prototype.cancel = function() {
    this.state = ZoomState.hidden();
    this.removeBox();
}

ZoomBox.prototype.move = function(x, y) {
    if (this.state == ZoomState.select()) {
        this.moveBox(x, y);
    }
}

ZoomBox.prototype.moveBox = function(x, y) {
    var minX = Math.min(x, this.boxAnchor.x);
    var minY = Math.min(y, this.boxAnchor.y);

    var maxX = Math.max(x, this.boxAnchor.x);
    var maxY = Math.max(y, this.boxAnchor.y);

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

    this.boxTopLeft = {}
    this.boxTopLeft.x = minX;
    this.boxTopLeft.y = minY;
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
    this.boxTopLeft = null;
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