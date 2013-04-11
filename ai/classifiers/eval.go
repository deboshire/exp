package classifiers

import (
	"github.com/deboshire/exp/ai"
	"github.com/deboshire/exp/ai/data"
)

// Evaluate binary classifier on a given data.
// Returns percentage of correct hits.
func Evaluate(c ai.Classifier, table data.Table, classAttr data.Attr) float64 {
	if table.Len() == 0 {
		return 0
	}

	featureAttrs := c.Features()
	successes := 0

	for it := table.Iterator([]data.Attributes{[]data.Attr{classAttr}, featureAttrs}); ; {
		row, ok := it()
		if !ok {
			break
		}

		class, _ := c.Classify(row[1]).MostLikelyClass()
		if class == row[0][0] {
			successes++
		}
	}

	return float64(successes) / float64(table.Len())
}

func HoldoutTest(trainer ai.ClassifierTrainer, table data.Table, classAttr data.Attr, testingFraction float64) float64 {
	shuffledData := table.Shuffled()

	idx := int(float64(shuffledData.Len()) * (1 - testingFraction))

	trainingData, testingData := shuffledData.Split(idx)
	classifier := trainer.Train(trainingData, classAttr)
	return Evaluate(classifier, testingData, classAttr)
}
