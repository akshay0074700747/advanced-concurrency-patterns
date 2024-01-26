package main

import (
	"log"
	"math/rand"
	"time"
)

type Aucted struct {
	id    int
	price string
}

func main() {

	const (
		aucted_item = "I am worth 500000000000 Billion"
		bidders     = 5
	)

	winnerChan := make(chan Aucted)
	quitChan := make(chan bool)
	resMap := make(map[int]<-chan int)

	bidGenerator := func() int {
		return rand.Intn(50000000)
	}

	for i := 1; i <= bidders; i++ {

		resMap[i] = bid(i, bidGenerator, winnerChan)

	}

	go func() {
		maxBid := -1
		maxBiddedBidder := -1
		auctionEnd := time.NewTimer(time.Second * 20)
		for {
			for id, channel := range resMap {
				select {
				case bid := <-channel:
					log.Printf("%s %d %s %d", "bidder with id", id, "bidded", bid)
					if bid > maxBid {
						maxBid = bid
						maxBiddedBidder = id
					}
				case <-auctionEnd.C:
					if maxBiddedBidder == -1 {
						log.Println("no one bidded... Executing without conducting Auction")
						return
					}
					priceGoesTo := Aucted{
						id:    maxBiddedBidder,
						price: aucted_item,
					}
					log.Printf("%s %d %s %d", "highest bid of", maxBid, "by the bidder with id", maxBiddedBidder)
					winnerChan <- priceGoesTo
					time.Sleep(time.Second * 1)
					quitChan <- true
				default:
					continue
				}
			}
		}
	}()
	<-quitChan
	log.Println("Auction Finished...")
}

func bid(id int, bidGenerator func() int, winnerChan chan Aucted) <-chan int {

	log.Printf("%d %s /n", id, "th bidder came for Auction")
	bidder := make(chan int)

	go func() {
		for {
			bid := bidGenerator()

			select {
			case price := <-winnerChan:
				if price.id == id {
					log.Printf("%s %s %d", price.price, "was won by the bidder with id", id)
					close(bidder)
					return
				} else {
					winnerChan <- price
				}
			case bidder <- bid:
			}
			restTime := rand.Intn(4)
			time.Sleep(time.Second * time.Duration(restTime))
		}
	}()

	return bidder

}
