package main

import (
	"flag"
	"fmt"
	"github.com/deboshire/exp/ai"
	"github.com/deboshire/exp/ai/classifiers"
	"github.com/deboshire/exp/ai/classifiers/knn"
	"github.com/deboshire/exp/ai/classifiers/logit"
	"github.com/deboshire/exp/ai/data"
	"github.com/deboshire/exp/io/csv"
	"github.com/deboshire/exp/math/opt"
	"os"
	"runtime/pprof"
)

var trainCsvPath = flag.String("train-csv", "", "Path to train.csv file from kaggle")

func benchmarkClassifier(trainer ai.ClassifierTrainer, table data.Table, labelAttr data.Attr) {
	fmt.Print("Benchmarking ", trainer.Name(), "...")

	testingFraction := 0.2
	result := classifiers.HoldoutTest(trainer, table, labelAttr, testingFraction)
	fmt.Println(" testingFraction =", testingFraction, " precision =", result)
}

func readTrainData() data.Table {
	if *trainCsvPath == "" {
		panic("--train-csv not set.")
	}
	fmt.Print("Reading training data...")

	csvData, err := csv.ReadFile(*trainCsvPath, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(csvData.Len(), "rows")

	labelAttr := csvData.Attrs().ByName("label")
	csvData.TransformAttr(labelAttr, data.TO_NOMINAL)
	return csvData
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

	csvData := readTrainData()
	labelAttr := csvData.Attrs().ByName("label")
	// TODO: add bias
	benchmarkClassifier(logit.Trainer{TermCrit: &opt.NumIterationsCrit{NumIterations: 100}}, csvData, labelAttr)
	benchmarkClassifier(knn.Trainer{K: 3}, csvData, labelAttr)
}
