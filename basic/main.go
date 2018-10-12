// basic is a program used to define a goroutine leak.
package main

import (
	"fmt"
	"runtime"
)

func main() {

	// Report number of goroutines. Should be 1.
	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())

	leak()

	// Report new number of goroutines. Will be 2.
	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
}

// leak is a buggy function. It launches a goroutine that blocks reading from a
// channel. Nothing will ever be sent on that channel and the channel is never
// closed so that goroutine will be blocked forever.
func leak() {
	ch := make(chan int)

	go func() {
		val := <-ch
		fmt.Println("We received a value:", val)
	}()
}
