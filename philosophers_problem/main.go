package main

import (
	"fmt"
	"sync/atomic"
)

type Philosopher struct {
	id       int
	hungry   int
	forks    int
	leftover int
}

type Priority struct {
	queue []Philosopher
}

func NewPriority() *Priority {
	return &Priority{
		queue: []Philosopher{},
	}
}

func (pr *Priority) add(p Philosopher) {
	pr.queue = append(pr.queue, p)
	n := len(pr.queue) - 1
	if n == 1 {
		return
	}
	pr.heapifyup(n)
}

func (pr *Priority) heapifyup(i int) {
	p := getParent(i)
	for pr.queue[i].hungry > pr.queue[p].hungry {
		pr.queue[i], pr.queue[p] = pr.queue[p], pr.queue[i]
		i = p
		p = getParent(i)
	}
}

func (pr *Priority) delete() Philosopher {
	n := pr.queue[0]
	pr.queue[0] = pr.queue[len(pr.queue)-1]
	pr.queue = pr.queue[:len(pr.queue)-1]
	pr.heapifydown(0)
	return n
}

func (pr *Priority) heapifydown(i int) {
	left := getLeft(i)
	right := getRight(i)
	largest := i
	if left < len(pr.queue) && pr.queue[left].hungry > pr.queue[largest].hungry {
		largest = left
	}
	if right < len(pr.queue) && pr.queue[right].hungry > pr.queue[largest].hungry {
		largest = right
	}
	if largest != i {
		pr.queue[largest], pr.queue[i] = pr.queue[i], pr.queue[largest]
		pr.heapifydown(largest)
	}
}

func (pr *Priority) print() {
	for _, v := range pr.queue {
		fmt.Printf("%d => ", v.hungry)
	}
	fmt.Println()
}

func getParent(i int) int {
	if i == 0 {
		return 0
	}
	return ((i - 1) / 2)
}

func getLeft(i int) int {
	return ((i * 2) + 1)
}

func getRight(i int) int {
	return ((i * 2) + 2)
}

func main() {
	done := make(chan Philosopher)
	forkchan := make(chan int)
	leftoverchan := make(chan int)
	needmorechan := make(chan Philosopher)
	var forks atomic.Int32
	forks.Store(5)
	pr := NewPriority()
	hungry := []int{45, 23, 0, 65, 100}
	for i, v := range hungry {
		pr.add(Philosopher{
			id:     i,
			hungry: v,
		})
	}
	pr.print()
	go func() {
		for {
			n := <-forkchan
			forks.Add(int32(n))
		}
	}()
	for _, p := range pr.queue {
		fork := forks.Load()
		if fork == 0 || fork == 1 {
			p := <-done
			if p.leftover != 0 {
				leftoverchan <- p.leftover
			}
			fork = forks.Load()
		}
		go Eat(p, done, forkchan, leftoverchan,needmorechan)
		forks.Store(fork - 2)
	}
}

func Eat(p Philosopher, done chan<- Philosopher, forkchan chan<- int, leftoverchan <-chan int, needmore chan<- Philosopher) {
	select {
	case left := <-leftoverchan:
		if left < p.hungry {
			p.hungry -= left
			needmore <- p
		} else {
			p.hungry = 0
			done <- p
		}
	default:
		p.forks = 2
	}
}
