package main_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sterlingdevils/pipelines"
)

type Node struct {
	id int
}

func (n Node) Key() int {
	return n.id
}

// bufferpipe
func ExampleBufferPipe() {
	b1, err := pipelines.BufferPipe[int]{}.New(1)
	if err != nil {
		log.Fatalln(err)
	}
	b2, err := pipelines.BufferPipe[int]{}.NewWithPipeline(1, b1)
	if err != nil {
		log.Fatalln(err)
	}

	b2.Close()

	// Output:
}

// containerpipe
//    CP -> LogPipe -> CP
func ExampleContainerPipe() {
	cp := pipelines.ContainerPipe[int, Node]{}
	b1 := cp.New()
	b2 := cp.NewWithPipeline(pipelines.LogPipe[Node]{}.NewWithPipeline("log1", b1))

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
	p := pipelines.Packet{DataSlice: []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x31, 0x32, 0x33,
	}}

	k := pipelines.KeyablePacket(p)
	fmt.Println(k.Key())

	// Output:
	// 578437695752307201
}

type DataHolder struct {
	data []byte
}

func (d DataHolder) Data() []byte {
	return d.data
}

func ExampleConverterPipe() {
	// Create Logger
	logger := pipelines.LogPipe[DataHolder]{}.New("inlog")

	// Can't Use logger directly as a Dataer,  Use a convert pipe to change the channel
	//	conv := converterpipe.NewWithPipeline[DataHolder](logger, func(i DataHolder) (pipelines.Dataer, error) { return i, nil })

	conv2 := pipelines.TypeConverterPipe[DataHolder, pipelines.Dataer]{}.NewWithPipeline(logger)

	os.Chdir("/tmp")
	fd := pipelines.FileDump{}.NewWithPipeline(conv2)

	// Send a Packet
	fd.InChan() <- pipelines.Packet{DataSlice: []byte("Hello, World!")}
	fd.InChan() <- pipelines.Packet{DataSlice: []byte("Gimme Jimmy")}
	fd.InChan() <- pipelines.Packet{DataSlice: []byte("See what happens with special characters\nOn this line")}

	fd.InChan() <- DataHolder{data: []byte("This is another type of input")}

	logger.InChan() <- DataHolder{data: []byte("Input to Logger to make sure it flows")}

	time.Sleep(1 * time.Second)

	fd.Close()

	// Output:
	//
}
