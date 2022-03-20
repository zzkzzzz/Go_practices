package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// Interface for semaphores we will implement
type SemaphoreInterface interface {
	Acquire()
	Release()
}

///////////////////////////////////////////////////////////////

// Simplest semaphore design using a buffered channel
// Initially, the buffered channel is empty, so send will not block
// However, each send fills up the buffered channel until it's full,
// so further sends will block
//
// This is a correct implementation of a semaphore.
// While it is FIFO in practice due to an implementation detail
// in the Go runtime, its FIFO behaviour is not actually guaranteed
// by the Go specification
type Semaphore1 struct {
	sem chan struct{}
}

func NewSemaphore1(capacity int, initial_count int) *Semaphore1 {
	sem := Semaphore1{
		sem: make(chan struct{}, capacity),
	}
	for ; initial_count < capacity; initial_count++ {
		sem.Acquire()
	}
	return &sem
}

func (s *Semaphore1) Acquire() {
	// Send to the channel to decrement the number of empty slots
	// Blocks if there are no slots remaining
	s.sem <- struct{}{}
	// Blocked goroutines will be unblocked in FIFO order as of Go 1.17
}

func (s *Semaphore1) Release() {
	// Receive from the channel to increment the number of empty slots
	<-s.sem
}

///////////////////////////////////////////////////////////////

func stresstest(s SemaphoreInterface, num_releasers int, num_goroutines int) []int {
	releasersCtx, releasersCancel := context.WithCancel(context.Background())
	acquirersCtx, acquirersCancel := context.WithCancel(context.Background())
	var acquirersWg sync.WaitGroup

	opsCh := make(chan int)

	for i := 0; i < num_releasers; i++ {
		release_msg := fmt.Sprintf("T%d: Release\n", i)
		go func() {
			ops := 0
		loop:
			for {
				select {
				case <-releasersCtx.Done():
					break loop
				default:
					fmt.Print(release_msg)
					s.Release()
					ops++
				}
			}
			opsCh <- ops
		}()
	}

	for i := num_releasers; i < num_goroutines; i++ {
		waiting_msg := fmt.Sprintf("T%d: Waiting\n", i)
		unblocked_msg := fmt.Sprintf("T%d: Unblocked\n", i)

		acquirersWg.Add(1)
		go func() {
			ops := 0
		loop:
			for {
				select {
				case <-acquirersCtx.Done():
					break loop
				default:
					fmt.Print(waiting_msg)
					s.Acquire()
					fmt.Print(unblocked_msg)
					ops++
				}
			}
			acquirersWg.Done()
			opsCh <- ops
		}()
	}

	time.Sleep(time.Second)

	acquirersCancel()
	acquirersWg.Wait()
	releasersCancel()

	ops := make([]int, 0, num_goroutines)
	for i := 0; i < num_goroutines; i++ {
		ops = append(ops, <-opsCh)
	}

	return ops
}

func fifotest(sem SemaphoreInterface, num_acquirers int) {
	for i := 0; i < num_acquirers; i++ {
		i := i
		go func() {
			time.Sleep(time.Duration(i) * 50 * time.Millisecond)
			sem.Acquire()
			fmt.Printf("hello from thread %d\n", i)
		}()
	}

	for i := 0; i < num_acquirers; i++ {
		sem.Release()
		time.Sleep(50 * time.Millisecond)
	}
}

func main() {
	if len(os.Args) <= 3 {
		fmt.Fprintln(os.Stderr, "Please specify test ID, num_releasers, and num_goroutines")
		os.Exit(1)
	}

	testId, _ := strconv.Atoi(os.Args[1])
	num_releasers, _ := strconv.Atoi(os.Args[2])
	num_goroutines, _ := strconv.Atoi(os.Args[3])

	var ops []int
	switch testId {
	case 1:
		ops = stresstest(NewSemaphore1(1000000, 0), num_releasers, num_goroutines)
		fmt.Fprint(os.Stderr, "Semaphore1")
		for _, i := range ops {
			fmt.Fprintf(os.Stderr, "\t%d", i)
		}
		fmt.Fprintln(os.Stderr)
	case 2:
		fifotest(NewSemaphore1(1000000, 0), num_goroutines)
	}
}
