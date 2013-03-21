package ai

import (
	v "github.com/deboshire/exp/math/vector"
	"github.com/deboshire/exp/optimization/sgrad"
	"math"
)

type logisticRegressionClassifier struct {
	cost  float64
	theta v.F64
}

func TrainLogisticRegressionClassifier(
	features []v.F64,
	labels v.B,
	lambda float64,
	termCrit sgrad.TermCrit,
	eps float64) BinaryClassifier {
	y, x := sgrad.Minimize(
		logisticRegressionCostFunction(features, labels, lambda),
		v.Zeroes(len(features[0])),
		eps,
		termCrit,
		nil)

	return &logisticRegressionClassifier{cost: y, theta: x}
}

func NewLogisticRegressionTrainer(
	lambda float64,
	termCrit sgrad.TermCrit,
	eps float64) BinaryClassifierTrainer {
	return func(features []v.F64, labels []bool) BinaryClassifier {
		return TrainLogisticRegressionClassifier(features, labels, lambda, termCrit, eps)
	}
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// http://mathurl.com/bmfs3db
func logisticRegressionCostFunction(features []v.F64, labels v.B, lambda float64) sgrad.ObjectiveFunc {
	f := func(idx int, x v.F64, gradient v.F64) (value float64) {
		feature := features[idx]
		label := labels[idx]

		h := sigmoid(x.DotProduct(feature))
		feature.CopyTo(gradient)

		if label {
			value = -math.Log(h)
			gradient.Mul(h - 1.0)
		} else {
			value = -math.Log(1.0 - h)
			gradient.Mul(h)
		}

		if lambda != 0.0 {
			// apply regularizaiton.
			for i := 1; i < len(x); i++ {
				value += 0.5 * lambda * x[i] * x[i]
				gradient[i] += lambda * x[i]
			}
		}

		return
	}

	return sgrad.ObjectiveFunc{Terms: len(features), F: f}
}

func (c *logisticRegressionClassifier) Classify(features v.F64) (res bool, confidence float64) {
	h := sigmoid(c.theta.DotProduct(features))
	return h >= 0.5, math.Abs(0.5-h) * 2.0
}
