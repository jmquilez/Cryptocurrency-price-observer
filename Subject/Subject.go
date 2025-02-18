// Credits: Juan José Serrano, José Miguel Quilez.
// Package Subject defines the Subject interface.
package Subject

import "p1/Observer"

// Subject interface defines the methods that all subjects must implement
type Subject interface {
	// Observers returns the list of observers attached to the subject
	Observers() []Observer.Observer

	// Attach adds an observer to the subject
	Attach(observer Observer.Observer) (bool, error)

	// Detach removes an observer from the subject
	Detach(observer Observer.Observer) (bool, error)

	// Notify updates all observers with the given data
	Notify(data string) (bool, error)
}

