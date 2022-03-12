package main

import (
	"fmt"
	"time"
)

func producer(done chan struct{}, q chan<- int) {
	for {
		select {
		case q <- 1: // keeps sending 1 to q
		case <-done: // should exit
			return
		}
	}
}

func consumer(done chan struct{}, q <-chan int, sumCh chan int) {
	for {
		select {
		case num := <-q: // keeps receiving data from q
			select {
			case sum := <-sumCh:
				sumCh <- num + sum // sequentially increments sum
			case <-done: // should exit
				return
			}
		case <-done: // should exit
			return
		}
	}
}

func main() {
	done := make(chan struct{})
	q, sumCh := make(chan int), make(chan int)

	for i := 0; i < 10; i++ {
		go producer(done, q)
	}
	for j := 0; j < 5; j++ {
		go consumer(done, q, sumCh)
	}

	sumCh <- 0                  // sends initial sum to unblock all consumers and producers
	time.Sleep(2 * time.Second) // runs for 2 seconds
	close(done)                 // signals	 to all goroutines they should exit

	// Once the program is done with producing and consuming the data, main will close(done).
	// close of a channel will free up any goroutines that are blocked at reading from or writing to this channel.
	// Any reader can repeatedly read the zero value of the channel’s data type from a closed channel.

	// Usually, case <-done: is blocked and will not be chosen;
	// however, when main closes done, any goroutines blocked at both cases, most,
	// if not all, of the producers and consumers, will choose case <-done: return.

	// The reason why it is “most” but not all is that: when a select has more than one case that are non-blocking,
	// the select will randomly (with equal weightage) choose one of the cases to proceed.
	// This mechanism is introduced to prevent starvation for some of the cases, and ensure fairness to any blocked cases.

	// That said, more accurately, all goroutines will exit from case <-done.

	fmt.Println("Sum: ", <-sumCh)
}
