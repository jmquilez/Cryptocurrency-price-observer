// Credits: Juan José Serrano, José Miguel Quilez.
// Main file.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"p1/Observer"
	"p1/Subject"
	"strings"
	"sync"
	"syscall"
)

var (
	// Global flag to indicate shutdown is in progress
	isShuttingDown bool
	// Mutex to protect the flag
	shutdownMutex sync.Mutex
	// WaitGroup to track active file operations
	activeFileOps sync.WaitGroup
)

// parseObserver transforms a string into an Observer.Observer.
// The string must have the format: "observerID,CURRENCY1,CURRENCY2,..."
// Example: "observer1,BTC,ETH,ADA"
func parseObserver(input string) Observer.Observer {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}
	tokens := strings.Split(input, ",")
	if len(tokens) < 1 {
		return nil
	}
	id := strings.TrimSpace(tokens[0])
	btcOk, ethOk, adaOk := false, false, false
	for _, token := range tokens[1:] {
		token = strings.TrimSpace(strings.ToUpper(token))
		if token == "BTC" {
			btcOk = true
		} else if token == "ETH" {
			ethOk = true
		} else if token == "ADA" {
			adaOk = true
		}
	}
	return Observer.NewConcreteObserver(id, btcOk, ethOk, adaOk, &activeFileOps, &shutdownMutex, &isShuttingDown)
}

// observersFromInput processes the input and returns a slice of Observer.Observer.
// First, it separates the input by ";" (each observer) and then processes each one.
func observersFromInput(input string) []Observer.Observer {
	parts := strings.Split(input, ";")
	observers := []Observer.Observer{}
	for _, part := range parts {
		if obs := parseObserver(part); obs != nil {
			observers = append(observers, obs)
		}
	}
	return observers
}

func main() {
	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle shutdown signals
	go func() {
		<-sigChan
		fmt.Println("\nShutdown signal received, cleaning up...")

		// Set shutdown flag
		shutdownMutex.Lock()
		isShuttingDown = true
		shutdownMutex.Unlock()

		// Wait for all file operations to complete
		fmt.Println("Waiting for file operations to complete...")
		activeFileOps.Wait()

		fmt.Println("Clean shutdown complete")
		os.Exit(0)
	}()

	// Create a reader for input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please specify the observers and their preferences in the format as done in this example:")
	fmt.Println("observer1,BTC,ETH,ADA; observer2,BTC,ETH; observer3,ETH,ADA;")
	input, _ := reader.ReadString('\n')

	// Process the input to create the observers
	observers := observersFromInput(input)

	// Initialize the subject; NewConcreteSubject does not receive arguments.
	subject := Subject.NewConcreteSubject()

	// Attach the observers to the subject to receive updates
	for _, observer := range observers {
		subject.Attach(observer)
	}

	// Start listening to the websockets in a goroutine
	go subject.StartListening()

	// Keep the program running until Enter is pressed
	fmt.Println("Press Enter to exit...")
	reader.ReadString('\n')
}
