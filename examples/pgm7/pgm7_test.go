package pgm7

import (
	"fmt"
	"github.com/deboshire/exp/ai"
	"github.com/deboshire/exp/ai/classifiers"
	"github.com/deboshire/exp/ai/data"
	"github.com/deboshire/exp/io/mat"
	"github.com/deboshire/exp/math/opt/gssearch"
	"github.com/deboshire/exp/math/opt/sgrad"
	"math/rand"
)

func readTrainData() (d data.Table, classAttribute data.Attr) {
	features := mat.MustRead("Train1X.mat").Array("Train1X").RowsToTable()
	labels := mat.MustRead("Train1Y.mat").Array("Train1Y").RowsToTable()

	d = data.Zip(features, labels)
	classAttribute = d.Attrs().ByName("Train1Y.0")
	return
}

func readBenchmarkData() data.Table {
	benchmarkFeatures := mat.MustRead("Validation1X.mat").Array("Validation1X").Rename("Train1X").RowsToTable()
	benchmarkLabels := mat.MustRead("Validation1Y.mat").Array("Validation1Y").Rename("Train1Y").RowsToTable()
	return data.Zip(benchmarkFeatures, benchmarkLabels)
}

func ExamplePGM7_LogisticRegression_Iterations() {
	rand.Seed(98765)
	trainData, labelAttr := readTrainData()
	benchData := readBenchmarkData()

	for _, iterations := range []int{1, 10, 100, 1000} {
		fmt.Println("---\niterations: ", iterations)
		trainer := &ai.LogisticRegressionTrainer{
			Lambda:   0,
			TermCrit: &sgrad.NumIterationsCrit{NumIterations: iterations * trainData.Len()},
			Eps:      1e-8}
		classifier := trainer.Train(trainData, labelAttr)
		fmt.Println("classifier:", classifier)
		fmt.Println("train set: ", classifiers.Evaluate(classifier, trainData, labelAttr))
		fmt.Println("benchmark set: ", classifiers.Evaluate(classifier, benchData, labelAttr))
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


func ExamplePGM7_LogisticRegression_HoldoutTesting() {
	rand.Seed(98765)
	trainData, labelAttr := readTrainData()

	for _, testingFraction := range []float64{0.5, 0.25, 0.1, 0.05} {
		fmt.Println("---\nfraction: ", testingFraction)
		score := classifiers.HoldoutTest(
			&ai.LogisticRegressionTrainer{
				Lambda:   0,
				TermCrit: &sgrad.NumIterationsCrit{NumIterations: 10 * trainData.Len()},
				Eps:      1e-8},
			trainData,
			labelAttr,
			testingFraction)
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

func ExamplePGM7_LogisticRegression_Epsilon() {
	rand.Seed(98765)
	trainData, labelAttr := readTrainData()
	benchmarkData := readBenchmarkData()

	for _, epsilon := range []float64{1e-1, 1e-2, 1e-3} {
		fmt.Println("---\nepsilon: ", epsilon)
		trainer := &ai.LogisticRegressionTrainer{
			Lambda:   0,
			TermCrit: &sgrad.AbsDistanceCrit{},
			Eps:      epsilon,
		}
		classifier := trainer.Train(trainData, labelAttr)
		fmt.Println("train set: ", classifiers.Evaluate(classifier, trainData, labelAttr))
		fmt.Println("benchmark set: ", classifiers.Evaluate(classifier, benchmarkData, labelAttr))
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
	trainData, labelAttr := readTrainData()
	benchmarkData := readBenchmarkData()

	for _, lambda := range []float64{0, 0.1, 0.2, 0.4, 0.8, 1} {
		fmt.Println("---\nlambda: ", lambda)
		classifier := ai.LogisticRegressionTrainer{
			Lambda:   lambda,
			TermCrit: &sgrad.AbsDistanceCrit{},
			Eps:      1e-4,
		}.Train(trainData, labelAttr)
		fmt.Println("train set: ", classifiers.Evaluate(classifier, trainData, labelAttr))
		fmt.Println("benchmark set: ", classifiers.Evaluate(classifier, benchmarkData, labelAttr))
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
	trainData, labelAttr := readTrainData()
	benchmarkData := readBenchmarkData()

	goalFunc := func(lambda float64) float64 {
		score := classifiers.HoldoutTest(
			ai.LogisticRegressionTrainer{
				Lambda:   lambda,
				TermCrit: &sgrad.NumIterationsCrit{NumIterations: 10 * len(trainData.Attrs())},
				Eps:      1e-8},
			trainData,
			labelAttr,
			.1,
		)
		return -score
	}

	lambda := gssearh.Minimize(0, 10, goalFunc, &gssearh.AbsoluteErrorTermCrit{}, .1)
	fmt.Println("Optimal lambda:", lambda)
	classifier := ai.LogisticRegressionTrainer{
		Lambda:   lambda,
		TermCrit: &sgrad.NumIterationsCrit{NumIterations: 10 * len(trainData.Attrs())},
		Eps:      1e-8}.Train(trainData, labelAttr)
	fmt.Println("train set: ", classifiers.Evaluate(classifier, trainData, labelAttr))
	fmt.Println("benchmark set: ", classifiers.Evaluate(classifier, benchmarkData, labelAttr))

	// Output:
	// Optimal lambda: 1.4841053312063623
	// train set:  0.895
	// benchmark set:  0.9
}
