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
	Trans() []float64
}

type markovChain struct {
	states int
	trans  []float64
}

func (ch *markovChain) States() int {
	return ch.states
}

func (ch *markovChain) Trans() []float64 {
	return ch.trans
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

func Square(n int, p []float64) []float64 {
	res := make([]float64, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			ind := i*n + j
			for k := 0; k < n; k++ {
				res[ind] += p[i*n+k] * p[k*n+j]
			}
		}
	}
	return res
}

func StationaryDistr2(ch MarkovChain, count int) []float64 {
	p := ch.Trans()
	for i := 0; i < count; i++ {
		p = Square(ch.States(), p)
	}
	return p[0:ch.States()]
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

func PrintDistr(p []float64, names []string) {
	for i, v := range names {
		fmt.Printf("%s: %f\n", v, p[i])
	}
	fmt.Println()
}

// Exercise 2.8.2 (c)
func PrintStationaryDistr(ch MarkovChain, names []string, count int) {
	PrintDistr(StationaryDistr(ch, count), names)
}

func PrintStationaryDistr2(ch MarkovChain, names []string, count int) {
	PrintDistr(StationaryDistr2(ch, count), names)
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
	PrintStationaryDistr2(ch, weather, 0)
	PrintStationaryDistr2(ch, weather, 1)
	PrintStationaryDistr2(ch, weather, 2)
	PrintStationaryDistr2(ch, weather, 3)
	PrintStationaryDistr2(ch, weather, 4)
	PrintStationaryDistr2(ch, weather, 5)
	PrintStationaryDistr2(ch, weather, 6)

}
