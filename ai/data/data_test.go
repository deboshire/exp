package data

import (
	"fmt"
	"github.com/bmizerany/assert"
	"testing"
)

type WeatherData struct {
	Outlook     string // sunny/overcast/rainy
	Temperature string // hot/mild/cool
	Humidity    string // high/normal/
	Windy       bool
	Play        bool // should the game be played or not
}

func TestDataDefinition(t *testing.T) {
	d := []WeatherData{
		WeatherData{"sunny", "hot", "high", false, false},
		WeatherData{"sunny", "hot", "high", true, false},
		WeatherData{"overcast", "hot", "high", false, true},
		WeatherData{"rainy", "mild", "high", false, true},
		WeatherData{"rainy", "cool", "normal", false, true},
		WeatherData{"rainy", "cool", "normal", true, false},
		WeatherData{"overcast", "cool", "normal", true, true},
		WeatherData{"sunny", "mild", "high", false, false},
		WeatherData{"sunny", "cool", "normal", false, true},
		WeatherData{"rainy", "mild", "normal", false, true},
		WeatherData{"sunny", "mild", "normal", true, true},
		WeatherData{"overcast", "mild", "high", true, true},
		WeatherData{"overcast", "hot", "normal", false, true},
		WeatherData{"rainy", "mild", "high", true, false},
	}

	table := Of(d)
	assert.Equal(t, 14, table.Len())

	s := ""
	it := table.Iterator([]Attributes{
		[]Attr{table.Attrs().ByName("Windy")},
		[]Attr{table.Attrs().ByName("Windy"), table.Attrs().ByName("Play")},
	})
	for {
		row, ok := it()
		if !ok {
			break
		}

		s += fmt.Sprint(row)
		s += ", "
	}
	assert.Equal(t, "[[0] [0 0]], [[1] [1 0]], [[0] [0 1]], [[0] [0 1]], [[0] [0 1]], [[1] [1 0]], [[1] [1 1]], [[0] [0 0]], [[0] [0 1]], [[0] [0 1]], [[1] [1 1]], [[1] [1 1]], [[0] [0 1]], [[1] [1 0]], ", s)
}
