// Package Observer implements the Observer pattern for cryptocurrency price tracking
package Observer

// ConcreteObserver represents a specific observer that tracks cryptocurrency prices
type ConcreteObserver struct {
	id     string  // Unique identifier for the observer
	Btc    float64 // Bitcoin price
	Eth    float64 // Ethereum price
	Ada    float64 // Cardano price
	Btc_Ok bool    // Flag indicating if observer is subscribed to BTC updates
	Eth_Ok bool    // Flag indicating if observer is subscribed to ETH updates
	Ada_Ok bool    // Flag indicating if observer is subscribed to ADA updates
}

// Update receives new cryptocurrency prices and updates the observer's values
func (p *ConcreteObserver) Update(Btc, Eth, Ada float64) {
	p.Btc = Btc
	p.Eth = Eth
	p.Ada = Ada
}

// GetID returns the observer's unique identifier
func (p *ConcreteObserver) GetID() string {
	return p.id
}

// GetBtc returns the current Bitcoin price
func (p *ConcreteObserver) GetBtc() float64 {
	return p.Btc
}

// GetEth returns the current Ethereum price
func (p *ConcreteObserver) GetEth() float64 {
	return p.Eth
}

// GetAda returns the current Cardano price
func (p *ConcreteObserver) GetAda() float64 {
	return p.Ada
}

// GetBtc_Ok returns whether the observer is subscribed to Bitcoin updates
func (p *ConcreteObserver) GetBtc_Ok() bool {
	return p.Btc_Ok
}

// GetEth_Ok returns whether the observer is subscribed to Ethereum updates
func (p *ConcreteObserver) GetEth_Ok() bool {
	return p.Eth_Ok
}

// GetAda_Ok returns whether the observer is subscribed to Cardano updates
func (p *ConcreteObserver) GetAda_Ok() bool {
	return p.Ada_Ok
}
