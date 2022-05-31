package main_test

import (
	"fmt"
	"log"
	"time"

	"github.com/sterlingdevils/pipelines/bufferpipe"
	"github.com/sterlingdevils/pipelines/containerpipe"
	"github.com/sterlingdevils/pipelines/logpipe"
	"github.com/sterlingdevils/pipelines/udppipe"
)

type Node struct {
	id int
}

func (n Node) Key() int {
	return n.id
}

// bufferpipe
func Example_buffer() {
	b1, err := bufferpipe.New[int](1)
	if err != nil {
		log.Fatalln(err)
	}
	b2, err := bufferpipe.NewWithPipeline[int](1, b1)
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
		logpipe.NewWithPipeline[Node]("log1", b1),
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

func Example_packet() {
	p := udppipe.Packet{DataSlice: []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x31, 0x32, 0x33,
	}}

	k := udppipe.KeyablePacket(p)
	fmt.Println(k.Key())

	// Output:
	// 578437695752307201
}
