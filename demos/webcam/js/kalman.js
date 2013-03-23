/*
 * Implementation of Kalman filter.
 * A - a matrix, so that x_t = Ax_{t-1} + eps_t
 * R - covariance matrix for state-transition noise
 * C - a matrix, so that z_t = C * x_t + delta_t
 * Q - covariance matrix for measurement noise
 * mu - vector of mean values on start
 * sigma - covariance matrix on start
 */
Kalman = function(A, R, C, Q, mu, sigma) {
    this.A = A;
    this.R = R;
    this.C = C;
    this.Q = Q;
    this.mu = mu;
    this.sigma = sigma;
};

/*
 * Make prediction about mu and sigma.
 * Return [mu_pred, sigma_pred].
 */
Kalman.prototype.predict = function() {
    console.log("Kalman.predict(A=", this.A.inspect(), "C=", this.C.inspect(),
		"mu=", this.mu.inspect());
    var mu_pred = this.A.multiply(this.mu);
    console.log("mu_pred: ", mu_pred.inspect());
    var sigma_pred = this.A.multiply(this.sigma.multiply(this.A.transpose())).add(this.R);
    console.log("sigma_pred: ", sigma_pred.inspect());
    return [mu_pred, sigma_pred];
}

/*
 * Update mu and sigma by looking at the current measurements vector z.
 */
Kalman.prototype.update = function(z) {
    console.log("Kalman.update(A=", this.A.inspect(), "C=", this.C.inspect(),
		"mu=", this.mu.inspect(), "z=", z.inspect());
    var pred = this.predict();
    var mu_pred = pred[0];
    var sigma_pred = pred[1];
    console.log("mu_pred: ", mu_pred.inspect());
    console.log("sigma_pred: ", sigma_pred.inspect());
    var tmp = (this.C.multiply(sigma_pred.multiply(this.C.transpose())).add(this.Q)).inverse();
    var K = sigma_pred.multiply(this.C.transpose().multiply(tmp));
    var mu_next = mu_pred.add(K.multiply(z.subtract(this.C.multiply(mu_pred))));
    console.log("mu_next: ", mu_next.inspect());

    // sigma_next = (I - K*C) * sigma_pred
    var tmp2 = K.multiply(this.C);
    var I = Sylvester.Matrix.I(tmp2.rows());
    var sigma_next = I.subtract(tmp2).multiply(sigma_pred)
    console.log("sigma_next: ", sigma_next.inspect());
    this.mu = mu_next;
    this.sigma = sigma_next;
}
