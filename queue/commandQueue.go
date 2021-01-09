package commandQueue

import (
	"sync"

	. "../command"
)

type CommandQueue struct {
	sync.Mutex
	storage []Command
	request chan struct{}
	waiting bool
}

func (cq *CommandQueue) Push(cmd Command) {
	cq.Lock()
	defer cq.Unlock()
	cq.storage = append(cq.storage, cmd)
	if cq.waiting {
		cq.request <- struct{}{}
		cq.waiting = false
	}
	return
}

func (cq *CommandQueue) Pull() Command {
	cq.Lock()
	defer cq.Unlock()
	if cq.Empty() {
		cq.waiting = true
		cq.Unlock()
		<-cq.request
		cq.Lock()
	}
	cmd := cq.storage[0]
	cq.storage[0] = nil
	cq.storage = cq.storage[1:]
	return cmd
}

func (cq *CommandQueue) Empty() bool {
	return len(cq.storage) == 0
}

func InitCommandQueue() CommandQueue {
	return CommandQueue{
		storage: []Command{},
		request: make(chan struct{}),
		waiting: false,
	}
}
