package commands

import (
	"fmt"

	"../command"
)

type PrintCommand struct {
	Arg string
}

func (p *PrintCommand) Execute(loop command.Handler) {
	fmt.Println(p.Arg)
}
