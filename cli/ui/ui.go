package ui

import (
	"io"
	"os"

	"github.com/fatih/color"
)

type UI interface {
	Printf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func Printf(format string, args ...interface{}) {
	defaultUI.Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultUI.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultUI.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultUI.Errorf(format, args...)
}

type CLI struct {
	Stdout io.Writer
	Stderr io.Writer
}

var defaultUI = &CLI{Stdout: os.Stdout, Stderr: os.Stderr}

func (c *CLI) Printf(format string, args ...interface{}) {
	color.New(color.FgHiBlue).Fprintf(c.Stdout, format, args...)
}

func (c *CLI) Infof(format string, args ...interface{}) {
	color.New(color.FgGreen).Fprintf(c.Stdout, format, args...)
}

func (c *CLI) Warnf(format string, args ...interface{}) {
	color.New(color.FgYellow).Fprintf(c.Stdout, format, args...)
}

func (c *CLI) Errorf(format string, args ...interface{}) {
	color.New(color.FgRed).Fprintf(c.Stderr, format, args...)
}
