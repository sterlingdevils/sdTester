package main_test

import (
	"log"

	"github.com/sterlingdevils/pipelines/pkg/bufferpipe"
)

type Node struct {
}

func Example_checkinchan_buffer() {
	b1, err := bufferpipe.New[Node](1)
	if err != nil {
		log.Fatalln(err)
	}
	b2, err := bufferpipe.NewWithPipeline[Node](1, b1)
	if err != nil {
		log.Fatalln(err)
	}

	c1 := b1.InChan()
	c2 := b2.InChan()

	if c1 != c2 {
		log.Fatalln("inchan is not the same channel")
	}

	// Output:
}
