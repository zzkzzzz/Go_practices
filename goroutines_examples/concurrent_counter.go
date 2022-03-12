package main

import (
	"fmt"
	"sync"
)

func main() {

	goroutines_without_sync()
	// goroutines_with_sync()
	// goroutines_with_sync2()
	// goroutines_with_channel_but_deadlock()
	goroutines_with_channel_no_deadlock()
	goroutines_with_channel_with_buffer()
}

func goroutines_without_sync() {
	count := 0
	go func() {
		count++
	}()
	go func() {
		count++
	}()
	// “(the) exit of a goroutine is not guaranteed to happen before any event in the program.”
	// count is very likely going to be 0
	fmt.Println("Count1: ", count)
}

func goroutines_with_sync() {
	// data race may happen
	count := 0

	var wg sync.WaitGroup
	wg.Add(2) // 2 goroutines to wait
	go func() {
		count++
		wg.Done()
	}()
	go func() {
		count++
		wg.Done()
	}()
	wg.Wait() // wait until the 2 goroutines have done

	fmt.Println("Count2: ", count)
}

func goroutines_with_sync2() {
	// data race may happen
	count := 0

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1) // goroutine to wait
		go func() {
			count++
			wg.Done()
		}()
	}
	wg.Wait() // wait until the 1000 goroutines are done

	fmt.Println("Count3: ", count)
}

func goroutines_with_channel_but_deadlock() {
	ch := make(chan int) // make unbuffered channel
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1) // make main wait for 1 more goroutine
		go func() {
			defer wg.Done() // defer to release a wait
			// read count, increment, send away
			count := <-ch // blocked until the previous goroutine sent
			count++
			ch <- count // deadlock will happen here
			// the last goroutine is waiting for a receiver before it signals the waitgroup
			// but the main also wait for the waitGroup signal before reading from channel
		}()
	}
	ch <- 0   // main sends initial value; block until received
	wg.Wait() // wait for all goroutines
	fmt.Println("Count: ", <-ch)
}

func goroutines_with_channel_no_deadlock() {
	ch := make(chan int) // make unbuffered channel
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1) // make main wait for 1 more goroutine
		go func() {
			// read count, increment, send away
			count := <-ch // blocked until sent
			count++
			// Resolve deadlock by reversing the order
			wg.Done()   // B
			ch <- count // A
		}()
	}
	ch <- 0   // main sends initial value; block until received
	wg.Wait() // wait for all goroutines
	// main will be released from wg blockage, and when main tries to read the count off ch,
	// it can only read from the last goroutine that unblocks wg.Wait() and will therefore get the final count.
	fmt.Println("Count: ", <-ch)
}

func goroutines_with_channel_with_buffer() {
	ch := make(chan int, 1) // make buffered channel of size 1
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1) // make main wait for 1 more goroutine
		go func() {
			// any sender will be able to send to the channel in a non-blocking manner
			// as long as the channel is not full.
			ch <- 1 + <-ch
			wg.Done()
		}()
	}
	ch <- 0   // main sends initial value; block until received
	wg.Wait() // wait for all goroutines
	fmt.Println("Count: ", <-ch)
}
