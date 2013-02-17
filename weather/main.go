// This code is a playground for exersises 2.8 from the book
// Probabilistic Robotics by Sebastian Thrun, Wolfram Burgard, Dieter Fox
// ISBN-13: 978-0-262-20162-9
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type MarkovChain interface {
	States() int
	Next(state int) int
	Trans() []float64
	P(a, b int) float64
}

type MeasurementModel interface {
	StatesA() int
	StatesB() int
	Next(a int) int
	P(a, b int) float64
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

func (ch *markovChain) P(a, b int) float64 {
	return ch.trans[a*ch.states+b]
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

type measurementModel struct {
	statesA int
	statesB int
	trans   []float64
}

func (m *measurementModel) StatesA() int { return m.statesA }
func (m *measurementModel) StatesB() int { return m.statesB }

func (m *measurementModel) Next(a int) int {
	r := rand.Float64()
	var cur float64
	st := a * m.statesB
	for i, v := range m.trans[st : st+m.statesB] {
		cur += v
		if r <= cur {
			return i
		}
	}
	return m.statesB - 1
}

func (m *measurementModel) P(a, b int) float64 {
	return m.trans[a*m.statesB+b]
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
	fmt.Print("[ ")
	for i, v := range names {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%s: %f", v, p[i])
	}
	fmt.Print(" ] \n")
}

// Exercise 2.8.2 (c)
func PrintStationaryDistr(ch MarkovChain, names []string, count int) {
	PrintDistr(StationaryDistr(ch, count), names)
}

func PrintStationaryDistr2(ch MarkovChain, names []string, count int) {
	PrintDistr(StationaryDistr2(ch, count), names)
}

func PrintEntropy(p []float64) {
	var res float64
	for _, v := range p {
		res -= v * math.Log2(v)
	}
	fmt.Printf("Entropy: %f\n", res)
}

// Exercise 2.8.2 (e)
func PrintStationaryEntropy(ch MarkovChain) {
	PrintEntropy(StationaryDistr2(ch, 10))
}

func Predict(ch MarkovChain, bel []float64) []float64 {
	res := make([]float64, ch.States())
	for i := 0; i < ch.States(); i++ {
		for j := 0; j < ch.States(); j++ {
			res[j] += bel[i] * ch.P(i, j)
		}
	}
	return res
}

// Simulate Hidden Markov Model using Bayes Filter
func SimulateMeasurements(ch MarkovChain, mm MeasurementModel, names []string, count int) {
	cur := 0 // Sunny
	bel := []float64{1, 0, 0}
	fmt.Printf("Initially, weather is %s\n", names[cur])
	for i := 1; i <= count; i++ {
		predict := Predict(ch, bel)
		cur = ch.Next(cur)
		curm := mm.Next(cur)
		bel = []float64{0, 0, 0}
		var sum float64
		for j := 0; j < len(bel); j++ {
			bel[j] = mm.P(j, curm) * predict[j]
			sum += bel[j]
		}
		// Normalize belief
		for j := 0; j < len(bel); j++ {
			bel[j] = bel[j] / sum
		}

		fmt.Printf("Step %d, weather: %s, we see: %s, belief: ", i, names[cur], names[curm])
		PrintDistr(bel, names)
		fmt.Println()
	}
	fmt.Println()
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

	real2view := []float64{
		/* Sunny */ 0.6, 0.4, 0,
		/* Cloudy */ 0.3, 0.7, 0,
		/* Rainy */ 0, 0, 1,
	}

	mm := &measurementModel{
		statesA: len(weather),
		statesB: len(weather),
		trans:   real2view,
	}

	Simulate(ch, weather, 20)
	PrintStationaryDistr(ch, weather, 100000000)
	PrintStationaryDistr2(ch, weather, 10)
	PrintStationaryEntropy(ch)

	SimulateMeasurements(ch, mm, weather, 30)
}
