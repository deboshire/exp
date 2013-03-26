package ai

import (
	"github.com/deboshire/exp/ai/data"
	"github.com/deboshire/exp/math/opt/sgrad"
	v "github.com/deboshire/exp/math/vector"
	"math"
)

type logisticRegressionClassifier struct {
	cost         float64
	theta        v.F64
	featureAttrs data.Attributes
}

type LogisticRegressionTrainer struct {
	Lambda   float64
	TermCrit sgrad.TermCrit
	Eps      float64
}

func (t LogisticRegressionTrainer) Train(instances data.Instances, classAttr data.Attr) Classifier {
	featureAttrs := instances.Attrs().Without(classAttr)

	y, x := sgrad.Minimize(
		logisticRegressionCostFunction(instances, classAttr, t.Lambda),
		v.Zeroes(len(instances.Attrs())-1),
		t.Eps,
		t.TermCrit,
		nil)

	return &logisticRegressionClassifier{cost: y, theta: x, featureAttrs: featureAttrs}
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// http://mathurl.com/bmfs3db
func logisticRegressionCostFunction(instances data.Instances, classAttr data.Attr, lambda float64) sgrad.ObjectiveFunc {
	featureAttrs := instances.Attrs().Without(classAttr)

	features := instances.View(featureAttrs)
	labels := instances.View([]data.Attr{classAttr})

	f := func(idx int, x v.F64, gradient v.F64) (value float64) {
		feature := features[idx]
		label := labels[idx][0]

		h := sigmoid(x.DotProduct(feature))
		feature.CopyTo(gradient)

		if label != 0 {
			value = -math.Log(h)
			gradient.Mul(h - 1.0)
		} else {
			value = -math.Log(1.0 - h)
			gradient.Mul(h)
		}

		if lambda != 0.0 {
			// apply regularizaiton.
			for i := 1; i < len(x); i++ {
				value += 0.5 * lambda * x[i] * x[i]
				gradient[i] += lambda * x[i]
			}
		}

		return
	}

	return sgrad.ObjectiveFunc{Terms: instances.Len(), F: f}
}

func (c *logisticRegressionClassifier) ClassType() data.AttrType {
	return data.TYPE_BOOL
}

type logitClassification struct {
	h float64
}

func (c *logitClassification) MostLikelyClass() (class float64, probability float64) {
	class = 0
	if c.h >= 0.5 {
		class = 1
	}
	probability = math.Abs(0.5-c.h) * 2.0
	return
}

func (c *logisticRegressionClassifier) Classify(instance data.Instance) Classification {
	features := instance.View(c.featureAttrs)
	h := sigmoid(c.theta.DotProduct(features))
	return &logitClassification{h: h}
}
