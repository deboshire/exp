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
			&sgrad.NumIterationsCriterion{NumIterations: iterations},
			1e-8)
		fmt.Println("train set: ", ai.EvaluateBinaryClassifier(classifier, trainFeatures, trainLabels))
		fmt.Println("benchmark set: ", ai.EvaluateBinaryClassifier(classifier, benchmarkFeatures, benchmarkLabels))
	}

	// Output:
	// 	---
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

}

func init() {
}