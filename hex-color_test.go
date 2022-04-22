package txtp

import (
	"image/color"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	d := map[string]color.NRGBA{
		"#000000": {0, 0, 0, 255},
		"#ffffff": {255, 255, 255, 255},
	}

	for s, c := range d {
		r, g, b, a := parseHexColor(s)
		c1 := color.NRGBA{r, g, b, a}
		if c1 != c {
			t.Errorf("parseHexColor() failed, %s => %v, want %v\n", s, c1, c)
		}
	}
}
