package main

import (
	"math/rand"
	"sync"
)

func Generator[T, K any](done <-chan T, fn func() K) <-chan K {
	stream := make(chan K)
	go func() {
		defer func() {
			close(stream)
		}()
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()

	return stream
}

func Creator() int {
	return rand.Intn(500000000)
}

func primeFinder(rantInts <-chan int) <-chan int {
	res := make(chan int)

	isPrime := func(prime int) bool {
		if prime == 2 {
			return true
		}
		for i := 2; i <= prime/2; i++ {
			if prime%i == 0 {
				return false
			}
		}
		return true
	}

	go func() {
		defer close(res)
		for num := range rantInts {
			go func() {
				if isPrime(num) {
					res <- num
				}
			}()

		}
	}()
	return res
}

func take[T any](stream <-chan T, n int) <-chan T {
	take := make(chan T)
	go func() {
		defer close(take)
		for i := 0; i < n; i++ {
			select {
			case take <- <-stream:
			}
		}
	}()
	return take
}

//fan in

func fanIn[T any](channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	res := make(chan T)

	transfer := func(c <-chan T) {
		defer wg.Done()
		for v := range c {
			res <- v
		}
	}

	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}
