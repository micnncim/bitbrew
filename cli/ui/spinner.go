package ui

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
)

const (
	symbol = 14
)

type Spinner struct {
	spinner *spinner.Spinner
	text    string
}

func NewSpinner(text string) *Spinner {
	return &Spinner{
		spinner: spinner.New(spinner.CharSets[symbol], 100*time.Millisecond),
		text:    text,
	}
}

func (s *Spinner) Start() {
	s.spinner.Writer = os.Stdout
	s.spinner.Prefix = "\r"
	if len(s.text) > 0 {
		s.spinner.Suffix = "  " + s.text
	}
	s.spinner.HideCursor = true
	s.spinner.Start()
}

func (s *Spinner) Stop() {
	s.spinner.Stop()
}
