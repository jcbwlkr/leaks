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

	// Start a worker to do some work then either send on the channel or abort if
	// the context was canceled.
	go func() {

		// Get the result of the work. Ideally ctx could be passed down into
		// doSomeWork which would know how to cancel early. Realistically, in many
		// cases it is not cancellable and the worker must wait for it to finish.
		result := doSomeWork()

		// Send on the channel if possible or abort if the context is canceled.
		select {
		case ch <- result:
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

// doSomeWork simulates a function that takes up to 200ms to finish some work.
func doSomeWork() string {
	delay := time.Duration(200 * time.Millisecond)
	time.Sleep(delay)
	return "some value"
}
