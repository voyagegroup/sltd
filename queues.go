package main

import ()

type queues struct {
	parser   *chan message
	transfer *chan message
	remover  *chan message
}

type message struct {
	filepath string
	content  string
}

func NewQueues() *queues {
	q := new(queues)

	qp := make(chan message)
	q.parser = &qp

	qt := make(chan message)
	q.transfer = &qt

	qr := make(chan message)
	q.remover = &qr

	return q
}
