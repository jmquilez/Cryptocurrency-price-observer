package main

import (
	"Observer/Observer"
	"Subject/Subject"
	"fmt"
)

func initialiseConcreteImplementors() (Subject.Subject, Observer.Observer, Observer.Observer, Observer.Observer) {
	subject := Subject.NewConcreteSubject()
	observer1 := Observer.NewConcreteObserver("observer1", true, true, false)
	observer2 := Observer.NewConcreteObserver("observer2", true, false, true)
	observer3 := Observer.NewConcreteObserver("observer3", true, true, false)
	return subject, observer1, observer2, observer3
}

func main() {
	// We initialise the subject and the observers
	subject, observer1, observer2, observer3 := initialiseConcreteImplementors()
	
	// We attach the observers to the subject in order to receive the updates
	subject.Attach(observer1)
	subject.Attach(observer2)
	subject.Attach(observer3)

	// Start listening to the websockets
	go subject.StartListening()

	// When the enter key is pressed, the program ends
	fmt.Scanln()
}
