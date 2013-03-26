package ai

import (
)

//------------------------------------------------------------------------------

// func TrainNominalClassifierFromBinary(
// 	features []v.F64,
// 	labels []int,
// 	labelsCardinality int,
// 	binaryTrainer BinaryClassifierTrainer) NominalClassifier {

// 	classifiers := make([]BinaryClassifier, labelsCardinality)

// 	// TODO(mike): parallelize
// 	for i := 0; i < labelsCardinality; i++ {
// 		boolLabels := make([]bool, len(labels))
// 		for j, l := range labels {
// 			boolLabels[j] = l == i
// 		}

// 		classifiers[i] = binaryTrainer(features, boolLabels)
// 	}

// 	return &compoundNominalClassifier{classifiers: classifiers}
// }

// func shuffleFeaturesAndLabels(features []v.F64, labels v.B) {
// 	for i := len(features) - 1; i > 0; i-- {
// 		j := rand.Intn(i)
// 		features[i], features[j] = features[j], features[i]
// 		labels[i], labels[j] = labels[j], labels[i]
// 	}
// }
