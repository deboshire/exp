package ai

import (
	"fmt"
	"github.com/deboshire/exp/ai/data"
	"github.com/deboshire/exp/math/opt/sgrad"
	v "github.com/deboshire/exp/math/vector"
	"math"
)

type logisticRegressionClassifier struct {
	theta        v.F64
	featureAttrs data.Attributes
	minimizer    sgrad.Minimizer
}

type LogisticRegressionTrainer struct {
	Lambda   float64
	TermCrit sgrad.TermCrit
	Eps      float64
}

func (t LogisticRegressionTrainer) Train(table data.Table, classAttr data.Attr) Classifier {
	featureAttrs := table.Attrs().Without(classAttr)
	minimizer := sgrad.Minimizer{
		F:       logisticRegressionCostFunction(table, classAttr, t.Lambda),
		Initial: v.Zeroes(len(table.Attrs()) - 1),
	}

	x := minimizer.Minimize(t.Eps, t.TermCrit)

	return &logisticRegressionClassifier{theta: x, featureAttrs: featureAttrs, minimizer: minimizer}
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// http://mathurl.com/bmfs3db
func logisticRegressionCostFunction(table data.Table, classAttr data.Attr, lambda float64) sgrad.ObjectiveFunc {
	featureAttrs := table.Attrs().Without(classAttr)
	i := table.CyclicIterator([]data.Attributes{[]data.Attr{classAttr}, featureAttrs})

	grad := v.Zeroes(len(table.Attrs()) - 1)

	f := func(x v.F64) (value float64, gradient v.F64, ok bool) {
		row, ok := i()

		if !ok {
			return 0, nil, false
		}

		gradient = grad

		label := row[0][0]
		features := row[1]

		h := sigmoid(x.DotProduct(features))
		features.CopyTo(gradient)

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
	return f
}

func (c *logisticRegressionClassifier) ClassType() data.AttrType {
	return data.TYPE_BOOL
}

func (c *logisticRegressionClassifier) String() string {
	return fmt.Sprintf("logisticRegressionClassifier{totalIter: %d}", c.minimizer.State.TotalIter)
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

func (c *logisticRegressionClassifier) Classify(row v.F64) Classification {
	h := sigmoid(c.theta.DotProduct(row))
	return &logitClassification{h: h}
}

func (c *logisticRegressionClassifier) Features() data.Attributes {
	return c.featureAttrs
}
