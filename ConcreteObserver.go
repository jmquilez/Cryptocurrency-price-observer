package main

import "fmt"

type ConcreteObserver struct {
	id string
}

func (p *ConcreteObserver) Update(data string) {
	fmt.Println("Received data:", data)
}

func (p *ConcreteObserver) GetID() string {
	return p.id
}
