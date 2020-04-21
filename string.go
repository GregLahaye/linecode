package main

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

func PadString(str string, max int, left bool) string {
	length := len(str)
	if length > max {
		return str
	}

	difference := max - length
	padding := strings.Repeat(" ", difference)
	if left {
		str = padding + str
	} else {
		str += padding
	}

	return str
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
