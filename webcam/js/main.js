'use strict';

HTMLElement.prototype.hasClass = function(className) {
    var pattern = new RegExp('\\b' + className + '\\b');

    return pattern.test(this.className);
};

NodeList.prototype.each = function(callback) {
    for (var i = 0, numberOfElements = this.length; i < numberOfElements; i++) {
        callback(this[i]);
    }
};

function drawFaces(ctx, comp, scale) {
    ctx.lineWidth = 2;
    ctx.strokeStyle = 'rgba(230,87,0,0.8)';
    /* draw detected area */
    for (var i = 0; i < comp.length; i++) {
	ctx.beginPath();
	ctx.arc((comp[i].x + comp[i].width * 0.5) * scale, (comp[i].y + comp[i].height * 0.5) * scale,
		(comp[i].width + comp[i].height) * 0.25 * scale * 1.2, 0, Math.PI * 2);
	ctx.stroke();
    }
}

function scheduleSearch() {
    setTimeout(function() {
	document.getElementById('findFacesButton').click();
    }, 100);
}

var FaceModel = function() {
    this.mu = $M([
	[0],
	[0],
	[0],
	[0],
	[320],
	[240]
    ]);
    this.sigma = $M([
	[1600, 0, 0, 0, 0, 0],
	[0, 1600, 0, 0, 0, 0],
	[0, 0, 1600, 0, 0, 0],
	[0, 0, 0, 1600, 0, 0],
	[0, 0, 0, 0, 1600, 0],
	[0, 0, 0, 0, 0, 1600],
    ]);
    this.lastUpdated = performance.now();
};

FaceModel.prototype.update = function(ctx, comp) {
    var now = performance.now();
    var dt = (now - this.lastUpdated)/1000;
    this.lastUpdated = now;

    var A = $M([
	[1, 0, dt, 0, 0, 0],
	[0, 1, 0, dt, 0, 0],
	[0, 0, 1, 0, 0, 0],
	[0, 0, 0, 1, 0, 0],
	[0, 0, 0, 0, 1, 0],
	[0, 0, 0, 0, 0, 1]
    ]);
    var R = $M([
	[400,   0,   0,   0,   0,   0],
	[  0, 400,   0,   0,   0,   0],
	[  0,   0, 200,   0,   0,   0],
	[  0,   0,   0, 200,   0,   0],
	[  0,   0,   0,   0, 400,   0],
	[  0,   0,   0,   0,   0, 400]
    ]);
    var C = $M([
	[1, 0, 0, 0, 0, 0],
	[0, 1, 0, 0, 0, 0],
	[0, 0, 0, 0, 1, 0],
	[0, 0, 0, 0, 0, 1]
    ]);

    var Q;
    var z;
    if (comp.length == 0) {
	z = $M([
	    [0],
	    [0],
	    [320],
	    [240]
	    ]);
	Q = $M([
	    [1600, 0, 0, 0],
	    [0, 1600, 0, 0],
	    [0, 0, 1600, 0],
	    [0, 0, 0, 1600]
	]);
    } else {
	z = $M([
	    [comp[0].x],
	    [comp[0].y],
	    [comp[0].width],
	    [comp[0].height]
	]);
	Q = $M([
	    [90, 0, 0, 0],
	    [ 0,90, 0, 0],
	    [ 0, 0,90, 0],
	    [ 0, 0, 0,90]
	]);
    }
    var kalman = new Kalman(A, R, C, Q, this.mu, this.sigma);
    kalman.update(z);
    this.mu = kalman.mu;
    this.sigma = kalman.sigma;
    console.log("new mu: ", this.mu.inspect(), "new sigma: ", this.sigma.inspect());

    ctx.lineWidth = 2;
    ctx.strokeStyle = 'rgba(87,0,230, 0.8)';
    ctx.beginPath();
    var kx = this.mu.e(1, 1);
    var ky = this.mu.e(2, 1);
    var kw = this.mu.e(4, 1);
    var kh = this.mu.e(5, 1);
    console.log("kx: ", kx, "ky: ", ky, "kw: ", kw, "kh: ", kh);
    ctx.arc((kx + kw * 0.5), (ky + kh * 0.5),
	    (kw + kh) * 0.25 * 1.2, 0, Math.PI * 2);
    ctx.stroke();

};

(function() {
    var videoElement = document.querySelector('video');
    var c1 = document.getElementById("c1");
    var ctx1 = c1.getContext("2d");

    var screenCast = new ScreenCast(videoElement);

    var faceModel = new FaceModel();

    screenCast.start();

    var effectButtons = document.querySelectorAll('ul.effects li a');
    for (var i = 0, numberOfButtons = effectButtons.length; i < numberOfButtons; i++) {
        effectButtons[i].addEventListener('click', function(e) {
            e.preventDefault();
	    var start = performance.now();
	    var width = videoElement.videoWidth;
	    var height = videoElement.videoHeight;
	    ctx1.drawImage(videoElement, 0, 0, width, height);
	    var comp = ccv.detect_objects({ "canvas" : ccv.grayscale(c1),
					    "cascade" : cascade,
					    "interval" : 5,
					    "min_neighbors" : 1 });
	    console.log(comp);
	    drawFaces(ctx1, comp, 1);
	    faceModel.update(ctx1, comp);
			
	    var end = performance.now();
	    document.getElementById('stats').innerHTML = 'Total time: ' + (end - start) + ' ms';
	    scheduleSearch();
        });
    }

}());
