package util

import "sync/atomic"

type Task struct {
	running atomic.Bool
}

func (t *Task) Start() bool {
	return t.running.CompareAndSwap(false, true)
}

func (t *Task) Finish() {
	t.running.Store(false)
}
