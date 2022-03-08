package main

import "fmt"

func mayPanic() {
	panic("a problem")
}

func main() {

	// recover must be called within a deferred function.
	// When the enclosing function panics, the defer will activate and a recover call within it will catch the panic.
	defer func() {
		if r := recover(); r != nil {

			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	// execution stop here
	mayPanic()

	// this code will not run due to Panic
	fmt.Println("After mayPanic()")
}
