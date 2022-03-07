package main

import "fmt"

func plus(a int, b int) int {

    return a + b
}

// you may omit the type name for the like-typed parameters up to the final parameter that declares the type
func plusPlus(a, b, c int) int {
    return a + b + c
}

func main() {

    res := plus(1, 2)
    fmt.Println("1+2 =", res)

    res = plusPlus(1, 2, 3)
    fmt.Println("1+2+3 =", res)
}