#pragma once


double dist2(void* ap, void* bp, int len) {
	int i;
	double* a = (double*)ap;
	double* b = (double*)bp;

	double result = 0.0;

	for (i = 0; i < len; ++i) {
		double x = a[i] - b[i];
		result += x * x;
	}

	return result;
}

double dot(void* ap, void* bp, int len) {
	int i;
	double* a = (double*)ap;
	double* b = (double*)bp;

	double result = 0.0;

	for (i = 0; i < len; ++i) {
		result += a[i] * b[i];
	}

	return result;
}
