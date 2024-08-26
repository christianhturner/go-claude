package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/christianhturner/go-claude/logger"
	"golang.org/x/term"
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

type Option struct {
	ID          interface{}
	Description string
}

func (t *Terminal) PromptOptionsSelect(options map[interface{}]string) Option {
	optionSlice := make([]Option, 0, len(options))
	for id, desc := range options {
		optionSlice = append(optionSlice, Option{ID: id, Description: desc})
	}

	selectedIndex := 0
	optionsCount := len(optionSlice)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		logger.PanicError(err, "Panic occurred registering terminal old state.")
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		clearScreen()
		for i, opt := range optionSlice {
			if i == selectedIndex {
				fmt.Fprintf(t.writer, "> %s\r\n", opt.Description)
			} else {
				fmt.Fprintf(t.writer, "  %s\r\n", opt.Description)
			}
		}
		fmt.Fprint(t.writer, "\r\nUse h (left), j (down), k (up), l (right) to navigate; Enter to select; q to quit\r\n")

		b := make([]byte, 3)
		os.Stdin.Read(b)

		switch {
		case b[0] == 27 && b[1] == 91: // Arrow keys
			switch b[2] {
			case 65, 68: // Up arrow or Left arrow
				selectedIndex = (selectedIndex - 1 + optionsCount) % optionsCount
			case 66, 67: // Down arrow or Right arrow
				selectedIndex = (selectedIndex + 1) % optionsCount
			}
		case b[0] == 'h' || b[0] == 'H' || b[0] == 'k' || b[0] == 'K':
			selectedIndex = (selectedIndex - 1 + optionsCount) % optionsCount
		case b[0] == 'l' || b[0] == 'L' || b[0] == 'j' || b[0] == 'J':
			selectedIndex = (selectedIndex + 1) % optionsCount
		case b[0] == 13: // Enter
			return optionSlice[selectedIndex]
		case b[0] == 'q' || b[0] == 'Q':
			return Option{} // Return empty option if user quits
		}
	}
}

func clearScreen() {
	time.Sleep(50 * time.Millisecond)
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[2J\033[H")
	}
}
