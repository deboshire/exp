package pgm7

import (
	"fmt"
	"github.com/deboshire/exp/ai/classifiers"
	"github.com/deboshire/exp/ai/classifiers/logit"
	"github.com/deboshire/exp/ai/data"
	"github.com/deboshire/exp/io/mat"
	"github.com/deboshire/exp/math/opt"
	"github.com/deboshire/exp/math/opt/gssearch"
	"testing"

	// "github.com/deboshire/exp/math/opt/gssearch"

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

type Score struct {
	trainScore     float64
	benchmarkScore float64
}

func verifyScore(t *testing.T, testName string, expectedScore Score, actualScore Score) {
	if actualScore.trainScore != expectedScore.trainScore {
		t.Errorf("[%s] Bad train score: expected=%v, actual=%v", testName, expectedScore.trainScore, actualScore.trainScore)
	}
	if actualScore.benchmarkScore != expectedScore.benchmarkScore {
		t.Errorf("[%s] Bad benchmark score: expected=%v, actual=%v", testName, expectedScore.benchmarkScore, actualScore.benchmarkScore)
	}
}

func TestIncreasingNumberOfIterations(t *testing.T) {
	trainData, labelAttr := readTrainData()
	benchData := readBenchmarkData()

	tests := []struct {
		iterations int
		score      Score
	}{
		{1, Score{0.815, 0.825}},
		{10, Score{0.96, 0.92}},
		{100, Score{1.0, 0.92}},
		{1000, Score{1.0, 0.92}},
	}

	for _, test := range tests {
		iterations := test.iterations
		trainer := &logit.Trainer{
			Lambda:   0,
			TermCrit: &opt.NumIterationsCrit{NumIterations: iterations},
			Eps:      1e-8}
		classifier := trainer.Train(trainData, labelAttr)

		actualIterations := classifier.(*logit.LogitClassifier).Minimizer.State.TotalIter
		if actualIterations != test.iterations {
			t.Errorf("Unexpected number of iterations: %d vs %d", actualIterations, iterations)
		}

		trainScore := classifiers.Evaluate(classifier, trainData, labelAttr)
		benchmarkScore := classifiers.Evaluate(classifier, benchData, labelAttr)
		verifyScore(t, fmt.Sprint(iterations), test.score, Score{trainScore, benchmarkScore})
	}
}

func TestDecreasingEpsilon(t *testing.T) {
	trainData, labelAttr := readTrainData()
	benchData := readBenchmarkData()

	tests := []struct {
		epsilon    float64
		iterations int
		score      Score
	}{
		{1e-1, 13, Score{0.965, 0.915}},
		{1e-2, 33, Score{0.99, 0.92}},
		{1e-3, 115, Score{1.0, 0.92}},
		{1e-4, 839, Score{1.0, 0.925}},
	}

	for _, test := range tests {
		trainer := &logit.Trainer{
			Lambda:   0,
			TermCrit: &opt.MaxRelativeChangeCrit{},
			Eps:      test.epsilon,
		}

		classifier := trainer.Train(trainData, labelAttr)
		actualIterations := classifier.(*logit.LogitClassifier).Minimizer.State.TotalIter
		if actualIterations != test.iterations {
			t.Errorf("Unexpected number of iterations: %d vs %d", actualIterations, test.iterations)
		}

		trainScore := classifiers.Evaluate(classifier, trainData, labelAttr)
		benchmarkScore := classifiers.Evaluate(classifier, benchData, labelAttr)
		verifyScore(t, fmt.Sprintf("epsilon = %v", test.epsilon), test.score, Score{trainScore, benchmarkScore})
	}
}

func TestChangingHoldoutFraction(t *testing.T) {
	trainData, labelAttr := readTrainData()

	tests := []struct {
		fraction float64
		score    float64
	}{
		{0.5, 0.9},
		{0.25, 0.74},
		{0.1, 0.95},
		{0.05, 0.9},
	}

	for _, test := range tests {
		actualScore := classifiers.HoldoutTest(
			&logit.Trainer{
				Lambda:   0,
				TermCrit: &opt.MaxRelativeChangeCrit{},
				Eps:      1e-4},
			trainData,
			labelAttr,
			test.fraction)

		if actualScore != test.score {
			t.Errorf("[fraction = %v] Bad score: expected=%v, actual=%v", test.fraction, test.score, actualScore)
		}
	}
}

func TestChangingLambda(t *testing.T) {
	trainData, labelAttr := readTrainData()
	benchData := readBenchmarkData()

	tests := []struct {
		lambda float64
		score  Score
	}{
		{0, Score{1, 0.925}},
		{0.1, Score{0.955, 0.925}},
		{0.2, Score{0.945, 0.925}},
		{0.3, Score{0.93, 0.915}},
		{0.4, Score{0.915, 0.905}},
		{0.8, Score{0.905, 0.88}},
		{1, Score{0.905, 0.88}},
	}

	for _, test := range tests {
		trainer := &logit.Trainer{
			Lambda:   test.lambda,
			TermCrit: &opt.MaxRelativeChangeCrit{},
			Eps:      1e-4,
		}

		classifier := trainer.Train(trainData, labelAttr)

		trainScore := classifiers.Evaluate(classifier, trainData, labelAttr)
		benchmarkScore := classifiers.Evaluate(classifier, benchData, labelAttr)
		verifyScore(t, fmt.Sprintf("lambda = %v", test.lambda), test.score, Score{trainScore, benchmarkScore})
	}
}

// TODO: convert into test
func ExamplePGM7_LogisticRegression_OptimizeLambda() {
	trainData, labelAttr := readTrainData()
	benchmarkData := readBenchmarkData()

	goalFunc := func(lambda float64) float64 {
		score := classifiers.HoldoutTest(
			&logit.Trainer{
				Lambda:   lambda,
				TermCrit: &opt.NumIterationsCrit{NumIterations: 10},
				Eps:      1e-8},
			trainData,
			labelAttr,
			.1,
		)
		return -score
	}

	lambda := gssearh.Minimize(0, 1, goalFunc, &gssearh.AbsoluteErrorTermCrit{}, .1)
	fmt.Println("Optimal lambda:", lambda)
	classifier := (&logit.Trainer{
		Lambda:   lambda,
		TermCrit: &opt.NumIterationsCrit{NumIterations: 10 * len(trainData.Attrs())},
		Eps:      1e-8}).Train(trainData, labelAttr)
	fmt.Println("train set: ", classifiers.Evaluate(classifier, trainData, labelAttr))
	fmt.Println("benchmark set: ", classifiers.Evaluate(classifier, benchmarkData, labelAttr))

	// Output:
	// Optimal lambda: 0.24013328687768137
	// train set:  0.945
	// benchmark set:  0.925
}

func init() {
	rand.Seed(98765)
}
