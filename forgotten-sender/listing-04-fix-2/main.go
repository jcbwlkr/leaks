// Example program with the goroutine leak fixed. The code cancels the worker
// on its way out of the process function.
package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {

	// Report number of goroutines. Will be 1.
	fmt.Println("Number of goroutines:", runtime.NumGoroutine())

	process(context.Background())

	// Hold the program from terminating for 1 second to see if the goroutine
	// created by process will terminate.
	time.Sleep(time.Second)

	// Report number of goroutines. Will be 1.
	fmt.Println("Number of goroutines:", runtime.NumGoroutine())
}

func process(ctx context.Context) {

	// Wrap the parent context in a new context with a 100ms deadline. The
	// context will be canceled if 100ms passes or if the parent function's
	// context is canceled for any reason.
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	// Make a channel for our goroutine to report its result.
	ch := make(chan string)

	// Start the worker. Block waiting to send on the channel or hear that the
	// context is canceled.
	go func() {
		select {
		case ch <- doSomeWork(ctx):
			fmt.Println("Worker completed")
		case <-ctx.Done():
			fmt.Println("Worker canceled")
		}
	}()

	// Wait for a result from the goroutine or the context to be canceled.
	select {
	case result := <-ch:
		fmt.Println("Received:", result)
	case <-ctx.Done():
		fmt.Println("Receiver canceled")
	}
}

// doSomeWork simulates a function that may take up to 200ms to perform some
// processing but can be canceled early.
func doSomeWork(ctx context.Context) string {
	delay := time.Duration(200 * time.Millisecond)

	select {
	case <-time.After(delay):
		return "some value"
	case <-ctx.Done():
		return ""
	}
}
