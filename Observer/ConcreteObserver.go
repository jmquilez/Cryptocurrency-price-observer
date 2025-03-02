// Credits: Juan José Serrano, José Miguel Quilez.
// Package Observer implements the Observer pattern for cryptocurrency price tracking
package Observer

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// ConcreteObserver represents a specific observer that tracks cryptocurrency prices
type ConcreteObserver struct {
	id             string  // Unique identifier for the observer
	Btc            float64 // Bitcoin price
	Eth            float64 // Ethereum price
	Ada            float64 // Cardano price
	Btc_Ok         bool    // Flag indicating if observer is subscribed to BTC updates
	Eth_Ok         bool    // Flag indicating if observer is subscribed to ETH updates
	Ada_Ok         bool    // Flag indicating if observer is subscribed to ADA updates
	activeFileOps  *sync.WaitGroup
	shutdownMutex  *sync.Mutex
	isShuttingDown *bool
}

// NewConcreteObserver creates a new ConcreteObserver
func NewConcreteObserver(id string, Btc_Ok bool, Eth_Ok bool, Ada_Ok bool, activeFileOps *sync.WaitGroup, shutdownMutex *sync.Mutex, isShuttingDown *bool) *ConcreteObserver {
	return &ConcreteObserver{
		id:             id,
		Btc_Ok:         Btc_Ok,
		Eth_Ok:         Eth_Ok,
		Ada_Ok:         Ada_Ok,
		activeFileOps:  activeFileOps,
		shutdownMutex:  shutdownMutex,
		isShuttingDown: isShuttingDown,
	}
}

// Update receives new cryptocurrency prices and updates the observer's values
func (p *ConcreteObserver) Update(Btc, Eth, Ada float64) {
	if Btc >= 0 {
		p.Btc = Btc
	}
	if Eth >= 0 {
		p.Eth = Eth
	}
	if Ada >= 0 {
		p.Ada = Ada
	}

	// Print the graph
	if p.Btc_Ok && Btc >= 0 {
		fmt.Println("BTC: ", Btc)
		go p.PrintGraph("BTC", p.Btc)
	}
	if p.Eth_Ok && Eth >= 0 {
		fmt.Println("ETH: ", Eth)
		go p.PrintGraph("ETH", p.Eth)
	}
	if p.Ada_Ok && Ada >= 0 {
		fmt.Println("ADA: ", Ada)
		go p.PrintGraph("ADA", p.Ada)
	}
}

// GetID returns the observer's unique identifier
func (p *ConcreteObserver) GetID() string {
	return p.id
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

var priceHistories = map[string]plotter.XYs{}
var mutex sync.Mutex

// PrintGraph creates and saves a line chart of cryptocurrency prices
func (p *ConcreteObserver) PrintGraph(option string, price float64) {
	// Check if we're shutting down
	p.shutdownMutex.Lock()
	if *p.isShuttingDown {
		p.shutdownMutex.Unlock()
		return // Skip file operations during shutdown
	}
	p.shutdownMutex.Unlock()

	// Indicate we're starting a file operation
	p.activeFileOps.Add(1)
	defer p.activeFileOps.Done()

	// Use mutex to prevent concurrent map writes
	mutex.Lock()
	// Append current price to the price history
	priceHistories[option] = append(priceHistories[option], plotter.XY{X: float64(len(priceHistories[option])), Y: price})

	// Create a local copy of the price history for this option
	localPriceHistory := make(plotter.XYs, len(priceHistories[option]))
	copy(localPriceHistory, priceHistories[option])
	mutex.Unlock()

	// Create a new plot
	plt := plot.New()

	plt.Title.Text = fmt.Sprintf("%s Price Chart - Observer %s", option, p.id)
	plt.X.Label.Text = "Time (samples)"
	plt.Y.Label.Text = "Price (USDT)"

	// Calculate min and max prices for better scaling
	mutex.Lock()
	minPrice := localPriceHistory[0].Y
	maxPrice := localPriceHistory[0].Y
	for _, point := range localPriceHistory {
		if point.Y < minPrice {
			minPrice = point.Y
		}
		if point.Y > maxPrice {
			maxPrice = point.Y
		}
	}
	mutex.Unlock()

	// Add padding to the price range (0.1% of the range)
	padding := (maxPrice - minPrice) * 0.001
	plt.Y.Min = minPrice - padding
	plt.Y.Max = maxPrice + padding

	// Set X axis range
	plt.X.Min = 0
	mutex.Lock()
	plt.X.Max = float64(len(localPriceHistory) * 3)

	// Create and add the line plot
	line, err := plotter.NewLine(localPriceHistory)
	mutex.Unlock()
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

	// Save the plot directly to the final filename
	filename := fmt.Sprintf("Results/%s_chart_%s.png", option, p.id)
	if err := plt.Save(10*vg.Inch, 6*vg.Inch, filename); err != nil {
		log.Printf("Error saving plot for observer %s: %v", p.id, err)
		return
	}
}
