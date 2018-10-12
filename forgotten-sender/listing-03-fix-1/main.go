// Example program with the goroutine leak fixed. We create capacity in our
// channel so the goroutine can place its value somewhere and die off.
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func main() {

	// Seed the random number generator so we get different results.
	// The seed 42 always results in a timeout.
	// The seed 99 always results in success.
	// In production we may seed using the current time like
	// rand.Seed(time.Now().UnixNano())
	rand.Seed(42)

	fmt.Printf("Number of goroutines: %d\n\n", runtime.NumGoroutine())

	process()

	// Sleep long enough to ensure the goroutine has finished.
	time.Sleep(200 * time.Millisecond)

	fmt.Printf("\nNumber of goroutines: %d\n", runtime.NumGoroutine())
}

func process() {

	// Make a channel for our goroutine to report its result. It has a capacity
	// of 1 so the goroutine will not be blocked sending.
	ch := make(chan string, 1)

	// Start a worker to do some work then send on the channel.
	go func() {
		ch <- doSomeWork()
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

// doSomeWork simulates a function that may take up to 200ms to do something.
func doSomeWork() string {
	delay := time.Duration(rand.Intn(200)) * time.Millisecond
	time.Sleep(delay)
	return "some value"
}
