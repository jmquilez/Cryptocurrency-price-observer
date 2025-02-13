package main

type ConcreteSubject struct {
	Observers []Observer
}

func (p *ConcreteSubject) Attach(observer Observer) (bool, error) {
    p.Observers = append(p.Observers, observer)
    return true, nil
}

func (p *ConcreteSubject) Detach(observer Observer) (bool, error) {
    for i, o := range p.Observers {
        if o.GetID() == observer.GetID() {
            p.Observers = append(p.Observers[:i], p.Observers[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

func (p *ConcreteSubject) Notify(data string) (bool, error) {
	for _, observer := range p.Observers {
		observer.Update(data)
	}
	return true, nil
}
