package ui

import (
	"os"

	"github.com/fatih/color"
)

func Printf(format string, args ...interface{}) {
	color.New(color.FgHiBlue).Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	color.New(color.FgGreen).Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	color.New(color.FgYellow).Printf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	color.New(color.FgRed).Fprintf(os.Stderr, format, args...)
}
