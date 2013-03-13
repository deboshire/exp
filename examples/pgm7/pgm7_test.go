package pgm7

import (
	"fmt"
	"github.com/deboshire/exp/ai"
	"github.com/deboshire/exp/io/mat"
	"github.com/deboshire/exp/optimization/sgrad"
)

func benchmark(iterations int)

func ExamplePGM7_LogisticRegression() {
	trainFeatures := mat.ReadFileOrPanic("Train1X.mat").GetArray("Train1X").RowsToVectors()
	trainLabels := mat.ReadFileOrPanic("Train1Y.mat").GetArray("Train1Y").ToVector().F64ToB()

	benchmarkFeatures := mat.ReadFileOrPanic("Validation1X.mat").GetArray("Validation1X").RowsToVectors()
	benchmarkLabels := mat.ReadFileOrPanic("Validation1Y.mat").GetArray("Validation1Y").ToVector().F64ToB()

	// This will be severely overfitted without regularization.
	// benchmark precision is supposed to be low.
	// In fact it is almost 50% :)
	for _, iterations := range []int{100, 1000, 10000} {
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
	// ---
	// iterations:  100
	// train set:  0.965
	// benchmark set:  0.515
	// ---
	// iterations:  1000
	// train set:  0.995
	// benchmark set:  0.51
	// ---
	// iterations:  10000
	// train set:  1
	// benchmark set:  0.52
}
