package main

import (
	"fmt"
)

const ESC = "\x1B"
const DEFAULT = "9"

const (
	Black   = "0"
	Red     = "1"
	Green   = "2"
	Yellow  = "3"
	Blue    = "4"
	Magenta = "5"
	Cyan    = "6"
	White   = "7"
)

func SetBold() {
	fmt.Printf("%s[1m", ESC)
}

func SetUnderline() {
	fmt.Printf("%s[4m", ESC)
}

func SetBlink() {
	fmt.Printf("%s[5m", ESC)
}

func Bright(color string) string {
	return fmt.Sprintf("%s;1", color)
}

func TwoFiftySix(color string) string {
	return fmt.Sprintf("8;5;%s", color)
}

func Foreground(color string) string {
	return fmt.Sprintf("%s[3%sm", ESC, color)
}

func Background(color string) string {
	return fmt.Sprintf("%s[4%sm", ESC, color)
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

func ClearScreen() {
	fmt.Printf("%s[2J", ESC)
}

func ClearLine() {
	fmt.Printf("%s[K", ESC)
}

func ResetCursor() {
	fmt.Printf("%s[00H", ESC)
}
