package terminal

import (
	"fmt"
	"strings"
)

func (t *Terminal) Prompt(prompt string) (string, error) {
	fmt.Fprint(t.writer, prompt+" ")
	input, err := t.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func (t *Terminal) PromptConfirm(prompt string) (bool, error) {
	for {
		input, err := t.Prompt(prompt + " (y/n)")
		if err != nil {
			return false, err
		}
		switch strings.ToLower(input) {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
	}
}
