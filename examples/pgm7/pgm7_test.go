package pgm7

import (
	"fmt"
	"github.com/deboshire/exp/ai"
	"github.com/deboshire/exp/io/mat"
	v "github.com/deboshire/exp/math/vector"
	"github.com/deboshire/exp/optimization/gcsearch"
	"github.com/deboshire/exp/optimization/sgrad"
	"math/rand"
)

func readTrainData() (trainFeatures []v.F64, trainLabels v.B) {
	trainFeatures = mat.MustRead("Train1X.mat").Array("Train1X").RowsToVectors()
	trainLabels = mat.MustRead("Train1Y.mat").Array("Train1Y").ToVector().F64ToB()
	return
}

func readBenchmarkData() (benchmarkFeatures []v.F64, benchmarkLabels v.B) {
	benchmarkFeatures = mat.MustRead("Validation1X.mat").Array("Validation1X").RowsToVectors()
	benchmarkLabels = mat.MustRead("Validation1Y.mat").Array("Validation1Y").ToVector().F64ToB()
	return
}

func ExamplePGM7_LogisticRegression_HoldoutTesting() {
	rand.Seed(98765)
	trainFeatures, trainLabels := readTrainData()

	for _, testingFraction := range []float64{0.5, 0.25, 0.1, 0.05} {
		fmt.Println("---\nfraction: ", testingFraction)
		score := ai.HoldoutTestBinaryClassifier(
			trainFeatures,
			trainLabels,
			testingFraction,
			ai.NewLogisticRegressionTrainer(0,
				&sgrad.NumIterationsCrit{NumIterations: 10},
				1e-8))
		fmt.Println("Holdout testing:", score)
	}

	// Output:
	// ---
	// fraction:  0.5
	// Holdout testing: 0.84
	// ---
	// fraction:  0.25
	// Holdout testing: 0.92
	// ---
	// fraction:  0.1
	// Holdout testing: 0.85
	// ---
	// fraction:  0.05
	// Holdout testing: 1
}

func ExamplePGM7_LogisticRegression_Iterations() {
	rand.Seed(98765)
	trainFeatures, trainLabels := readTrainData()
	benchmarkFeatures, benchmarkLabels := readBenchmarkData()

	for _, iterations := range []int{1, 10, 100, 1000} {
		fmt.Println("---\niterations: ", iterations)
		classifier := ai.TrainLogisticRegressionClassifier(
			trainFeatures,
			trainLabels,
			0,
			&sgrad.NumIterationsCrit{NumIterations: iterations},
			1e-8)
		fmt.Println("train set: ", ai.EvaluateBinaryClassifier(classifier, trainFeatures, trainLabels))
		fmt.Println("benchmark set: ", ai.EvaluateBinaryClassifier(classifier, benchmarkFeatures, benchmarkLabels))
	}

	// Output:
	// ---
	// iterations:  1
	// train set:  0.955
	// benchmark set:  0.94
	// ---
	// iterations:  10
	// train set:  0.985
	// benchmark set:  0.925
	// ---
	// iterations:  100
	// train set:  1
	// benchmark set:  0.925
	// ---
	// iterations:  1000
	// train set:  1
	// benchmark set:  0.93
}

func ExamplePGM7_LogisticRegression_Epsilon() {
	rand.Seed(98765)
	trainFeatures, trainLabels := readTrainData()
	benchmarkFeatures, benchmarkLabels := readBenchmarkData()

	for _, epsilon := range []float64{1e-1, 1e-2, 1e-3} {
		fmt.Println("---\nepsilon: ", epsilon)
		classifier := ai.TrainLogisticRegressionClassifier(
			trainFeatures,
			trainLabels,
			0,
			&sgrad.RelativeMeanImprovementCrit{},
			epsilon)
		fmt.Println("train set: ", ai.EvaluateBinaryClassifier(classifier, trainFeatures, trainLabels))
		fmt.Println("benchmark set: ", ai.EvaluateBinaryClassifier(classifier, benchmarkFeatures, benchmarkLabels))
	}

	// Output:
	// ---
	// epsilon:  0.1
	// train set:  0.99
	// benchmark set:  0.94
	// ---
	// epsilon:  0.01
	// train set:  0.99
	// benchmark set:  0.925
	// ---
	// epsilon:  0.001
	// train set:  1
	// benchmark set:  0.93
}

func ExamplePGM7_LogisticRegression_Lambda() {
	rand.Seed(98765)
	trainFeatures, trainLabels := readTrainData()
	benchmarkFeatures, benchmarkLabels := readBenchmarkData()

	for _, lambda := range []float64{0, 0.1, 0.2, 0.4, 0.8, 1} {
		fmt.Println("---\nlambda: ", lambda)
		classifier := ai.TrainLogisticRegressionClassifier(
			trainFeatures,
			trainLabels,
			lambda,
			&sgrad.RelativeMeanImprovementCrit{},
			1e-2)
		fmt.Println("train set: ", ai.EvaluateBinaryClassifier(classifier, trainFeatures, trainLabels))
		fmt.Println("benchmark set: ", ai.EvaluateBinaryClassifier(classifier, benchmarkFeatures, benchmarkLabels))
	}

	// Output:
	// ---
	// lambda:  0
	// train set:  1
	// benchmark set:  0.93
	// ---
	// lambda:  0.1
	// train set:  0.94
	// benchmark set:  0.93
	// ---
	// lambda:  0.2
	// train set:  0.94
	// benchmark set:  0.925
	// ---
	// lambda:  0.4
	// train set:  0.925
	// benchmark set:  0.915
	// ---
	// lambda:  0.8
	// train set:  0.875
	// benchmark set:  0.86
	// ---
	// lambda:  1
	// train set:  0.905
	// benchmark set:  0.905
}

func ExamplePGM7_LogisticRegression_OptimizeLambda() {
	rand.Seed(98765)
	trainFeatures, trainLabels := readTrainData()
	benchmarkFeatures, benchmarkLabels := readBenchmarkData()

	goalFunc := func(lambda float64) float64 {
		score := ai.HoldoutTestBinaryClassifier(
			trainFeatures,
			trainLabels,
			.1,
			ai.NewLogisticRegressionTrainer(
				lambda,
				&sgrad.NumIterationsCrit{NumIterations: 10},
				1e-8))
		return -score
	}

	lambda := gcsearh.Minimize(0, 10, goalFunc, &gcsearh.AbsoluteErrorTermCrit{}, .1)
	fmt.Println("Optimal lambda:", lambda)
	classifier := ai.TrainLogisticRegressionClassifier(
		trainFeatures,
		trainLabels,
		lambda,
		&sgrad.RelativeMeanImprovementCrit{},
		1e-2)
	fmt.Println("train set: ", ai.EvaluateBinaryClassifier(classifier, trainFeatures, trainLabels))
	fmt.Println("benchmark set: ", ai.EvaluateBinaryClassifier(classifier, benchmarkFeatures, benchmarkLabels))

	// Output:
	// Optimal lambda: 1.4841053312063623
	// train set:  0.895
	// benchmark set:  0.9
}
