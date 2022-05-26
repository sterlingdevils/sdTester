package main

import (
	"fmt"
	"log"

	"github.com/sterlingdevils/gobase/pkg/serialnum"
	"github.com/sterlingdevils/pipelines"
	"github.com/sterlingdevils/pipelines/converterpipe"
	"github.com/sterlingdevils/pipelines/logpipe"
	"github.com/sterlingdevils/pipelines/retrypipe"
	"github.com/sterlingdevils/pipelines/udppipe"
	"github.com/sterlingdevils/ratelimiterpipe/pkg/ratelimiterpipe"
)

// func p2r() func(*udppipe.Packet) *retrypipe.RetryThing[uint64, *udppipe.Packet] {
// 	mysn := serialnum.New()
// 	return func(i *udppipe.Packet) *retrypipe.RetryThing[uint64, *udppipe.Packet] {
// 		return retrypipe.NewRetryThing(mysn.Next(), i)
// 	}
// }

func r2p(r *retrypipe.RetryThing[uint64, *udppipe.Packet]) *udppipe.Packet {
	return r.Thing
}

// Creates a pipeline of UDP pipes connected by a rate limiter pipe
func createpipeline() pipelines.Pipeliner[*udppipe.Packet] {
	mysn := serialnum.New()

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

	l1 := logpipe.NewWithPipeline[*udppipe.Packet](rlpipe)

	p2 := converterpipe.NewWithPipeline[*udppipe.Packet](
		func(i *udppipe.Packet) *retrypipe.RetryThing[uint64, *udppipe.Packet] {
			return retrypipe.NewRetryThing(mysn.Next(), i)
		},
		l1)

	l2 := logpipe.NewWithPipeline[*retrypipe.RetryThing[uint64, *udppipe.Packet]](p2)
	r2 := converterpipe.NewWithPipeline[*retrypipe.RetryThing[uint64, *udppipe.Packet]](r2p, l2)
	l3 := logpipe.NewWithPipeline[*udppipe.Packet](r2)

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

	//	time.Sleep(2 * time.Second)

	p.Close()
}
