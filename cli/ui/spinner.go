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
	suffix  string
}

func NewSpinner(suffix string) *Spinner {
	return &Spinner{
		spinner: spinner.New(spinner.CharSets[symbol], 100*time.Millisecond),
		suffix:  suffix,
	}
}

func (s *Spinner) Start() {
	s.spinner.Writer = os.Stdout
	s.spinner.Prefix = "\r"
	if len(s.suffix) > 0 {
		s.spinner.Suffix = "  " + s.suffix
	}
	s.spinner.HideCursor = true
	s.spinner.Start()
}

func (s *Spinner) Stop() {
	s.spinner.Stop()
}
