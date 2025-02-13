package main

type Subject interface {
	Observers() []Observer

	Attach(observer Observer) (bool, error)
	Detach(observer Observer) (bool, error)
	Notify(data string) (bool, error)
}

