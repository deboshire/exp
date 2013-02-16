// This code is a playground for exersises 2.8 from the book
// Probabilistic Robotics by Sebastian Thrun, Wolfram Burgard, Dieter Fox
// ISBN-13: 978-0-262-20162-9
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type MarkovChain struct {
	States int
	Trans  []float64
}

func (ch *MarkovChain) Next(state int) int {
	r := rand.Float64()
	var cur float64
	st := ch.States * state
	for i, v := range ch.Trans[st : st+ch.States] {
		cur += v
		if r <= cur {
			return i
		}
	}
	return ch.States - 1
}

func main() {
	rand.Seed(time.Now().UnixNano())
	weather := []string{"Sunny", "Cloudy", "Rainy"}
	day2day := []float64{
		/* Sunny */ 0.8, 0.2, 0,
		/* Cloudy */ 0.4, 0.4, 0.2,
		/* Rainy */ 0.2, 0.6, 0.2,
	}
	chain := &MarkovChain{
		States: len(weather),
		Trans:  day2day,
	}
	cur := 0
	for i := 0; i < 20; i++ {
		fmt.Printf("%s\n", weather[cur])
		cur = chain.Next(cur)
	}
	fmt.Printf("\n")
}
