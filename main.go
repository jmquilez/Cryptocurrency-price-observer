package main

import (
	Observer "Observer/Observer"
	Subject "Subject/Subject"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func observersFromInput(input string) []Observer.Observer {
	// We parse the input and create the observers
	observers := []Observer.Observer{}

	// We split the input by the semicolon
	observers = strings.Split(input, ";")

	// We split the input by the comma
	observers = strings.Split(observers[0], ",")
	
	// We create the observers
	for _, observer := range observers {
		observers = append(observers, Observer.NewConcreteObserver(observer[0], observer[1], observer[2], observer[3]))
	}
	
	// We return the observers
	return observers
}

// We will now read the arguments from the command line and act accordingly
func main() {
	// Crear un lector de entrada
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please specify the number of observers you want to create and their preferences according to this format: observer1,BTC,ETH,ADA; observer2,BTC,ETH; observer3,BTC,ADA;")
	input, _ := reader.ReadString('\n')

	// We parse the input and create the observers
	observers := observersFromInput(input)

	// We initialise the subject and the observers
	subject := Subject.NewConcreteSubject("endpoints.json")

	// We attach the observers to the subject in order to receive the updates
	for _, observer := range observers {
		subject.Attach(observer)
	}

	// Start listening to the websockets
	go subject.StartListening()
}
