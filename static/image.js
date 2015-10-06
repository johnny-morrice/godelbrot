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