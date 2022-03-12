package main

import (
	"context"
	"fmt"
	"time"
)

// Due to the ubiquity of the use of done channels, Go 1.7 introduces the context package that does the same thing and more.
// When we write a goroutine that spawns a number of goroutines that might each acquire some resources
//  (e.g. memory, file descriptors, database connection) and will exit during the program lifetime,
// we want to release the resources held as soon as the former exits.
type queue chan int

func (q queue) try_enqueue(num int) bool {
	select {
	case q <- num:
		return true
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

func producer(ctx context.Context, q queue) {
	for {
		select {
		case q <- 1: // normal enqueue to q with type alias to chan int
		case <-ctx.Done():
			return
		}
	}
}

func consumer(ctx context.Context, q queue, sumCh chan<- int) {
	sum := 0
	for {
		select {
		case <-ctx.Done():
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
	ctx, cancel := context.WithCancel(context.Background())

	start := make(chan struct{})
	q := make(chan int, 10)
	sumChs := make([]chan int, 0, NumConsumer)
	for i := 0; i < NumConsumer; i++ {
		sumCh := make(chan int, 1)
		sumChs = append(sumChs, sumCh)
	}

	for i := 0; i < NumProducer; i++ {
		go func() {
			<-start
			producer(ctx, q)
		}()
	}
	for j := 0; j < NumConsumer; j++ {
		j := j
		go func() {
			<-start
			consumer(ctx, q, sumChs[j])
		}()
	}

	close(start)
	time.Sleep(time.Second)
	cancel() // cancel the context of this run

	sum := 0
	for _, ch := range sumChs {
		sum += <-ch
	}
	fmt.Println("Sum: ", sum)
}
