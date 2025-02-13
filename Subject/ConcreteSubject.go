// Package Subject implements the Subject part of the Observer pattern
package Subject

import "Observer/Observer"

// ConcreteSubject represents the subject that maintains a list of observers
// and notifies them of cryptocurrency price changes
type ConcreteSubject struct {
	Observers []Observer.Observer
}

// Attach adds a new observer to the notification list
// Returns true if successful and nil error
func (p *ConcreteSubject) Attach(observer Observer.Observer) (bool, error) {
	p.Observers = append(p.Observers, observer)
	return true, nil
}

// Detach removes an observer from the notification list
// Returns true if observer was found and removed, false if not found
func (p *ConcreteSubject) Detach(observer Observer.Observer) (bool, error) {
	for i, o := range p.Observers {
		if o.GetID() == observer.GetID() {
			p.Observers = append(p.Observers[:i], p.Observers[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

// Notify updates all relevant observers with new cryptocurrency prices
// A value of -1 indicates no update for that cryptocurrency
// Only notifies observers that are subscribed to the updated currencies
func (p *ConcreteSubject) Notify(Btc, Eth, Ada float64) (bool, error) {
	// If the given data is not -1, then we intended to give that data
	Btc_Ok := Btc != -1
	Eth_Ok := Eth != -1
	Ada_Ok := Ada != -1

	// We only notify the observers that are "subscribed" to the given data
	for _, observer := range p.Observers {
		if Btc_Ok && observer.GetBtc_Ok() || Eth_Ok && observer.GetEth_Ok() || Ada_Ok && observer.GetAda_Ok() {
			observer.Update(Btc, Eth, Ada)
		}
	}
	return true, nil
}
