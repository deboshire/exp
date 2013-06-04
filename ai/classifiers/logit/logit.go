// logistic regression classifier
package logit

import (
	"fmt"
	"github.com/deboshire/exp/ai"
	"github.com/deboshire/exp/ai/data"
	"github.com/deboshire/exp/math/opt/sgd"
	"github.com/deboshire/exp/math/vector"
	"math"
)

type LogitClassifier struct {
	Theta        vector.F64
	FeatureAttrs data.Attributes
	Minimizer    sgd.PushMinimizer
}

type Trainer struct {
	Lambda   float64
	TermCrit sgd.TermCrit
	Eps      float64
}

func (t Trainer) Train(table data.Table, classAttr data.Attr) ai.Classifier {
	featureAttrs := table.Attrs().Without(classAttr)

	minimizer := sgd.PushMinimizer{
		X0: vector.Zeroes(len(table.Attrs()) - 1),
	}

	x := minimizer.Minimize(t.Eps, t.TermCrit, func(pushFn sgd.PushFunction) {
		// TODO: randomize order
		preallocatedGradient := vector.Zeroes(len(table.Attrs()) - 1)
		table.Shuffled().Do(func(row []vector.F64) {
			pushFn(logitFn(minimizer.State.X, row, t.Lambda, preallocatedGradient))
		}, []data.Attributes{[]data.Attr{classAttr}, featureAttrs})
	})

	return &LogitClassifier{Theta: x, FeatureAttrs: featureAttrs, Minimizer: minimizer}
}

func (t Trainer) Name() string {
	return "Logistic Regression"
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// http://mathurl.com/bmfs3db
func logitFn(x vector.F64, row []vector.F64, lambda float64, preallocatedGradient vector.F64) (value float64, gradient vector.F64) {
	gradient = preallocatedGradient

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

func (c *LogitClassifier) ClassType() data.AttrType {
	return data.TYPE_BOOL
}

func (c *LogitClassifier) String() string {
	return fmt.Sprintf("LogitClassifier{totalIter: %d}", c.Minimizer.State.TotalIter)
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

func (c *LogitClassifier) Classify(row vector.F64) ai.Classification {
	h := sigmoid(c.Theta.DotProduct(row))
	return &logitClassification{h: h}
}

func (c *LogitClassifier) Features() data.Attributes {
	return c.FeatureAttrs
}
