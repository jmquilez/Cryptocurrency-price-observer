package Observer

type Observer interface {
	Update(data string)
	GetID() string
}
