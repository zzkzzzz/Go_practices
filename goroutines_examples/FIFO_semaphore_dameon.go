package main

import (
	"container/list"
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

// Helper struct that extends list.List with a helper method
type chanQueue struct{ list.List }

func NewChanQueue() *chanQueue {
	q := new(chanQueue)
	q.Init()
	return q
}

func (q *chanQueue) Pop() chan struct{} {
	ele := q.Front()
	q.Remove(ele)
	return ele.Value.(chan struct{})
}

///////////////////////////////////////////////////////////////

// Centralised semaphore design using a daemon goroutine
// Each Acquire and Release communicates with the daemon,
// which keeps waiters blocked if necessary,
// and chooses which waiters to unblock based on a FIFO queue
type Semaphore2 struct {
	acquireCh chan chan struct{}
	releaseCh chan struct{}
}

func NewSemaphore2(initial_count int) *Semaphore2 {
	sem := new(Semaphore2)
	sem.acquireCh = make(chan chan struct{}, 100)
	sem.releaseCh = make(chan struct{}, 100)

	go func() {
		count := initial_count
		// The FIFO queue that stores the channels used to unblock waiters
		waiters := NewChanQueue()

		for {
			select {
			case <-sem.releaseCh: // Increment or unblock a waiter
				if waiters.Len() > 0 {
					ch := waiters.Pop()
					ch <- struct{}{} // Unblocks the oldest waiter
				} else {
					count++
				}

			case ch := <-sem.acquireCh: // Decrement or add a waiter
				if count > 0 {
					count--
					ch <- struct{}{} // Don't keep waiter blocked
				} else {
					waiters.PushBack(ch) // Add waiter to back of queue
				}
			}
		}
	}()

	return sem
}

// Technically, it is possible that an acquire is blocked on the first send to s.acquireCh,
// even before it’s able to send its channel to the daemon. If we assume that channels do not unblock in FIFO order,
// it’s possible that it remains blocked on this first send forever while other goroutines are constantly sending new acquire requests to the daemon.

// So the answer is no, it’s not actually FIFO in the sense that a goroutine A that calls Acquire
// can be blocked before another goroutine B, and yet goroutine B unblocks before goroutine A.

// However, the ordering is enforced from the moment that the first send actually succeeds.
// Since the daemon is capable of emptying its request queues relatively quickly, and the request
// queues can be buffered to a length where it does not block in practice, it is possible to make an
// argument that this semaphore is FIFO under certain conditions.
func (s *Semaphore2) Acquire() {
	ch := make(chan struct{})
	// Send daemon a channel that can be used to unblock us
	s.acquireCh <- ch
	// Block until daemon decides to unblock us
	<-ch
}

func (s *Semaphore2) Release() {
	s.releaseCh <- struct{}{}
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
		ops = stresstest(NewSemaphore2(0), num_releasers, num_goroutines)
		fmt.Fprint(os.Stderr, "Semaphore2")
		for _, i := range ops {
			fmt.Fprintf(os.Stderr, "\t%d", i)
		}
		fmt.Fprintln(os.Stderr)
	case 2:
		fifotest(NewSemaphore2(0), num_goroutines)
	}
}
