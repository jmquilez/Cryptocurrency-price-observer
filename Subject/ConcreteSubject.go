// Package Subject implements the Subject part of the Observer pattern
package Subject

import (
	"Observer/Observer"
	"encoding/json"
	"log"

	"golang.org/x/net/websocket"
)

// ConcreteSubject represents the subject that maintains a list of observers
// and notifies them of cryptocurrency price changes
type ConcreteSubject struct {
	// List of observers that are notified of price changes
	Observers []Observer.Observer

	// Prices of the cryptocurrencies
	Btc_Price float64
	Eth_Price float64
	Ada_Price float64

	// Websocket endpoints for the cryptocurrencies
	Btc_Socket string
	Eth_Socket string
	Ada_Socket string
}

// Attach adds a new observer to the notification list
// Returns true if successful and nil error
func (p *ConcreteSubject) Attach(observer Observer.Observer) (bool, error) {
	p.Observers = append(p.Observers, observer)
	return true, nil
}

// Detach removes an observer from the notification list
// Returns true if observer was found and removed, false if not found
func (p *ConcreteSubject) Detach(observer Observer.Observer) (bool, error) {
	for i, o := range p.Observers {
		if o.GetID() == observer.GetID() {
			p.Observers = append(p.Observers[:i], p.Observers[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

// Notify updates all relevant observers with new cryptocurrency prices
// A value of -1 indicates no update for that cryptocurrency
// Only notifies observers that are subscribed to the updated currencies
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

// GetBtc_Price returns the current Bitcoin price
func (p *ConcreteSubject) GetBtc_Price() float64 {
	return p.Btc_Price
}

// GetEth_Price returns the current Ethereum price
func (p *ConcreteSubject) GetEth_Price() float64 {
	return p.Eth_Price
}

// GetAda_Price returns the current Cardano price
func (p *ConcreteSubject) GetAda_Price() float64 {
	return p.Ada_Price
}

// GetBtc_Socket returns the websocket endpoint for Bitcoin
func (p *ConcreteSubject) GetBtc_PriceFromSocket() float64 {
	
	// Connect to the websocket
	conn, err := websocket.Dial(p.Btc_Socket, "", "http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Read the message from the websocket
	var msg = make([]byte, 512)
	var n int
	if _, err := conn.Read(msg[:n]); err != nil {
		log.Fatal(err)
	}

	// Parse the message as a JSON object
	var data map[string]interface{}
	if err := json.Unmarshal(msg[:n], &data); err != nil {
		log.Fatal(err)
	}

	// Return the price of Bitcoin
	return data["price"].(float64)
}

// GetEth_Socket returns the websocket endpoint for Ethereum
func (p *ConcreteSubject) GetEth_PriceFromSocket() float64 {
	// Connect to the websocket
	conn, err := websocket.Dial(p.Eth_Socket, "", "http://localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Read the message from the websocket
	var msg = make([]byte, 512)
	var n int
	if _, err := conn.Read(msg[:n]); err != nil {
		log.Fatal(err)
	}

	// Parse the message as a JSON object
	var data map[string]interface{}
	if err := json.Unmarshal(msg[:n], &data); err != nil {
		log.Fatal(err)
	}

	return data["price"].(float64)
}

// GetAda_Socket returns the websocket endpoint for Cardano
func (p *ConcreteSubject) GetAda_PriceFromSocket() float64 {
	// Connect to the websocket
	conn, err := websocket.Dial(p.Ada_Socket, "", "http://localhost:8082")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Read the message from the websocket
	var msg = make([]byte, 512)
	var n int
	if _, err := conn.Read(msg[:n]); err != nil {
		log.Fatal(err)
	}

	// Parse the message as a JSON object
	var data map[string]interface{}
	if err := json.Unmarshal(msg[:n], &data); err != nil {
		log.Fatal(err)
	}

	// Return the price of Cardano
	return data["price"].(float64)
}

// listenToCrypto is a helper function that listens to a cryptocurrency websocket and updates the prices
func (p *ConcreteSubject) listenToCrypto(wsURL string, updates chan<- float64) {
	conn, err := websocket.Dial(wsURL, "", "http://localhost")
	if err != nil {
		log.Printf("Error connecting to %s: %v", wsURL, err)
		return
	}
	defer conn.Close()

	for {
		var msg = make([]byte, 512)
		if _, err := conn.Read(msg); err != nil {
			log.Printf("Error reading from %s: %v", wsURL, err)
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Printf("Error parsing JSON from %s: %v", wsURL, err)
			continue
		}

		if price, ok := data["price"].(float64); ok {
			updates <- price
		}
	}
}

// StartListening starts the listening process for cryptocurrency prices.
// Since there are not concurrent updates, we can use a single channel for each cryptocurrency.
func (p *ConcreteSubject) StartListening() {
	// Create go channels for each cryptocurrency
	btcUpdates := make(chan float64)
	ethUpdates := make(chan float64)
	adaUpdates := make(chan float64)

	// Start go routines for each websocket
	go p.listenToCrypto(p.Btc_Socket, btcUpdates)
	go p.listenToCrypto(p.Eth_Socket, ethUpdates)
	go p.listenToCrypto(p.Ada_Socket, adaUpdates)

	// Listen to the updates from each channel
	for {
		select {
		case btcPrice := <-btcUpdates:
			p.Btc_Price = btcPrice
			p.Notify(btcPrice, -1, -1)
		case ethPrice := <-ethUpdates:
			p.Eth_Price = ethPrice
			p.Notify(-1, ethPrice, -1)
		case adaPrice := <-adaUpdates:
			p.Ada_Price = adaPrice
			p.Notify(-1, -1, adaPrice)
		}
	}
}
