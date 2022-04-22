package txtp

import (
	"strings"
	"testing"
)

func TestBreakWords(t *testing.T) {
	d := map[string][]string{
		"Hello World!":        {"Hello ", "World!"},
		"吃了吗，世界？":             {"吃", "了", "吗，", "世", "界？"},
		"中English中English":    {"中", "English", "中", "English"},
		"中 English 中 English": {"中 ", "English ", "中 ", "English"},
		"中English 中English":   {"中", "English ", "中", "English"},
		"中English中 English":   {"中", "English", "中 ", "English"},
	}

	for s, r := range d {
		r1 := breakWords(s)
		if len(r1) != len(r) || strings.Join(r1, "++") != strings.Join(r, "++") {
			t.Errorf("breakWords() %s \n=> %s \nwant %s", s, r1, r)
		}
	}
}
