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
    
}
