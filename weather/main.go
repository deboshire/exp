// This code is a playground for exersises 2.8 from the book
// Probabilistic Robotics by Sebastian Thrun, Wolfram Burgard, Dieter Fox
// ISBN-13: 978-0-262-20162-9
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type MarkovChain interface {
	States() int
	Next(state int) int
}

type markovChain struct {
	states int
	trans  []float64
}

func (ch *markovChain) States() int {
	return ch.states
}

func (ch *markovChain) Next(state int) int {
	r := rand.Float64()
	var cur float64
	st := ch.states * state
	for i, v := range ch.trans[st : st+ch.states] {
		cur += v
		if r <= cur {
			return i
		}
	}
	return ch.states - 1
}

func StationaryDistr(ch MarkovChain, count int) []float64 {
	res := make([]float64, ch.States())
	cur := 0
	for i := 0; i < count; i++ {
		cur = ch.Next(cur)
		res[cur]++
	}
	for i := 0; i < ch.States(); i++ {
		res[i] /= float64(count)
	}
	return res
}

// Exercise 2.8.2 (b)
func Simulate(ch MarkovChain, names []string, count int) {
	cur := 0
	for i := 0; i < 20; i++ {
		fmt.Printf("%s\n", names[cur])
		cur = ch.Next(cur)
	}
	fmt.Printf("\n")
}

// Exercise 2.8.2 (c)
func PrintStationaryDistr(ch MarkovChain, names []string, count int) {
	distr := StationaryDistr(ch, count)
	for i, v := range names {
		fmt.Printf("%s: %f\n", v, distr[i])
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	weather := []string{"Sunny", "Cloudy", "Rainy"}
	day2day := []float64{
		/* Sunny */ 0.8, 0.2, 0,
		/* Cloudy */ 0.4, 0.4, 0.2,
		/* Rainy */ 0.2, 0.6, 0.2,
	}
	ch := &markovChain{
		states: len(weather),
		trans:  day2day,
	}
	Simulate(ch, weather, 20)
	PrintStationaryDistr(ch, weather, 100000000)
}
