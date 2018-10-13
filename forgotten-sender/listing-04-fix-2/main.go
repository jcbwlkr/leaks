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
	/*
	   ardan-bkennedy 3 hours ago
	   This timeout should come from the caller. We need to teach that higher levels calls tell lower level calls how long they want to wait.
	*/

	// Make a channel for our goroutine to report its result.
	ch := make(chan string)
	/*
	   ardan-bkennedy 18 hours ago
	   You still need a buffer of 1

	   ardan-bkennedy 18 hours ago
	   Don't go backwards on code we jsut taught

	   jcbwlkr 11 hours ago  Owner
	   I disagree. This part of the code is presented as an alternative solution to using the buffered channel, not a further improvement. An unbuffered channel is a better choice here in terms of both resource usage and correctness.

	   From the resource usage stance: If the worker knows it is canceled it doesn't need the buffer because it can abort.

	   From the correctness stance: If there is buffer space the select block may choose to place the value there rather than executing the <-ctx.Done() case. We might interpret this that the worker successfully sent the value but really it didn't. See my updated code where I added fmt calls to print "Worker completed" and "Worker canceled".

	   ardan-bkennedy 3 hours ago
	   You have timing issues and therefore a race. if it is "correct" to do something, then don't teach against it for any reason.

	   ardan-bkennedy 3 hours ago
	   I will state again, you just taught, as we do in the classroom, if you dont see a buffer of 1 in a cancellation pattern, you will have problems. Then this solution states it's ok if using context, which is not true.
	*/

	// Start a worker to do some work then either send on the channel or abort if
	// the context was canceled.
	go func() {

		// Get the result of the work. Ideally ctx could be passed down into
		// doSomeWork which would know how to cancel early. Realistically, in many
		// cases it is not cancellable and the worker must wait for it to finish.
		result := doSomeWork()
		/*
		   ardan-bkennedy 3 hours ago
		   This goroutine should always look the same.
		*/

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
