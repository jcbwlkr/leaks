// Example program that performs serial work. It does not
// launch any goroutines or suffer from leaks. Yet.
package main

import (
	"fmt"
	"time"
)

func main() {
	process("gophers")
}

// process is the work for the program. It finds a record
// then prints it.
func process(term string) {
	result := search(term)
	fmt.Println("Received:", result)
}

// search simulates a function that finds a document based
// on a search term. It takes 200ms to perform this work.
func search(term string) string {
	delay := time.Duration(200 * time.Millisecond)
	time.Sleep(delay)
	return "some value"
}
