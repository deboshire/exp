package classifiers

import (
	"github.com/deboshire/exp/ai"
	"github.com/deboshire/exp/ai/data"
)

// Evaluate binary classifier on a given data.
// Returns percentage of correct hits.
func Evaluate(c ai.Classifier, instances data.Instances, classAttr data.Attr) float64 {
	if instances.Len() == 0 {
		return 0
	}

	successes := 0

	for i := 0; i < instances.Len(); i++ {
		instance := instances.Get(i)
		class, _ := c.Classify(instance).MostLikelyClass()
		if class == instance.Get(classAttr) {
			successes++
		}
	}

	return float64(successes) / float64(instances.Len())
}

func HoldoutTest(trainer ai.ClassifierTrainer, instances data.Instances, classAttr data.Attr, testingFraction float64) float64 {
	shuffledData := instances.Shuffled()

	idx := int(float64(shuffledData.Len()) * (1 - testingFraction))

	trainingData, testingData := shuffledData.Split(idx)
	classifier := trainer.Train(trainingData, classAttr)
	return Evaluate(classifier, testingData, classAttr)
}
