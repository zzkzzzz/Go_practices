package main

import (
	"fmt"
	"time"
)

// In `blocking_queue.go` and `non_blocking_queue.go`
// From the way consumer works, we can see that the consumers’ use of sumCh is sequential.
// This means that adding more consumers will not speed up the process.

func producer(done chan struct{}, q chan int) {
	for {
		select {
		case q <- 1: // normal enqueue to q with type alias to chan int
		case <-done:
			return
		}
	}
}

func consumer(done chan struct{}, q chan int, sumCh chan<- int) {
	// Instead of storing one sum used by consumers sequentially, we can store many instances of sum on all consumers.
	// local sum
	sum := 0
	for {
		select {
		case <-done:
			// At the end of the process, we can send the sum back to main for consumption
			sumCh <- sum
			close(sumCh)
			return
		case num := <-q:
			sum += num
		}
	}
}

var (
	NumProducer = 5
	NumConsumer = 5
)

func main() {
	start, done := make(chan struct{}), make(chan struct{})
	q := make(chan int)
	sumChs := make([]chan int, 0, NumConsumer)
	for i := 0; i < NumConsumer; i++ {
		sumCh := make(chan int, 1)
		sumChs = append(sumChs, sumCh)
	}

	for i := 0; i < NumProducer; i++ {
		go func() {
			<-start
			producer(done, q)
		}()
	}
	for j := 0; j < NumConsumer; j++ {
		j := j // capture j in the scope
		go func() {
			<-start
			consumer(done, q, sumChs[j])
		}()
	}

	close(start)            // signal to all goroutines to start
	time.Sleep(time.Second) // run for 2 seconds
	close(done)             // signal to all goroutines they should exit

	// collect all sums
	sum := 0
	for _, ch := range sumChs {
		sum += <-ch
	}
	fmt.Println("Sum: ", sum)
}
