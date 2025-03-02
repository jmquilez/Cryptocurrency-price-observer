// Credits: Juan José Serrano, José Miguel Quilez.
// Package Subject implements the Subject part of the Observer pattern
package Subject

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"golang.org/x/net/websocket"

	"p1/Observer"
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

// NewConcreteSubject creates a new ConcreteSubject
func NewConcreteSubject() *ConcreteSubject {
	subject := &ConcreteSubject{}
	subject.getEndpoints("endpoints.json")

	return subject
}

// getEndpoints reads the endpoints from the JSON file and assigns them to the corresponding variables
func (p *ConcreteSubject) getEndpoints(fileName string) error {
	// We read the endpoints from the JSON file
	content, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	// We parse the JSON
	var endpoints map[string]string
	if err := json.Unmarshal(content, &endpoints); err != nil {
		return err
	}

	// We assign the endpoints to the corresponding variables
	p.Btc_Socket = endpoints["Btc"]
	p.Eth_Socket = endpoints["Eth"]
	p.Ada_Socket = endpoints["Ada"]

	return nil
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

	// We only notify the observers that are subscribed to the given data
	for _, observer := range p.Observers {
		if (Btc_Ok && observer.GetBtc_Ok()) || (Eth_Ok && observer.GetEth_Ok()) || (Ada_Ok && observer.GetAda_Ok()) {
			observer.Update(Btc, Eth, Ada)
		}
	}
	return true, nil
}

// listenToCrypto es una función auxiliar que escucha un websocket de criptomonedas y envía actualizaciones a través de un canal
func (p *ConcreteSubject) listenToCrypto(wsURL string, updates chan<- float64) {
	conn, err := websocket.Dial(wsURL, "", "http://localhost")
	if err != nil {
		log.Printf("Error conectando a %s: %v", wsURL, err)
		return
	}
	defer conn.Close()

	for {
		msg := make([]byte, 512)
		n, err := conn.Read(msg)
		if err != nil {
			log.Printf("Error leyendo desde %s: %v", wsURL, err)
			continue
		}

		// Se muestra el mensaje recibido para depuración
		log.Printf("Mensaje crudo de %s: %s", wsURL, string(msg[:n]))

		var data map[string]interface{}
		if err := json.Unmarshal(msg[:n], &data); err != nil {
			log.Printf("Error parseando JSON desde %s: %v", wsURL, err)
			continue
		}

		// Binance envía el precio en la clave "p" en formato string
		if priceStr, ok := data["p"].(string); ok {
			price, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				log.Printf("Error convirtiendo el precio desde %s: %v", wsURL, err)
				continue
			}
			updates <- price
		} else {
			log.Printf("No se encontró el precio en los datos provenientes de %s: %v", wsURL, data)
		}
	}
}

// StartListening starts the listening process for cryptocurrency prices.
func (p *ConcreteSubject) StartListening() {
	// Create channels for price updates
	btcUpdates := make(chan float64)
	ethUpdates := make(chan float64)
	adaUpdates := make(chan float64)

	// Start goroutines for each websocket connection
	go p.listenToCrypto(p.Btc_Socket, btcUpdates)
	go p.listenToCrypto(p.Eth_Socket, ethUpdates)
	go p.listenToCrypto(p.Ada_Socket, adaUpdates)

	// Listen for incoming price updates and notify observers accordingly
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
