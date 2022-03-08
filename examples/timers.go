package main

import (
	"fmt"
	"time"
)

//  execute Go code at some point in the future, or repeatedly at some interval.
//  Go’s built-in timer and ticker features make both of these tasks easy.
func main() {

	timer1 := time.NewTimer(2 * time.Second)

	// The <-timer1.C blocks on the timer’s channel C until it sends a value indicating that the timer fired.
	<-timer1.C
	fmt.Println("Timer 1 fired")

	// If you just wanted to wait, you could have used time.Sleep.
	// time.Sleep(time.Second)
	// One reason a timer may be useful is that you can cancel the timer before it fires. Here’s an example of that.
	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 fired")
	}()
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}

	time.Sleep(2 * time.Second)
}
