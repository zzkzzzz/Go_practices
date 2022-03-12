package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Event struct {
	id       int64
	procTime time.Duration
}

type EventFunc func(Event) Event

type worker struct {
	inputCh  <-chan Event
	outputCh chan<- Event
}

func newWorker(inputCh, outputCh chan Event) *worker {
	return &worker{
		inputCh:  inputCh,
		outputCh: outputCh,
	}
}

func (w *worker) start(
	done <-chan struct{},
	fn EventFunc, wg *sync.WaitGroup,
) {
	go func() {
		defer wg.Done()
		for {
			select {
			case e, more := <-w.inputCh:
				if !more {
					return
				}
				select {
				case w.outputCh <- fn(e):
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()
}

func genEventsCh() chan Event {
	outputCh := make(chan Event)
	go func() {
		counter := int64(1)
		rand.Seed(time.Now().Unix())
		for i := 0; i < 30; i++ {
			outputCh <- Event{
				id:       counter,
				procTime: time.Duration(rand.Intn(100)) * time.Millisecond,
			}
			counter++
		}
		close(outputCh)
	}()
	return outputCh
}

func main() {
	done := make(chan struct{})
	outputCh := make(chan Event, 1)

	inputCh := genEventsCh()

	// Fan-out the stream of input to multiple workers
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		newWorker(inputCh, outputCh).
			start(done, func(e Event) Event {
				time.Sleep(e.procTime)
				return e
			}, &wg)
	}

	readerDone := make(chan struct{})
	go func() {
		for output := range outputCh {
			fmt.Printf("Event id: %d\n", output.id)
		}
		close(readerDone)
	}()

	wg.Wait()

	// Close outputCh and wait for reader to finish reading
	close(outputCh)
	<-readerDone
}
