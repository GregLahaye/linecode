package main

import (
	"fmt"
	"time"
)

type Spinner struct {
	Style     Style
	Condition bool
	Channel   chan bool
}

type Style struct {
	Milliseconds time.Duration
	Characters   []string
}

func Start(style Style) *Spinner {
	s := &Spinner{
		Style:     style,
		Condition: false,
		Channel:   make(chan bool),
	}
	go s.Spin()
	return s
}

func (s *Spinner) Spin() {
	for !s.Condition {
		for _, c := range s.Style.Characters {
			fmt.Print(c)
			time.Sleep(time.Millisecond * s.Style.Milliseconds)
			CursorBackward(len(c))
		}
	}

	s.Channel <- true
}

func (s *Spinner) End() {
	s.Condition = true
	<-s.Channel
}

var Simple = Style{200, []string{".  ", ".. ", "...", " ..", "  .", "   "}}
var Star = Style{80, []string{"+", "x", "*"}}
var Point = Style{125, []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"}}
var Balloon = Style{140, []string{" ", ".", "o", "O", "@", "*", " "}}
var Bounce = Style{80, []string{"( ●    )", "(  ●   )", "(   ●  )", "(    ● )", "(     ●)", "(    ● )", "(   ●  )", "(  ●   )", "( ●    )", "(●     )"}}
var Bar = Style{120, []string{"[    ]", "[=   ]", "[==  ]", "[=== ]", "[ ===]", "[  ==]", "[   =]", "[    ]", "[   =]", "[  ==]", "[ ===]", "[====]", "[=== ]", "[==  ]", "[=   ]"}}
var Box = Style{100, []string{"▌", "▀", "▐", "▄"}}
var Noise = Style{100, []string{"▓", "▒", "░"}}
