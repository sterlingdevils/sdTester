package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sterlingdevils/gobase"
	"github.com/sterlingdevils/pipelines"
)

// Creates a pipeline of UDP pipes connected by a rate limiter pipe
func createpipeline() pipelines.Pipeline[pipelines.Packetable] {
	mysn := (&gobase.SerialNum{}).New()

	// Creates new UDP pipe on port 9876
	inpipe, err := pipelines.UDPPipe{}.New(9876)
	if err != nil {
		log.Fatalln(err)
	}

	logpipetype := pipelines.LogPipe[pipelines.Packetable]{}
	// Create Logger
	l1 := logpipetype.NewWithPipeline("log1", inpipe)

	ptp := pipelines.ConverterPipe[pipelines.Packetable, pipelines.Packetable]{}
	// Create Convert pipe to go from Packet to KeyablePacket by inserting a serial number
	ptokp := ptp.NewWithPipeline(l1,
		func(i pipelines.Packetable) (pipelines.Packetable, error) {
			kp := &pipelines.KeyablePacket{
				Addr:      i.Address(),
				DataSlice: mysn.AddInc(i.Data()),
			}
			return kp, nil
		})

	// Create Logger
	l2 := logpipetype.NewWithPipeline("kp", ptokp)

	// Create Convert pipe to go from Keyable Packet to Packet
	kptop := ptp.NewWithPipeline(l2,
		func(i pipelines.Packetable) (pipelines.Packetable, error) {
			d, _, _ := (&gobase.SerialNum{}).Remove(i.Data())
			p := &pipelines.Packet{
				Addr:      i.Address(),
				DataSlice: d,
			}
			return p, nil
		})

	// Create Logger
	l3 := logpipetype.NewWithPipeline("log3", kptop)

	// Creates a new outpipe using the ratelimiter pipe on port 9999
	outpipe, err := pipelines.UDPPipe{}.NewWithPipeline(9999, l3)
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
