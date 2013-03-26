package data

import (
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

	instances := Of(d)

	if instances.Len() != 14 {
		t.Errorf("Bad length: %d", instances.Len())
	}
}
