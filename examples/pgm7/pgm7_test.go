package pgm7

import (
	"fmt"
	"github.com/deboshire/exp/ai"
	"github.com/deboshire/exp/io/mat"
	"github.com/deboshire/exp/optimization/sgrad"
)

func benchmark(iterations int)

func ExamplePGM7_LogisticRegression() {
	trainFeatures := mat.MustRead("Train1X.mat").Array("Train1X").RowsToVectors()
	trainLabels := mat.MustRead("Train1Y.mat").Array("Train1Y").ToVector().F64ToB()

	benchmarkFeatures := mat.MustRead("Validation1X.mat").Array("Validation1X").RowsToVectors()
	benchmarkLabels := mat.MustRead("Validation1Y.mat").Array("Validation1Y").ToVector().F64ToB()

	for _, iterations := range []int{1, 10, 100, 1000} {
		fmt.Println("---\niterations: ", iterations)
		classifier := ai.TrainLogisticRegressionClassifier(
			trainFeatures,
			trainLabels,
			0,
			&sgrad.NumIterationsCriterion{NumIterations: iterations},
			1e-8)
		fmt.Println("train set: ", ai.EvaluateBinaryClassifier(classifier, trainFeatures, trainLabels))
		fmt.Println("benchmark set: ", ai.EvaluateBinaryClassifier(classifier, benchmarkFeatures, benchmarkLabels))
	}

	for _, epsilon := range []float64{1e-1, 1e-2, 1e-3} {
		fmt.Println("---\nepsilon: ", epsilon)
		classifier := ai.TrainLogisticRegressionClassifier(
			trainFeatures,
			trainLabels,
			0,
			&sgrad.RelativeMeanImprovementCriterion{},
			epsilon)
		fmt.Println("train set: ", ai.EvaluateBinaryClassifier(classifier, trainFeatures, trainLabels))
		fmt.Println("benchmark set: ", ai.EvaluateBinaryClassifier(classifier, benchmarkFeatures, benchmarkLabels))
	}

	for _, alpha := range []float64{0, 0.1, 0.2, 0.4, 0.8, 1} {
		fmt.Println("---\nalpha: ", alpha)
		classifier := ai.TrainLogisticRegressionClassifier(
			trainFeatures,
			trainLabels,
			alpha,
			&sgrad.RelativeMeanImprovementCriterion{},
			1e-2)
		fmt.Println("train set: ", ai.EvaluateBinaryClassifier(classifier, trainFeatures, trainLabels))
		fmt.Println("benchmark set: ", ai.EvaluateBinaryClassifier(classifier, benchmarkFeatures, benchmarkLabels))
	}

	// Output:
	// ---
	// iterations:  1
	// train set:  0.925
	// benchmark set:  0.915
	// ---
	// iterations:  10
	// train set:  0.985
	// benchmark set:  0.93
	// ---
	// iterations:  100
	// train set:  1
	// benchmark set:  0.925
	// ---
	// iterations:  1000
	// train set:  1
	// benchmark set:  0.92
	// ---
	// epsilon:  0.1
	// train set:  0.985
	// benchmark set:  0.925
	// ---
	// epsilon:  0.01
	// train set:  1
	// benchmark set:  0.925
	// ---
	// epsilon:  0.001
	// train set:  1
	// benchmark set:  0.93
	// ---
	// alpha:  0
	// train set:  1
	// benchmark set:  0.94
	// ---
	// alpha:  0.1
	// train set:  0.945
	// benchmark set:  0.93
	// ---
	// alpha:  0.2
	// train set:  0.94
	// benchmark set:  0.905
	// ---
	// alpha:  0.4
	// train set:  0.87
	// benchmark set:  0.855
	// ---
	// alpha:  0.8
	// train set:  0.91
	// benchmark set:  0.885
	// ---
	// alpha:  1
	// train set:  0.88
	// benchmark set:  0.87
}

func init() {
}
