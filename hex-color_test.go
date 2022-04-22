package txtp

import (
	"image/color"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	d := map[string]color.NRGBA{
		"#000000":   {0, 0, 0, 255},
		"#ffffff":   {255, 255, 255, 255},
		"ff0000":    {255, 0, 0, 255},
		"ff000000":  {255, 0, 0, 0},
		"#627318":   {98, 115, 24, 255},
		"318950ff":  {49, 137, 80, 255},
		"fff":       {255, 255, 255, 255},
		"#c00":      {204, 0, 0, 255},
		"#ffff8080": {255, 255, 128, 128},
	}

	for s, c := range d {
		r, g, b, a := parseHexColor(s)
		c1 := color.NRGBA{r, g, b, a}
		if c1 != c {
			t.Errorf("parseHexColor() failed, %s => %v, want %v\n", s, c1, c)
		}
	}
}
