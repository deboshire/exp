// k-nearest neighbors classifier
package knn

import (
	"container/heap"
	"fmt"
	"github.com/deboshire/exp/ai"
	"github.com/deboshire/exp/ai/data"
	"github.com/deboshire/exp/math/vector"
)

type Trainer struct {
	K int
}

type knnClassifier struct {
	k              int
	data           data.Table
	classAttribute data.Attr
}

func (t Trainer) Train(table data.Table, classAttribute data.Attr) ai.Classifier {
	if classAttribute.Type.Kind != data.KIND_NOMINAL {
		panic(fmt.Errorf("Class attribute should be of nominal type: %q", classAttribute))
	}

	k := t.K
	if k == 0 {
		k = 3
	}
	return &knnClassifier{k: k, data: table, classAttribute: classAttribute}
}

func (t Trainer) Name() string {
	return "K-nearest neighbors"
}

func (c *knnClassifier) ClassType() data.AttrType {
	return c.classAttribute.Type
}

func (c *knnClassifier) Features() data.Attributes {
	return c.data.Attrs().Without(c.classAttribute)
}

type distance struct {
	dist  float64
	class float64
}

type distanceQueue []*distance

func (q *distanceQueue) Len() int { return len(*q) }

func (q *distanceQueue) Less(i, j int) bool {
	// pop needs to return maximum value, use >
	return (*q)[i].dist > (*q)[j].dist
}

func (q *distanceQueue) Pop() interface{} {
	n := q.Len()
	result := (*q)[n-1]
	*q = (*q)[0 : n-1]
	return result
}

func (q *distanceQueue) Push(x interface{}) {
	a := *q
	n := len(a)
	a = a[0 : n+1]
	item := x.(*distance)
	a[n] = item
	*q = a
}

func (q *distanceQueue) Swap(i, j int) {
	a := *q
	a[i], a[j] = a[j], a[i]
}

func (q *distanceQueue) maxDist() float64 {
	a := *q
	if len(a) == 0 {
		panic("empty queue")
	}
	return a[0].dist
}

type knnClassification struct {
	class int
}

// Classify a single data row.
func (c *knnClassifier) Classify(instance vector.F64) ai.Classification {
	queue := make(distanceQueue, 0, c.k)
	heap.Init(&queue)

	iterator := c.data.Iterator([]data.Attributes{
		[]data.Attr{c.classAttribute},
		c.data.Attrs().Without(c.classAttribute)})

	for {
		row, ok := iterator()
		if !ok {
			break
		}

		class := row[0][0]

		d := row[1].Dist2(instance)

		if queue.Len() < c.k {
			heap.Push(&queue, &distance{dist: d, class: class})
		} else {
			//fmt.Println("d", d, "len", queue.Len(), "q", queue[0], queue[1], queue[2], "maxDist", queue.maxDist())
			if queue.maxDist() > d {
				heap.Pop(&queue)
				heap.Push(&queue, &distance{dist: d, class: class})
			}
		}
	}

	votes := make([]int, c.classAttribute.Type.NumValues)

	for _, d := range queue {
		votes[int(d.class)]++
	}

	maxVote := -1
	maxVoteIdx := -1

	for i, vote := range votes {
		if vote > maxVote {
			maxVote = vote
			maxVoteIdx = i
		}
	}

	return &knnClassification{class: maxVoteIdx}
}

func (c *knnClassification) MostLikelyClass() (class float64, probability float64) {
	return float64(c.class), -1 // todo: calculate probability
}
