package main

import (
	"fmt"
	"time"
)

func main() {
	var (
		global = make(chan int)
		even   = make(chan int)
		odd    = make(chan int)
		done   = make(chan int)
		output = make(chan int)
	)
	defer func() {
		close(even)
		close(odd)
		close(done)
		close(output)
	}()
	go evenProcessor(even, output)
	go oddProcessor(odd, output)
	go distributor(odd, even, global)
	go generator(global, done)
	go func() {
		time.Sleep(time.Second * 2)
		done <- 1
	}()
	for out := range output {
		if out == 0 {
			break
		}
		fmt.Println(out)
	}

}

func generator(global chan<- int, done <-chan int) {
	for i := 1; ; i++ {
		select {
		case <-done:
			close(global)
		default:
			global <- i
		}
	}
}

func evenProcessor(even <-chan int, output chan<- int) {
	for {

		eve := <-even
		output <- (eve * 1000)

	}
}

func oddProcessor(odd <-chan int, output chan<- int) {
	for {

		odd := <-odd
		output <- odd

	}
}

func distributor(odd chan<- int, even chan<- int, global <-chan int) {
	for {

		val := <-global
		if val%2 == 0 {
			even <- val
		} else {
			odd <- val
		}

	}
}
