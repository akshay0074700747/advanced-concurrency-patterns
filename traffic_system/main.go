package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	North = iota
	West
	East
	South
)

type Vehicle struct {
	from int
	to   int
	data int
}

type Signal struct {
	Discover map[int]*Road
}

func newVehicle(frm, to, dta int) Vehicle {
	return Vehicle{
		from: frm,
		to:   to,
		data: dta,
	}
}

func newRoad(dir int) *Road {
	return &Road{
		fromChan:  make(chan Vehicle, 10),
		toChan:    make(chan Vehicle, 10),
		direction: dir,
	}
}

type Road struct {
	fromChan  chan Vehicle
	toChan    chan Vehicle
	direction int
}

func main() {

	var roads []*Road

	var (
		done = make(chan int)
	)

	fmt.Println("starting....")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r, "panic occured ... Gracefully shutting down")
		}
		fmt.Println("stopping...")
	}()

	for i := 0; i < 4; i++ {
		road := newRoad(i)
		go road.startTraffic()
		roads = append(roads, road)
	}

	signal := newSignal(roads...)

	go signal.signaling(done)

	<-done

}

func (rd *Road) startTraffic() {

	defer func() {
		close(rd.fromChan)
		close(rd.toChan)
	}()
	randdta := func() int {
		return rand.Intn(1000)
	}
	roadSelector := func() int {
		return rand.Intn(4)
	}
	for {
		newVehicle := newVehicle(rd.direction, roadSelector(), randdta())
		select {
		case v := <-rd.toChan:
			fmt.Println("vehicle reached road", rd.direction, v)
		case rd.fromChan <- newVehicle:
			fmt.Println("vehicle sent from road", rd.direction)
			time.Sleep(time.Second * 2)
		}
	}
}

func newSignal(roads ...*Road) *Signal {
	signal := &Signal{
		Discover: make(map[int]*Road),
	}
	for _, road := range roads {
		switch road.direction {
		case West:
			signal.Discover[West] = road
		case East:
			signal.Discover[East] = road
		case North:
			signal.Discover[North] = road
		case South:
			signal.Discover[South] = road
		}
	}
	return signal
}

func (s *Signal) signaling(done chan<- int) {

	timee := time.NewTimer(time.Second * 30)
	go func() {
		<-timee.C
		done <- 1
	}()

	for {
		select {
		case v := <-s.Discover[West].fromChan:
			s.Discover[v.to].toChan <- v
		case v := <-s.Discover[East].fromChan:
			s.Discover[v.to].toChan <- v
		case v := <-s.Discover[North].fromChan:
			s.Discover[v.to].toChan <- v
		case v := <-s.Discover[South].fromChan:
			s.Discover[v.to].toChan <- v
		default:

		}
	}
}
