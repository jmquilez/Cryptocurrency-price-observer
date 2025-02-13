package Observer

// Observer interface defines the methods that all observers must implement
type Observer interface {
	// Update is called when the subject updates its data
	Update(data string)
	// GetID returns the unique identifier for the observer
	GetID() string
}
