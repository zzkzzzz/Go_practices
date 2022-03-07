package main

import "fmt"

func intSeq() func() int {
    i := 0
    return func() int {
        i++
        return i
    }
}

func main() {

    nextInt := intSeq()

	// This function value captures its own i value, which will be updated each time we call nextInt.
    fmt.Println(nextInt())
    fmt.Println(nextInt())
    fmt.Println(nextInt())

	// To confirm that the state is unique to that particular function, create and test a new one.
    newInts := intSeq()
    fmt.Println(newInts())
}