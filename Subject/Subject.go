package Subject

import "Observer/Observer"

type Subject interface {
	Observers() []Observer.Observer

	Attach(observer Observer.Observer) (bool, error)
	Detach(observer Observer.Observer) (bool, error)
	Notify(data string) (bool, error)
}

