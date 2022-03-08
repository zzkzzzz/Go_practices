// To wait for multiple goroutines to finish, we can use a wait group.
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int) {
	fmt.Printf("Worker %d starting\n", id)

	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {

	// This WaitGroup is used to wait for all the goroutines launched here to finish.
	// Note: if a WaitGroup is explicitly passed into functions, it should be done by pointer.
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		// why create a new variable?
		// What happens with closures running as goroutines?
		// https://go.dev/doc/faq#closures_and_goroutines
		i := i

		go func() {
			defer wg.Done()
			worker(i)
		}()
	}

	wg.Wait()

}
