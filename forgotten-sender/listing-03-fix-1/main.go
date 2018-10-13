// Example program with the goroutine leak fixed. We create capacity in our
// channel so the goroutine can place its value somewhere and die off.
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {

	// Report number of goroutines. Will be 1.
	fmt.Println("Number of goroutines:", runtime.NumGoroutine())

	process()

	// Hold the program from terminating for 1 second to see if the goroutine
	// created by process will terminate.
	time.Sleep(time.Second)

	// Report number of goroutines. Will be 1.
	fmt.Println("Number of goroutines:", runtime.NumGoroutine())
}

func process() {

	// Make a channel for our goroutine to report its result. It has a capacity
	// of 1 so the goroutine will not be blocked sending.
	ch := make(chan string, 1)

	// Start a worker to do some work then send on the channel.
	go func() {
		ch <- doSomeWork()
		fmt.Println("Worker terminated")
	}()

	// Create a timeout channel. In 100ms a value will be sent on this channel.
	timeout := time.After(100 * time.Millisecond)

	// Wait to receive from the goroutine's channel or the timeout channel,
	// whichever comes first.
	select {
	case result := <-ch:
		fmt.Println("Received:", result)
	case <-timeout:
		fmt.Println("Timed out")
	}
}

// doSomeWork simulates a function that takes up to 200ms to finish some work.
func doSomeWork() string {
	delay := time.Duration(200 * time.Millisecond)
	time.Sleep(delay)
	return "some value"
}
