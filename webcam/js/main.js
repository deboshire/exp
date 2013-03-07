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

(function() {
    var videoElement = document.querySelector('video');
    var c1 = document.getElementById("c1");
    var ctx1 = c1.getContext("2d");
    var c2 = document.getElementById("c2");
    var ctx2 = c2.getContext("2d");

    var screenCast = new ScreenCast(videoElement);

    screenCast.start();

    var effectButtons = document.querySelectorAll('ul.effects li a');
    for (var i = 0, numberOfButtons = effectButtons.length; i < numberOfButtons; i++) {
        effectButtons[i].addEventListener('click', function(e) {
            e.preventDefault();
	    var width = videoElement.videoWidth;
	    var height = videoElement.videoHeight;
	    ctx1.drawImage(videoElement, 0, 0, width, height);
	    var frame = ctx1.getImageData(0, 0, width, height);
	    var l = frame.data.length / 4;
	    for (var i = 0; i < l; i++) {
		var r = frame.data[i*4];
		var g = frame.data[i*4+1];
		var b = frame.data[i*4+2];
		if (g > 100 && r > 100 && b > 43) {
		    frame.data[i*4] = 0;
		    frame.data[i*4+1] = 0;
		    frame.data[i*4+2] = 0;
		}
	    }
	    ctx2.putImageData(frame, 0, 0);

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
