package ai

import (
	"github.com/deboshire/exp/ai/data"
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

	// Classify a single data row.
	Classify(row data.Row) Classification
}

type ClassifierTrainer interface {
	Train(table data.Table, classAttribute data.Attr) Classifier
}
