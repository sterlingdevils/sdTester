package main_test

import (
	"log"
	"time"

	"github.com/sterlingdevils/pipelines/bufferpipe"
	"github.com/sterlingdevils/pipelines/containerpipe"
	"github.com/sterlingdevils/pipelines/logpipe"
)

type Node struct {
	id int
}

func (n Node) Key() int {
	return n.id
}

// bufferpipe
func Example_buffer() {
	b1, err := bufferpipe.New[Node](1)
	if err != nil {
		log.Fatalln(err)
	}
	b2, err := bufferpipe.NewWithPipeline[Node](1, b1)
	if err != nil {
		log.Fatalln(err)
	}

	b2.Close()

	// Output:
}

// containerpipe
//    CP -> LogPipe -> CP
func Example_container() {
	b1 := containerpipe.New[int, Node]()
	b2 := containerpipe.NewWithPipeline[int, Node](
		logpipe.NewWithPipeline[Node](b1),
	)

	b1.InChan() <- Node{id: 1}

	time.Sleep(100 * time.Millisecond)
	<-b2.OutChan()
	b2.Close()

	// Output:
}

// dirscanpipe

// retrypipe

// udppipe
func Example_udp() {

	// Output:
}
