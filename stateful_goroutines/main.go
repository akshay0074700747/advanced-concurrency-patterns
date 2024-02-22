package main

import (
	"fmt"
	"time"
)

type readOp struct {
	key  int
	resp chan int
}
type writeOp struct {
	key  int
	val  int
	resp chan bool
}

func main() {
	reaChan := make(chan readOp)
	writeChan := make(chan writeOp)
	for i := 0; i < 5; i++ {
		go func() {
			for {
				res := writeOp{
					val:  1,
					resp: make(chan bool),
				}
				rea := readOp{
					resp: make(chan int),
				}
				select {
				case writeChan <- res:
					<-res.resp
					fmt.Println("write successful")
				case reaChan <- rea:
					fmt.Println("value read ", <-rea.resp)
				}
			}
		}()
	}
	go func() {
		var res int
		for {
			select {
			case read := <-reaChan:
				read.resp <- res
			case write := <-writeChan:
				res += write.val
			}
		}
	}()
	time.Sleep(time.Second * 15)
}
