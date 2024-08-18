package terminal

import (
	"fmt"
	"strconv"
)

func (t *Terminal) PromptSelect(prompt string, options []string) (int, string, error) {
	fmt.Fprintln(t.writer, prompt)
	for i, option := range options {
		fmt.Fprintln(t.writer, "[%d] %s\n", i+1, option)
	}

	for {
		input, err := t.Prompt("Enter the number of your choice:")
		if err != nil {
			return 0, "", err
		}

		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(options) {
			fmt.Fprintln(t.writer, "Invalid input. Please try again.")
			continue
		}

		return index - 1, options[index-1], nil

	}
}
