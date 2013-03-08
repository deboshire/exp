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

(function() {
    var videoElement = document.querySelector('video');
    var c1 = document.getElementById("c1");
    var ctx1 = c1.getContext("2d");

    var screenCast = new ScreenCast(videoElement);

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
			
	    var end = performance.now();
	    document.getElementById('stats').innerHTML = 'Total time: ' + (end - start) + ' ms';

            var effect = this.getAttribute('data-effect');

            if (this.hasClass('active')) {
                this.classList.remove('active');
                videoElement.classList.remove(effect);
            } else {
                document.querySelectorAll('ul.effects li a').each(function(element) {
                    element.classList.remove('active');
                });

                videoElement.setAttribute('class', '');

                this.classList.add('active');
                videoElement.classList.add(effect);
            }
        });
    }
}());
