package main

import (
	"fmt"
	"time"
)

type queue chan int

func (q queue) try_enqueue(num int) bool {
	select {
	case q <- num:
		return true
	// The default case in select provides a nice exit for the goroutine when all the cases are blocked.
	// It makes the channel access asynchronous.
	default:
	}
	return false
}

func (q queue) try_dequeue() (int, bool) {
	select {
	case num := <-q:
		return num, true
	default:
	}
	return 0, false
}

func producer(done chan struct{}, q queue) {
	for {
		select {
		case <-done:
			return
		default:
		}

		if ok := q.try_enqueue(1); !ok {
			// do something
		} else {
			// do something
		}
	}
}

func consumer(done chan struct{}, q queue, sumCh chan int) {
	for {
		select {
		case <-done:
			return
		default:
		}

		num, ok := q.try_dequeue()
		if ok {
			sumCh <- num + <-sumCh
		}
	}
}

var (
	NumProducer = 5
	NumConsumer = 10
)

func main() {
	start, done := make(chan struct{}), make(chan struct{})
	q, sumCh := make(queue, 10), make(chan int, 1)
	sumCh <- 0

	for i := 0; i < NumProducer; i++ {
		go func() {
			<-start
			producer(done, q)
		}()
	}
	for j := 0; j < NumConsumer; j++ {
		go func() {
			<-start
			go consumer(done, q, sumCh)
		}()
	}

	close(start)
	time.Sleep(time.Second)
	close(done)

	fmt.Println("Sum: ", <-sumCh)
	sumCh <- 0 // unblock any consumer
}
