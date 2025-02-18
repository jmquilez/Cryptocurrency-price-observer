# Cryptocurrency Price Observer

## Authors

#### José Miguel Quilez, Juan José Serrano.

## Description

A Go application that implements the Observer pattern to track real-time cryptocurrency prices from Binance websockets.

## Features

- Real-time price tracking for BTC, ETH and ADA
- Multiple observers can subscribe to different cryptocurrencies
- Automatic price chart generation for each observer
- Concurrent websocket connections using goroutines
- Clean architecture using the Observer design pattern

## Requirements

- Go 1.20+
- Required packages (see go.mod):
  - github.com/gorilla/websocket
  - gonum.org/v1/plot

## Installation

1. Clone the repository
2. Run `go mod download` to install dependencies

## Usage

1. Create a `endpoints.json` file with your Binance websocket endpoints of BTC, ETH or ADA.
2. Run the application:
```bash
go run main.go
```

3. Follow the prompts to create observers with their cryptocurrency preferences
4. Price charts will be automatically generated in PNG format

## Project Structure

- `Observer/` - Observer interface and concrete implementation
- `Subject/` - Subject interface and concrete implementation 
- `Results/` - Results of the price charts
- `endpoints.json` - Websocket endpoint configuration

## Design Pattern

This project implements the Observer pattern where:

- Subject: Maintains cryptocurrency prices and notifies observers
- Observers: Subscribe to price updates for specific cryptocurrencies
- Each observer generates its own price charts
