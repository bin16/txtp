package txtp

import (
	"fmt"
	"strings"
	"unicode"
)

func parseHexColor(s string) (r, g, b, a uint8) {
	s = strings.TrimPrefix(strings.ToLower(s), "#")
	a = 255
	if len(s) == 3 {
		format := "%1x%1x%1x"
		fmt.Sscanf(s, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
		return
	}

	if len(s) == 6 {
		format := "%02x%02x%02x"
		fmt.Sscanf(s, format, &r, &g, &b)
		return
	}

	if len(s) == 8 {
		format := "%02x%02x%02x%02x"
		fmt.Sscanf(s, format, &r, &g, &b, &a)
		return
	}

	return
}

func breakWords(s string) []string {
	s = strings.TrimSpace(s)
	results := []string{}

	word := []rune{}
	for i, r := range s {
		// Space
		if i > 0 && r != ' ' && s[i-1] == ' ' {
			// push word into results & clear word
			results = append(results, string(word))
			word = []rune{r}
			continue
		}

		// Han
		if unicode.In(r, unicode.Han) {
			// break left
			if len(word) > 0 {
				// push 字 into results & clear word
				results = append(results, string(word))
				word = []rune{}
			}
			word = append(word, r)
			continue
		}

		if unicode.In(r, unicode.Latin) && len(word) > 0 && unicode.In(word[len(word)-1], unicode.Han) {
			// push 字 into results & clear word
			results = append(results, string(word))
			word = []rune{}

			word = append(word, r)
			continue
		}

		word = append(word, r)
	}

	if len(word) > 0 {
		results = append(results, string(word))
	}

	return results
}

func (ctx *Context) textWrap(s string, maxWidth float64) []string {
	parts := breakWords(s)
	lines := []string{}
	line := ""
	for _, u := range parts {
		// will not appear
		if u == "\n" {
			lines = append(lines, line)
			lines = append(lines, "")
			line = ""
			continue
		}

		if w, _ := ctx.MeasureString(line + u); w >= maxWidth {
			lines = append(lines, line)
			line = u
			continue
		}

		line += u
	}

	if len(line) > 0 {
		lines = append(lines, line)
	}

	return lines
}
