package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	startTime := time.Now()

	done := make(chan int)

	defer func() {
		close(done)
		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)
		fmt.Println("Total runtime:", elapsedTime)
	}()
	stream := Generator(done, Creator)
	cpus := runtime.NumCPU()
	//fanOut
	var primeChans = make([]<-chan int, cpus)
	for i := 0; i < cpus; i++ {
		primeChans[i] = primeFinder(stream)
	}
	fannedInstream := fanIn(primeChans...)
	for val := range take(fannedInstream, 10) {
		fmt.Println(val)
	}
}
