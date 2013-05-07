package csv

import (
	"bufio"
	"encoding/csv"
	"github.com/deboshire/exp/ai/data"
	"github.com/deboshire/exp/math/vector"
	"os"
)

func ReadFile(fileName string, hasHeader bool) (table data.Table, err error) {
	if !hasHeader {
		panic("not implemented")
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(bufio.NewReader(file))
	csvRows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	header := csvRows[0]
	csvRows = csvRows[1:] // remove header

	rows := make([]vector.F64, len(csvRows))

	for i, csvRow := range csvRows {
		rows[i], err = vector.Parse(csvRow)
	}

	attrs := make([]data.Attr, len(header))
	for i, name := range header {
		attrs[i] = data.Attr{name, data.TYPE_NUMERIC}
	}

	return data.FromRows(rows, attrs), nil
}
