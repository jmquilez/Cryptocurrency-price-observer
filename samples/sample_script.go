package sample

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/websocket"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// AggTrade represents the structure of Binance's aggTrade stream data
type AggTrade struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
}

var prices plotter.XYs

func main() {
	// Set up WebSocket connection
	u := url.URL{Scheme: "wss", Host: "stream.binance.com:9443", Path: "/ws/btcusdt@aggTrade"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer c.Close()

	// Handle graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Read messages from WebSocket
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}

			// Parse aggTrade data
			var trade AggTrade
			if err := json.Unmarshal(message, &trade); err != nil {
				log.Println("JSON unmarshal error:", err)
				continue
			}

			// Convert price to float64
			price, err := strconv.ParseFloat(trade.Price, 64)
			if err != nil {
				log.Println("Price parse error:", err)
				continue
			}

			// Append price to the slice
			prices = append(prices, plotter.XY{X: float64(len(prices)), Y: price})

			// Update the chart every 10 data points
			if len(prices)%10 == 0 {
				plotChart()
			}
		}
	}()

	// Wait for interrupt signal
	<-interrupt
	log.Println("Shutting down...")
}

// plotChart creates and saves a line chart
func plotChart() {
	p := plot.New()

	p.Title.Text = "BTC/USDT Price Chart"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Price"

	// Calculate min and max prices for better scaling
	minPrice := prices[0].Y
	maxPrice := prices[0].Y
	for _, point := range prices {
		if point.Y < minPrice {
			minPrice = point.Y
		}
		if point.Y > maxPrice {
			maxPrice = point.Y
		}
	}

	// Add some padding to the price range (0.1% of the range)
	padding := (maxPrice - minPrice) * 0.001
	p.Y.Min = minPrice - padding
	p.Y.Max = maxPrice + padding

	// Set X axis range
	p.X.Min = 0
	p.X.Max = float64(len(prices))

	// Create a line plot
	line, err := plotter.NewLine(prices)
	if err != nil {
		log.Fatal("Line plot error:", err)
	}
	line.Color = plotutil.Color(0)

	p.Add(line)

	// Save the plot to a PNG file
	if err := p.Save(10*vg.Inch, 6*vg.Inch, "price_chart.png"); err != nil {
		log.Fatal("Save plot error:", err)
	}

	log.Println("Chart updated and saved to price_chart.png")
}
