package main

import (
	"fmt"
)

const ESC = "\x1B"

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

func Foreground(color string) string {
	return fmt.Sprintf("%s[38;5;%sm", ESC, color)
}

func ForegroundReset() string {
	return fmt.Sprintf("%s[39m", ESC)
}

func Background(color string) string {
	return fmt.Sprintf("%s[48;5;%sm", ESC, color)
}

func BackgroundReset() string {
	return fmt.Sprintf("%s[49m", ESC)
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
