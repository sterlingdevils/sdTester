package main_test

import (
	"log"

	"github.com/sterlingdevils/pipelines/pkg/bufferpipe"
	"github.com/sterlingdevils/pipelines/pkg/containerpipe"
	"github.com/sterlingdevils/pipelines/pkg/udppipe"
)

type Node struct {
	id int
}

func (n Node) Key() int {
	return n.id
}

// bufferpipe
func Example_checkinchan_buffer() {
	b1, err := bufferpipe.New[Node](1)
	if err != nil {
		log.Fatalln(err)
	}
	b2, err := bufferpipe.NewWithPipeline[Node](1, b1)
	if err != nil {
		log.Fatalln(err)
	}

	if b1.InChan() != b2.InChan() {
		log.Fatalln("inchan is not the same channel")
	}

	// Output:
}

// containerpipe
func Example_checkinchan_container() {
	b1, err := containerpipe.New[int, Node]()
	if err != nil {
		log.Fatalln(err)
	}
	b2, err := containerpipe.NewWithPipeline[int, Node](b1)
	if err != nil {
		log.Fatalln(err)
	}

	if b1.InChan() != b2.InChan() {
		log.Fatalln("inchan is not the same channel")
	}

	// Output:
}

// dirscanpipe

// retrypipe
// func Example_checkinchan_retry() {
// 	b1, err := retrypipe.New()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	b2, err := retrypipe.NewWithPipeline(b1)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	if b1.InChan() != b2.InChan() {
// 		log.Fatalln("inchan is not the same channel")
// 	}

// 	// Output:
// }

// udppipe
func Example_checkinchan_udp() {
	b1, err := udppipe.New(9876)
	if err != nil {
		log.Fatalln(err)
	}
	b2, err := udppipe.NewWithPipeline(9999, b1)
	if err != nil {
		log.Fatalln(err)
	}

	if b1.InChan() != b2.InChan() {
		log.Fatalln("inchan is not the same channel")
	}

	// Output:
}
