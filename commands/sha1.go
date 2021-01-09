package commands

import (
	"crypto/sha1"
	"encoding/hex"

	"../command"
)

type Sha1Command struct {
	Arg string
}

func (s *Sha1Command) Execute(loop command.Handler) {
	h := sha1.New()
	h.Write([]byte(s.Arg))
	sha1 := hex.EncodeToString(h.Sum(nil))
	loop.Post(&PrintCommand{Arg: sha1})
}
