package main

import "fmt"

func main() {

	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)

	// the iteration terminates after receiving the 2 elements.
	// This example also showed that it’s possible to close a non-empty channel but still have the remaining values be received.
	for elem := range queue {
		fmt.Println(elem)
	}
}
