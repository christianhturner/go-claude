package terminal

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/christianhturner/go-claude/logger"
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

func GetWidthAndHeight() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		logger.FatalError(err, "Error occured at GetWidthAndHeight; out, err: err returned a value")
	}

	s := string(out)
	s = strings.TrimSpace(s)
	sArr := strings.Split(s, " ")

	height, err := strconv.Atoi(sArr[0])
	if err != nil {
		logger.FatalError(err, "Error occured at GetWidthAndHeight; height, err: returned a value")
	}

	width, err := strconv.Atoi(sArr[1])
	if err != nil {
		logger.FatalError(err, "Error occurred at GetWidthAndHeight; width, err: returned a value")
	}
	return width, height
}
