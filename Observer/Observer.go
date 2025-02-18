// Credits: Juan José Serrano, José Miguel Quilez.
// Package Observer defines the Observer interface.
package Observer

// Observer interface defines the methods that all observers must implement
type Observer interface {
	// Update is called when the subject updates its data
	Update(Btc, Eth, Ada float64)
	// GetID returns the unique identifier for the observer
	GetID() string
	// GetBtc returns the current Bitcoin price
	GetBtc() float64
	// GetEth returns the current Ethereum price
	GetEth() float64
	// GetAda returns the current Cardano price
	GetAda() float64
	// GetBtc_Ok returns whether the observer is subscribed to Bitcoin updates
	GetBtc_Ok() bool
	// GetEth_Ok returns whether the observer is subscribed to Ethereum updates
	GetEth_Ok() bool
	// GetAda_Ok returns whether the observer is subscribed to Cardano updates
	GetAda_Ok() bool
}
