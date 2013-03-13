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
	terminationCriterion sgrad.TerminationCriterion,
	epsilon float64) BinaryClassifier {
	y, x := sgrad.Minimize(
		logisticRegressionCostFunction(features, labels),
		v.Zeroes(len(features[0])),
		epsilon,
		terminationCriterion,
		nil)

	return &logisticRegressionClassifier{cost: y, theta: x}
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// http://mathurl.com/bmfs3db
func logisticRegressionCostFunction(features []v.F64, labels v.B) sgrad.ObjectiveFunc {
	f := func(idx int, x v.F64) (value float64, gradient v.F64) {
		feature := features[idx]
		label := labels[idx]

		h := sigmoid(x.DotProduct(feature))
		gradient = feature.Copy()

		if label {
			value = -math.Log(h)
			gradient.Mul(h - 1.0)
		} else {
			value = -math.Log(1.0 - h)
			gradient.Mul(h)
		}

		return
	}

	return sgrad.ObjectiveFunc{Terms: len(features), F: f}
}

func (c *logisticRegressionClassifier) Classify(features v.F64) (result bool, confidence float64) {
	h := sigmoid(c.theta.DotProduct(features))
	return h >= 0.5, math.Abs(0.5-h) * 2.0
}
