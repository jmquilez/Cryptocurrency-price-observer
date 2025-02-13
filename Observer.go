package main

type Observer interface {
	Update(data string)
	GetID() string
}
