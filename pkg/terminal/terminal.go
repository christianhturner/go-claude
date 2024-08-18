package terminal

import (
	"bufio"
	"io"
	"os"
)

type Terminal struct {
	reader *bufio.Reader
	writer io.Writer
}

type TerminalInterface interface {
	Prompt(prompt string) (string, error)
	PromptPassword(prompt string) (string, error)
	PromptConfirm(prompt string) (bool, error)
	PromptSelect(prompt string, options []string) (int, string, error)
}

func New() *Terminal {
	return &Terminal{
		reader: bufio.NewReader(os.Stdin),
		writer: os.Stdout,
	}
}

func (t *Terminal) SetReader(reader io.Reader) {
	t.reader = bufio.NewReader(reader)
}

func (t *Terminal) SetWriter(writer io.Writer) {
	t.writer = writer
}
