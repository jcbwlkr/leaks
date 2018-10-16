// Example program showing a goroutine leak. It launches a
// goroutine that sends on a channel but sometimes there is
// no other goroutine available to receive.
package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"
)

func main() {

	// Report number of goroutines. Will be 1.
	fmt.Println("Number of goroutines:", runtime.NumGoroutine())

	process("gophers")

	// Hold the program from terminating for 1 second to see
	// if any goroutines created by process will terminate.
	time.Sleep(time.Second)

	// Report number of goroutines. Will be 2.
	fmt.Println("Number of goroutines:", runtime.NumGoroutine())
}

// process is the work for the program. It finds a record
// then prints it. It fails if it takes more than 100ms.
func process(term string) {

	// Create a context that will be canceled in 100ms.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Make a channel for the goroutine to report its result.
	ch := make(chan string)

	// Launch a goroutine to find the record. Send the return
	// value on the channel.
	go func() {
		ch <- search(term)
	}()

	// Block waiting to receive from the goroutine's channel
	// or for the context to be canceled.
	select {
	case result := <-ch:
		fmt.Println("Received:", result)
	case <-ctx.Done():
		log.Println("search canceled")
	}
}

// search simulates a function that finds a document based
// on a search term. It takes 200ms to perform this work.
func search(term string) string {
	delay := time.Duration(200 * time.Millisecond)
	time.Sleep(delay)
	return "some value"
}
