package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/deboshire/exp/ai"
	"github.com/deboshire/exp/ai/classifiers"
	"github.com/deboshire/exp/math/opt/sgrad"
	v "github.com/deboshire/exp/math/vector"
	"os"
	"runtime/pprof"
	"strconv"
)

var trainCsvPath = flag.String("train-csv", "", "Path to train.csv file from kaggle")

func parseVector(strs []string) (res v.F64, err error) {
	res = v.Zeroes(len(strs))

	for i, str := range strs {
		parsedFloat, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return res, err
		}
		res[i] = parsedFloat
	}
	return
}

func main() {
	flag.Parse()

	{
		f, err := os.Create("digits.prof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
	}
	defer pprof.StopCPUProfile()

	fmt.Print("Reading training data...")
	file, err := os.Open(*trainCsvPath)
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(bufio.NewReader(file))
	allData, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	allData = allData[1:] // remove header

	fmt.Println("done:", len(allData), "rows")
	labels := make([]int, len(allData))
	pixels := make([]v.F64, len(allData))

	// TODO(mike): add bias term
	for i, row := range allData {
		parsedLabel, err := strconv.ParseInt(row[0], 10, 32)
		if err != nil {
			panic(err)
		}
		labels[i] = int(parsedLabel)
		pixels[i], err = parseVector(row[1:])
	}

	// binTrainer := func(features []v.F64, labels []bool) ai.Classifier {
	// 	fmt.Println("Training binary classifier")
	// 	classifier := ai.TrainLogisticRegressionClassifier(
	// 		features,
	// 		labels,
	// 		0,
	// 		&sgrad.NumIterationsCrit{NumIterations: 10},
	// 		1e-8)
	// 	return classifier
	// }

	binTrainer := &ai.LogisticRegressionTrainer{
		Lambda: 0,
		TermCrit: &sgrad.NumIterationsCrit{NumIterations: 10},
		Eps: 1e-8}

	classifier := classifiers.NominalClassifierTrainerFromBinary(pixels, labels, 10, binTrainer)
	fmt.Println(classifier)
}
