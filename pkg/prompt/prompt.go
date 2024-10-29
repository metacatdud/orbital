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

func Bold(colorName ColorName, msg string, args ...interface{}) {
	c := newColor(colorName)
	line := c.Add(color.Bold)
	_, _ = line.Printf(msg, args...)
}

func Info(msg string, args ...interface{}) {
	white := newColor(ColorWhite)
	_, _ = white.Printf(msg, args...)
}

func OK(msg string, args ...interface{}) {
	green := newColor(ColorGreen)
	_, _ = green.Printf(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	yellow := newColor(ColorYellow)
	_, _ = yellow.Printf(msg, args...)
}

func Err(msg string, args ...interface{}) {
	red := newColor(ColorRed)
	_, _ = red.Printf(msg, args...)
}

func NewLine(s string) string {
	return fmt.Sprintf("\n%s", s)
}

func NewLineWithTab(s string) string {
	return fmt.Sprintf("\n\t%s", s)
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
