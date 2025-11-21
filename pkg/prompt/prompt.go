package prompt

import (
	"fmt"

	"github.com/fatih/color"
)

type ColorName string

const (
	ColorGreen  ColorName = "green"
	ColorYellow ColorName = "yellow"
	ColorRed    ColorName = "red"
	ColorWhite  ColorName = "white"
)

func Bold(colorName ColorName, msg string, args ...any) {
	c := newColor(colorName)
	line := c.Add(color.Bold)
	_, _ = line.Printf(msg, args...)
}

func Info(msg string, args ...any) {
	white := newColor(ColorWhite)
	_, _ = white.Printf(msg, args...)
}

func OK(msg string, args ...any) {
	green := newColor(ColorGreen)
	_, _ = green.Printf(msg, args...)
}

func Warn(msg string, args ...any) {
	yellow := newColor(ColorYellow)
	_, _ = yellow.Printf(msg, args...)
}

func Err(msg string, args ...any) {
	red := newColor(ColorRed)
	_, _ = red.Printf(msg, args...)
}

func NewLine(msg string, args ...any) string {
	return "\n" + fmt.Sprintf(msg, args...)
}

func NewLineWithTab(msg string, args ...any) string {
	return "\n\t" + fmt.Sprintf(msg, args...)
}

func newColor(c ColorName) *color.Color {
	switch c {
	case ColorGreen:
		return color.New(color.FgGreen)
	case ColorYellow:
		return color.New(color.FgYellow)
	case ColorRed:
		return color.New(color.FgRed)
	case ColorWhite:
		fallthrough
	default:
		return color.New(color.FgWhite)
	}
}
