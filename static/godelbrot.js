// 
function render(realMin, realMax, imagMin, imagMax) {
    var fractal = document.getElementById('fractal');
    var aspect = 16.0 / 9.0;
    var height = Math.round(fractal.width / aspect);
    var packet = {
        "Command": "render",
        "Render": {
            "ImageWidth": fractal.width,
            "ImageHeight": height,
            "RealMin": realMin,
            "RealMax": realMax,
            "ImagMin": imagMin,
            "ImagMax": imagMax
        }
    };
    var json = JSON.stringify(packet);
    var encoded = encodeURIComponent(json);
    var service = "/service?godelbrotPacket=" + encoded;
    fractal.src = service
}