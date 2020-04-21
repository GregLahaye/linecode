package main

import (
	"flag"
	"os"
	"strings"
)

type ArrayFlags []string

func (a *ArrayFlags) String() string {
	return strings.Join(*a, ", ")
}

func (a *ArrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func FlagStatus(positive, negative *bool) int {
	if *positive {
		return accept
	} else if *negative {
		return deny
	} else {
		return none
	}
}

func parseFilterFlags() Filter {
	var tags ArrayFlags

	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Var(&tags, "t", "tags")

	e := fs.Bool("e", false, "easy questions")
	E := fs.Bool("E", false, "not easy questions")

	m := fs.Bool("m", false, "medium questions")
	M := fs.Bool("M", false, "not medium questions")

	h := fs.Bool("h", false, "hard questions")
	H := fs.Bool("H", false, "not hard questions")

	a := fs.Bool("a", false, "accepted")
	A := fs.Bool("A", false, "not accepted")

	p := fs.Bool("p", false, "paid")
	P := fs.Bool("P", false, "not paid")

	s := fs.Bool("s", false, "starred")
	S := fs.Bool("S", false, "not starred")

	_ = fs.Parse(os.Args[2:])

	return Filter{
		Tags:     tags,
		Easy:     FlagStatus(e, E),
		Medium:   FlagStatus(m, M),
		Hard:     FlagStatus(h, H),
		Accepted: FlagStatus(a, A),
		Paid:     FlagStatus(p, P),
		Starred:  FlagStatus(s, S),
	}
}
