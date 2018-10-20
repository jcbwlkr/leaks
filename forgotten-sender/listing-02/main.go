// Example program that performs serial work. It does not
// launch any goroutines or suffer from leaks. Yet.
package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	if err := process("gophers"); err != nil {
		log.Print(err)
	}
}

// process is the work for the program. It finds a record
// then prints it.
func process(term string) error {
	record, err := search(term)
	if err != nil {
		return err
	}

	fmt.Println("Received:", record)
	return nil
}

// search simulates a function that finds a record based
// on a search term. It takes 200ms to perform this work.
func search(term string) (string, error) {
	time.Sleep(200 * time.Millisecond)
	return "some value", nil
}
