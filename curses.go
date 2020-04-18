package main

import (
	"fmt"
)

const ESC = "\x1B"
const ForegroundReset = ESC + "[39m"
const BackgroundReset = ESC + "[49m"

func Foreground(color string) string {
	return fmt.Sprintf("%s[38;5;%sm", ESC, color)
}

func Background(color string) string {
	return fmt.Sprintf("%s[48;5;%sm", ESC, color)
}

func SetCursor(line, col int) {
	fmt.Printf("%s[%d;%dH", ESC, line, col)
}

func CursorUp(value int) {
	fmt.Printf("%s[%dA", ESC, value)
}

func CursorDown(value int) {
	fmt.Printf("%s[%dB", ESC, value)
}

func CursorForward(value int) {
	fmt.Printf("%s[%dC", ESC, value)
}

func CursorBackward(value int) {
	fmt.Printf("%s[%dD", ESC, value)
}
