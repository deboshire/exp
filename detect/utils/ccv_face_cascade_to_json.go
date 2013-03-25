// Pretty prints json with cascades.
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Feature struct {
	Size int   `json:"size"`
	PX   []int `json:"px"`
	PY   []int `json:"py"`
	PZ   []int `json:"pz"`
	NX   []int `json:"nx"`
	NY   []int `json:"ny"`
	NZ   []int `json:"nz"`
}

type StageClassifier struct {
	Count     int       `json:"count"`
	Threshold float64   `json:"threshold"`
	Feature   []Feature `json:"feature"`
	Alpha     []float64 `json:"alpha"`
}

type Cascade struct {
	Length          int               `json:"length"`
	Width           int               `json:"width"`
	Height          int               `json:"height"`
	StageClassifier []StageClassifier `json:"stage_classifier"`
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Could not read stdin: ", err)
	}
	var cascade *Cascade
	if err := json.Unmarshal(data, &cascade); err != nil {
		log.Fatal("json.Unmarshal: ", err)
	}
	if out, err := json.MarshalIndent(cascade, "", "  "); err == nil {
		os.Stdout.Write(out)
	} else {
		log.Fatal("json.Masrshal: ", err)
	}
}
