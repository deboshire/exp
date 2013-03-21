package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/deboshire/exp/ai"
	v "github.com/deboshire/exp/math/vector"
	"github.com/deboshire/exp/optimization/sgrad"
	"os"
	"runtime/pprof"
	"strconv"
)

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
	{
		f, err := os.Create("digits.prof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
	}
	defer pprof.StopCPUProfile()

	fmt.Print("Reading training data...")
	file, err := os.Open(os.ExpandEnv("$HOME/Dropbox/Projects/kaggle/digits/train.csv"))
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

	binClassifierTrainer := func(features []v.F64, labels []bool) ai.BinaryClassifier {
		fmt.Println("Training binary classifier")
		classifier := ai.TrainLogisticRegressionClassifier(
			features,
			labels,
			0,
			&sgrad.NumIterationsCrit{NumIterations: 10},
			1e-8)
		return classifier
	}

	classifier := ai.TrainNominalClassifierFromBinary(pixels, labels, 10, binClassifierTrainer)
	fmt.Println(classifier)
}
