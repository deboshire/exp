// This utility takes ground truth from MIT+CMU frontal face test set
// and converts to json file.
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	X float64
	Y float64
}
type Face struct {
	LeftEye          Position
	RightEye         Position
	Nose             Position
	LeftCornerMouth  Position
	CenterMouth      Position
	RightCornerMouth Position
}

type File struct {
	Filename string
	Face     []*Face
}

type Set struct {
	File []*File
}

func NewFace(lineno int, cols []string) (face *Face) {
	face = new(Face)
	parse := func() float64 {
		val, err := strconv.ParseFloat(cols[0], 64)
		if err != nil {
			log.Printf("Line #%d: %v", lineno, err)
		}
		cols = cols[1:]
		return val
	}
	pos := func() Position { return Position{parse(), parse()} }

	face.LeftEye = pos()
	face.RightEye = pos()
	face.Nose = pos()
	face.LeftCornerMouth = pos()
	face.CenterMouth = pos()
	face.RightCornerMouth = pos()

	return
}

func main() {
	data, _ := ioutil.ReadAll(os.Stdin)
	text := string(data)
	files := make(map[string]*File)
	set := new(Set)
	for lineno, line := range strings.Split(text, "\n") {
		lineno++ // Editors start counting lines from 1
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cols := strings.Split(line, " ")
		if len(cols) != 13 {
			log.Fatalf("Line %d: unexpected number of columns. Want 13, got %d", lineno, len(cols))
		}
		filename := cols[0]
		file, ok := files[filename]
		if !ok {
			file = new(File)
			files[filename] = file
			file.Filename = filename
			set.File = append(set.File, file)
		}
		file.Face = append(file.Face, NewFace(lineno, cols[1:]))
	}

	b, err := json.MarshalIndent(set, "", "  ")
	if err != nil {
		log.Fatal("json.Marshal: ", err)
	}
	os.Stdout.Write(b)
}
