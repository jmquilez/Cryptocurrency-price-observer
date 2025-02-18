// Credits: Juan José Serrano, José Miguel Quilez.
// Package Observer implements the Observer pattern for cryptocurrency price tracking
package Observer

import (
	"fmt"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

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

// NewConcreteObserver creates a new ConcreteObserver
func NewConcreteObserver(id string, Btc_Ok bool, Eth_Ok bool, Ada_Ok bool) *ConcreteObserver {
	return &ConcreteObserver{
		id: id,
		Btc_Ok: Btc_Ok,
		Eth_Ok: Eth_Ok,
		Ada_Ok: Ada_Ok,
	}
}

// Update receives new cryptocurrency prices and updates the observer's values
func (p *ConcreteObserver) Update(Btc, Eth, Ada float64) {
	p.Btc = Btc
	p.Eth = Eth
	p.Ada = Ada

	// Print the graph
	if p.Btc_Ok {
		fmt.Println("BTC: ", Btc)
		go p.PrintGraph("BTC")
	}
	if p.Eth_Ok {
		fmt.Println("ETH: ", Eth)
		go p.PrintGraph("ETH")
	}
	if p.Ada_Ok {
		fmt.Println("ADA: ", Ada)
		go p.PrintGraph("ADA")
	}
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

var btcPrices plotter.XYs

// PrintGraph creates and saves a line chart of cryptocurrency prices
func (p *ConcreteObserver) PrintGraph(option string) {
	// Append current price to the price history
	btcPrices = append(btcPrices, plotter.XY{X: float64(len(btcPrices)), Y: p.Btc})

	// Create a new plot
	plt := plot.New()

	plt.Title.Text = fmt.Sprintf("%s Price Chart - Observer %s", option, p.id)
	plt.X.Label.Text = "Time"
	plt.Y.Label.Text = "Price (USDT)"

	// Calculate min and max prices for better scaling
	minPrice := btcPrices[0].Y
	maxPrice := btcPrices[0].Y
	for _, point := range btcPrices {
		if point.Y < minPrice {
			minPrice = point.Y
		}
		if point.Y > maxPrice {
			maxPrice = point.Y
		}
	}

	// Add padding to the price range (0.1% of the range)
	padding := (maxPrice - minPrice) * 0.001
	plt.Y.Min = minPrice - padding
	plt.Y.Max = maxPrice + padding

	// Set X axis range
	plt.X.Min = 0
	plt.X.Max = float64(len(btcPrices) * 3)

	// Create and add the line plot
	line, err := plotter.NewLine(btcPrices)
	if err != nil {
		log.Printf("Error creating line plot for observer %s: %v", p.id, err)
		return
	}
	line.Color = plotutil.Color(0)
	plt.Add(line)

	// Ensure the Results directory exists
	if err := os.MkdirAll("Results", 0755); err != nil {
		log.Printf("Error creating Results directory: %v", err)
		return
	}

	// Save the plot to a PNG file
	filename := fmt.Sprintf("Results/%s_chart_%s.png", option, p.id)
	if err := plt.Save(10*vg.Inch, 6*vg.Inch, filename); err != nil {
		log.Printf("Error saving plot for observer %s: %v", p.id, err)
		return
	}

	log.Printf("Chart updated and saved to %s", filename)
}