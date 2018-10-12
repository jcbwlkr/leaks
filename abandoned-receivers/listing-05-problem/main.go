// Example program to show a goroutine leak. We start goroutines that range
// over a channel but nothing ever closes the channel.
package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func main() {
	fmt.Printf("Number of goroutines: %d\n\n", runtime.NumGoroutine())

	names := []string{"Anna", "Jacob", "Kell", "Carter", "Rory"}
	processRecords(names)

	// Give goroutines a chance to return before reporting.
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("\nNumber of goroutines: %d\n", runtime.NumGoroutine())
}

// processRecords is given a slice of values such as lines from a file. The
// order of these values is not important so it can start multiple workers to
// perform some processing on each record then feed the results back.
func processRecords(records []string) {
	input := make(chan string)
	output := make(chan string)

	// Start multiple workers to process input and send results to output.
	const workers = 3
	for i := 0; i < workers; i++ {
		go worker(i, input, output)
	}

	// Start a goroutine to feed records to workers.
	go func() {
		for _, record := range records {
			input <- record
		}
	}()

	// Receive from output the expected number of times. If 10 records went in
	// then 10 will come out.
	for i := 0; i < len(records); i++ {
		fmt.Printf("[main    ]: output %s\n", <-output)
	}
}

// worker represents the work that I wish to do in parallel. This is a blog
// post so all the workers do is capitalize a string but you can imagine they
// are doing something more intensive.
//
// I don't know how many records each individual goroutine will need to process
// so they use the range keyword to receive in a loop.
func worker(id int, input <-chan string, output chan<- string) {
	for v := range input {
		fmt.Printf("[worker %d]: input %s\n", id, v)
		output <- strings.ToUpper(v)
	}
}
