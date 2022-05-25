package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sterlingdevils/pipelines/pkg/pipeline"
	"github.com/sterlingdevils/pipelines/pkg/udppipe"
	"github.com/sterlingdevils/ratelimiterpipe/pkg/ratelimiterpipe"
)

// Creates a pipeline of UDP pipes connected by a rate limiter pipe
func createpipeline() pipeline.Pipelineable[*udppipe.Packet] {
	// Creates new UDP pipe on port 9876
	inpipe, err := udppipe.New(9876)
	if err != nil {
		log.Fatalln(err)
	}

	// Creates a ratelimiter pipe using the UDP pipe on port 9876 with a rate and burst limit of
	// the maximum value of 16 bit int
	rlpipe, err := ratelimiterpipe.NewWithPipeline[*udppipe.Packet](1, 7, inpipe)
	if err != nil {
		log.Fatalln(err)
	}

	// Creates a new outpipe using the ratelimiter pipe on port 9999
	outpipe, err := udppipe.NewWithPipeline(9999, rlpipe)
	if err != nil {
		log.Fatalln(err)
	}

	return outpipe
}

func main() {
	defer log.Println("leaving main")

	fmt.Printf("Creating Pipeline")
	p := createpipeline()

	time.Sleep(2029 * time.Second)

	p.Close()
}
