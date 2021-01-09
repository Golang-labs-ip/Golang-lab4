package commands

import (
	"os"

	"../command"
)

type ErrorCommand struct {
	Msg string
}

func (p *ErrorCommand) Execute(loop command.Handler) {
	os.Stderr.WriteString(p.Msg)
}
