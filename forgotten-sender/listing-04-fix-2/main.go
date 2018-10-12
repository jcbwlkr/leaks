// Example program with the goroutine leak fixed. The code cancels the worker
// on its way out of the process function.
package main

import (
	"context"
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

	process(context.Background())

	// Sleep long enough to ensure the goroutine has finished.
	time.Sleep(200 * time.Millisecond)

	fmt.Printf("\nNumber of goroutines: %d\n", runtime.NumGoroutine())
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
		case <-ctx.Done():
		}
	}()

	// Wait for a result from the goroutine or the context to be canceled.
	select {
	case result := <-ch:
		fmt.Println("Received:", result)
	case <-ctx.Done():
		fmt.Println("Canceled")
	}
}

// doSomeWork simulates a function that may take up to 200ms to perform some
// processing but can be canceled early.
func doSomeWork(ctx context.Context) string {
	delay := time.Duration(rand.Intn(200)) * time.Millisecond

	select {
	case <-time.After(delay):
		return "some value"
	case <-ctx.Done():
		return ""
	}
}
