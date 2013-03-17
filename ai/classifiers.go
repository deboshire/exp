package ai

import (
	v "github.com/deboshire/exp/math/vector"
	"math/rand"
)

// Classifies given features into two boolean classes.
// 0 <= confidence <= 1 - classifier confidence of its result.
type BinaryClassifier interface {
	Classify(features v.F64) (result bool, confidence float64)
}

type BinarryClassifierTrainer func(features []v.F64, labels []bool) BinaryClassifier

// Evaluate binary classifier on a given data.
// Returns percentage of correct hits.
func EvaluateBinaryClassifier(c BinaryClassifier, features []v.F64, labels v.B) float64 {
	if len(features) == 0 {
		return 0
	}

	successes := 0

	for i := range features {
		feature := features[i]
		label := labels[i]
		l1, _ := c.Classify(feature)
		if l1 == label {
			successes++
		}
	}

	return float64(successes) / float64(len(features))
}

type MultinomialClassifier interface {
	Classify(features v.F64) (result int, confidence float64)
}

type compoundMultinomialClassifier struct {
	classifiers []BinaryClassifier
}

func (c *compoundMultinomialClassifier) Classify(features v.F64) (result int, confidence float64) {
	panic("not implemented")
}

func TrainMultinomialClassifierFromBinary(
	features []v.F64,
	labels []int,
	labelsCardinality int,
	binaryTrainer BinarryClassifierTrainer) MultinomialClassifier {

	classifiers := make([]BinaryClassifier, labelsCardinality)

	for i := 0; i < labelsCardinality; i++ {
		boolLabels := make([]bool, len(labels))
		for j, l := range labels {
			boolLabels[j] = l == i
		}

		classifiers[i] = binaryTrainer(features, boolLabels)
	}

	return &compoundMultinomialClassifier{classifiers: classifiers}
}

func shuffleFeaturesAndLabels(features []v.F64, labels v.B) {
	for i := len(features) - 1; i > 0; i-- {
		j := rand.Intn(i)
		features[i], features[j] = features[j], features[i]
		labels[i], labels[j] = labels[j], labels[i]
	}
}

func HoldoutTestBinaryClassifier(features []v.F64, labels v.B, testingFraction float64, binaryTrainer BinarryClassifierTrainer) float64 {
	shuffleFeaturesAndLabels(features, labels)

	idx := int(float64(len(features)) * (1 - testingFraction))

	trainingFeatures := features[:idx]
	trainingLabels := labels[:idx]

	testingFeatures := features[idx:]
	testingLabels := labels[idx:]

	classifier := binaryTrainer(trainingFeatures, trainingLabels)

	return EvaluateBinaryClassifier(classifier, testingFeatures, testingLabels)
}
