package convert

import (
	"golang.org/x/net/html"
	"strconv"
	"strings"
)

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func PadString(s string, max int, left bool) string {
	length := len(s)
	if length > max {
		return s
	}

	difference := max - length
	padding := strings.Repeat(" ", difference)
	if left {
		s = padding + s
	} else {
		s += padding
	}

	return s
}

func ParseHTML(h string) string {
	z := html.NewTokenizer(strings.NewReader(h))

	s := ""
	for {
		tt := z.Next()
		t := z.Token()

		switch tt {
		case html.ErrorToken:
			return s
		case html.TextToken:
			s += t.Data
		}
	}
}
