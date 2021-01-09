package command

type Command interface {
	Execute(handler Handler)
}

type Handler interface {
	Post(cmd Command)
}

type CommandFunc func(h Handler)

func (cmf CommandFunc) Execute(h Handler) {
	cmf(h)
}
