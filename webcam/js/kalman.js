/*
 * Implementation of Kalman filter.
 * A - a matrix, so that x_t = Ax_{t-1} + eps_t
 * R - covariance matrix for state-transition noise
 * C - a matrix, so that z_t = C * x_t + delta_t
 * Q - covariance matrix for measurement noise
 * mu - vector of mean values
 * sigma - covariance matrix
 * z - measurements vector
 * Returns [mu_next, sigma_next] array.
 */
function Kalman(A, R, C, Q, mu, sigma, z) {
    console.log("Kalman(A=", A.inspect(), "C=", C.inspect(),
		"mu=", mu.inspect(), "z=", z.inspect());
    var mu_pred = A.multiply(mu);
    console.log("mu_pred: ", mu_pred.inspect());
    var sigma_pred = A.multiply(sigma.multiply(A.transpose())).add(R);
    console.log("sigma_pred: ", sigma_pred.inspect());
    var tmp = (C.multiply(sigma_pred.multiply(C.transpose())).add(Q)).inverse();
    var K = sigma_pred.multiply(C.transpose().multiply(tmp));
    console.log("K: ", K.inspect());
    var mu_next = mu_pred.add(K.multiply(z.subtract(C.multiply(mu_pred))));
    console.log("mu_next: ", mu_next.inspect());

    // sigma_next = (I - K*C) * sigma_pred
    var tmp2 = K.multiply(C);
    var I = Sylvester.Matrix.I(tmp2.rows());
    var sigma_next = I.subtract(tmp2).multiply(sigma_pred)
    console.log("sigma_next: ", sigma_next.inspect());
    return [mu_next, sigma_next];
}
