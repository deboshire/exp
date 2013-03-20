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

function drawKalman(ctx, mu, sigma) {
    ctx.lineWidth = 2;
    ctx.strokeStyle = 'rgba(87,0,230, 0.8)';
    ctx.beginPath();
    var kx = mu.e(1, 1);
    var ky = mu.e(2, 1);
    var kw = mu.e(5, 1);
    var kh = mu.e(6, 1);
    console.log("kx: ", kx, "ky: ", ky, "kw: ", kw, "kh: ", kh);
    ctx.arc(kx, ky, (kw + kh) * 0.25 * 1.2, 0, Math.PI * 2);
    ctx.stroke();
}

var FaceModel = function() {
    this.mu = $M([
	[320],
	[240],
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
	[200,   0,   0,   0,   0,   0],
	[  0, 200,   0,   0,   0,   0],
	[  0,   0, 200,   0,   0,   0],
	[  0,   0,   0, 200,   0,   0],
	[  0,   0,   0,   0, 50,   0],
	[  0,   0,   0,   0,   0, 50]
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
	Q = $M([
	    [3600, 0, 0, 0],
	    [0, 3600, 0, 0],
	    [0, 0, 3600, 0],
	    [0, 0, 0, 3600]
	]);
	var kalman = new Kalman(A, R, C, Q, this.mu, this.sigma);
	var pred = kalman.predict();
	console.log("mu_pred: ", pred[0].inspect(), "sigma_pred: ", pred[1].inspect());
	drawKalman(ctx, pred[0], pred[1]);
    } else {
	z = $M([
	    [comp[0].x+comp[0].width/2],
	    [comp[0].y+comp[0].height/2],
	    [comp[0].width],
	    [comp[0].height]
	]);
	console.log("x: ", comp[0].x, "y: ", comp[0].y, "w: ", comp[0].width, "h: ", comp[0].height);
	Q = $M([
	    [900, 0, 0, 0],
	    [ 0,900, 0, 0],
	    [ 0, 0,900, 0],
	    [ 0, 0, 0,900]
	]);
	var kalman = new Kalman(A, R, C, Q, this.mu, this.sigma);
	kalman.update(z);
	this.mu = kalman.mu;
	this.sigma = kalman.sigma;
	console.log("new mu: ", this.mu.inspect(), "new sigma: ", this.sigma.inspect());
	drawKalman(ctx, this.mu, this.sigma);
    }
};

(function() {
    var videoElement = document.querySelector('video');
    var c1 = document.getElementById("c1");
    var ctx1 = c1.getContext("2d");

    var screenCast = new ScreenCast(videoElement);

    var faceModel = new FaceModel();

    screenCast.start();

    document.getElementById('findFacesButton').addEventListener('click', function(e) {
        e.preventDefault();
	var start = performance.now();
	console.time("total");
	var width = videoElement.videoWidth;
	var height = videoElement.videoHeight;
	console.time("drawImage");
	ctx1.drawImage(videoElement, 0, 0, width, height);
	console.timeEnd("drawImage");
	console.time("detectObjects");
	var comp = ccv.detect_objects({ "canvas" : ccv.grayscale(c1),
					"cascade" : cascade,
					"interval" : 5,
					"min_neighbors" : 1 });
	console.timeEnd("detectObjects");
	console.log(comp);
	console.time("drawFaces");
	drawFaces(ctx1, comp, 1);
	console.timeEnd("drawFaces");
	console.time("updateFaceModel");
	faceModel.update(ctx1, comp);
	console.timeEnd("updateFaceModel");
			
	var end = performance.now();
	console.timeEnd("total");
	document.getElementById('stats').innerHTML = 'Total time: ' + (end - start) + ' ms';
	scheduleSearch();
    });

}());
