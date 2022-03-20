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

type signal struct {
	isRelease bool
	releaseCh chan signal
}

func getSignal(s signal) (bool, chan signal) {
	return s.isRelease, s.releaseCh
}

type Semaphore3 struct {
	waitQueue   chan signal
	globalRelCh chan signal
}

func NewSemaphore3(capacity int, initial_count int) *Semaphore3 {
	s := Semaphore3{
		waitQueue:   make(chan signal),
		globalRelCh: make(chan signal),
	}

	// fill the waitQueue with initial_count Releases
	go func() {
		// force the first Acquirer to read from globalRelCh
		s.waitQueue <- signal{false, s.globalRelCh}
	}()
	for i := 0; i < initial_count; i++ {
		s.Release()
	}

	return &s
}

// Acquire forms a link in the waitQueue, and try to read from the latest releaseCh.
func (s *Semaphore3) Acquire() {
	// try to see if this is the waitQueue head
	isRelease, relCh := getSignal(<-s.waitQueue)

	// this is a new link in the waitQueue
	// prepare to pass globalRelCh to the next waiter
	nextRelCh := make(chan signal)
	go func() { s.waitQueue <- signal{false, nextRelCh} }()

	// if isRelease is false, releaseCh is not the global release chan
	// if releaseCh is not global release chan, will not read from Release
	for !isRelease {
		isRelease, relCh = getSignal(<-relCh)
	}
	// must have read from a Release, releaseCh must be globalRelCh
	// pass it to the next waiter
	go func() { nextRelCh <- signal{false, s.globalRelCh} }()
}

// Release sends release signal to globalRelCh held by the first waiter.
func (s *Semaphore3) Release() {
	go func() { s.globalRelCh <- signal{true, s.globalRelCh} }()
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
		ops = stresstest(NewSemaphore3(1000000, 0), num_releasers, num_goroutines)
		fmt.Fprint(os.Stderr, "Semaphore3")
		for _, i := range ops {
			fmt.Fprintf(os.Stderr, "\t%d", i)
		}
		fmt.Fprintln(os.Stderr)
	case 2:
		fifotest(NewSemaphore3(1000000, 0), num_goroutines)
	}
}
