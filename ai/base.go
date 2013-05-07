package ai

import (
	"github.com/deboshire/exp/ai/data"
	"github.com/deboshire/exp/math/vector"
)

//------------------------------------------------------------------------------

// Classifiers can perform classification in lots of different ways. They can
// determine a single class, probability distribution among the classes etc.
// Classification holds this results and allows to convert it to other values.
type Classification interface {
	// determines the most likely class
	MostLikelyClass() (class float64, probability float64)
}

type Classifier interface {
	// Describes the class attribute type that this classifier provides.
	ClassType() data.AttrType

	// List of features that classifier uses. Classify() accepts vectors strictly in this order.
	Features() data.Attributes

	// Classify a single data row.
	Classify(row vector.F64) Classification
}

type ClassifierTrainer interface {
	Name() string

	Train(table data.Table, classAttribute data.Attr) Classifier
}
