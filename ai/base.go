package ai

import (
	v "github.com/deboshire/exp/math/vector"
)

// Classifies given features into two boolean classes.
// 0 <= confidence <= 1 - classifier confidence of its result.
type BinaryClassifier interface {
	Classify(features v.F64) (result bool, confidence float64)
}

// Evaluate binary classifier on a given data.
// Returns percentage of correct hits.
func EvaluateBinaryClassifier(c BinaryClassifier, features []v.F64, labels v.B) float64 {
	successes := 0

	for i := range features {
		feature := features[i]
		label := labels[i]
		l1, _ := c.Classify(feature)
		if ; l1 == label {
			successes++
		}
	}

	return float64(successes) / float64(len(features))
}
