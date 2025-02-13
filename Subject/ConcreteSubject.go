package Subject

import "Observer/Observer"

type ConcreteSubject struct {
	Observers []Observer.Observer
}

func (p *ConcreteSubject) Attach(observer Observer.Observer) (bool, error) {
    p.Observers = append(p.Observers, observer)
    return true, nil
}

func (p *ConcreteSubject) Detach(observer Observer.Observer) (bool, error) {
    for i, o := range p.Observers {
        if o.GetID() == observer.GetID() {
            p.Observers = append(p.Observers[:i], p.Observers[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

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
