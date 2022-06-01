package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sterlingdevils/gobase/pkg/serialnum"
	"github.com/sterlingdevils/pipelines"
	"github.com/sterlingdevils/pipelines/converterpipe"
	"github.com/sterlingdevils/pipelines/logpipe"
	"github.com/sterlingdevils/pipelines/udppipe"
)

// Creates a pipeline of UDP pipes connected by a rate limiter pipe
func createpipeline() pipelines.Pipeline[udppipe.Packetable] {
	mysn := serialnum.New()

	// Creates new UDP pipe on port 9876
	inpipe, err := udppipe.New(9876)
	if err != nil {
		log.Fatalln(err)
	}

	// Create Logger
	l1 := logpipe.NewWithPipeline[udppipe.Packetable]("log1", inpipe)

	// Create Convert pipe to go from Packet to KeyablePacket by inserting a serial number
	ptokp := converterpipe.NewWithPipeline[udppipe.Packetable](
		func(i udppipe.Packetable) (udppipe.Packetable, error) {
			kp := &udppipe.KeyablePacket{
				Addr:      i.Address(),
				DataSlice: mysn.AddInc(i.Data()),
			}
			return kp, nil
		}, l1)

	// Create Logger
	l2 := logpipe.NewWithPipeline[udppipe.Packetable]("kp", ptokp)

	// Create Convert pipe to go from Keyable Packet to Packet
	kptop := converterpipe.NewWithPipeline[udppipe.Packetable](
		func(i udppipe.Packetable) (udppipe.Packetable, error) {
			d, _, _ := serialnum.Remove(i.Data())
			p := &udppipe.Packet{
				Addr:      i.Address(),
				DataSlice: d,
			}
			return p, nil
		}, l2)

	// Create Logger
	l3 := logpipe.NewWithPipeline[udppipe.Packetable]("log3", kptop)

	// Creates a new outpipe using the ratelimiter pipe on port 9999
	outpipe, err := udppipe.NewWithPipeline(9999, l3)
	if err != nil {
		log.Fatalln(err)
	}

	return outpipe
}

func main() {
	defer log.Println("leaving main")

	fmt.Printf("Creating Pipeline\n")
	p := createpipeline()

	time.Sleep(20 * time.Second)

	p.Close()
}
