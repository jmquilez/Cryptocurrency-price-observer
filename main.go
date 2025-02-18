package main

import (
	"bufio"
	"fmt"
	"os"
	"p1/Observer"
	"p1/Subject"
	"strings"
)

// parseObserver transforma una cadena en un Observer.Observer.
// La cadena debe tener el formato: "observerID,CURRENCY1,CURRENCY2,..."
// Ejemplo: "observer1,BTC,ETH,ADA"
func parseObserver(input string) Observer.Observer {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}
	tokens := strings.Split(input, ",")
	if len(tokens) < 1 {
		return nil
	}
	id := strings.TrimSpace(tokens[0])
	btcOk, ethOk, adaOk := false, false, false
	for _, token := range tokens[1:] {
		token = strings.TrimSpace(strings.ToUpper(token))
		if token == "BTC" {
			btcOk = true
		} else if token == "ETH" {
			ethOk = true
		} else if token == "ADA" {
			adaOk = true
		}
	}
	return Observer.NewConcreteObserver(id, btcOk, ethOk, adaOk)
}

// observersFromInput procesa la entrada y devuelve un slice de Observer.Observer.
// Primero separa la entrada por ";" (cada observador) y luego procesa cada uno.
func observersFromInput(input string) []Observer.Observer {
	parts := strings.Split(input, ";")
	observers := []Observer.Observer{}
	for _, part := range parts {
		if obs := parseObserver(part); obs != nil {
			observers = append(observers, obs)
		}
	}
	return observers
}

func main() {
	// Crear un lector de entrada
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please specify the observers and their preferences in the format as done in this example:")
	fmt.Println("observer1,BTC,ETH,ADA; observer2,BTC,ETH; observer3,ETH,ADA;")
	input, _ := reader.ReadString('\n')

	// Procesar la entrada para crear los observadores
	observers := observersFromInput(input)

	// Inicializar el subject; NewConcreteSubject ya no recibe argumentos.
	subject := Subject.NewConcreteSubject()

	// Adjuntar los observadores al subject para recibir actualizaciones
	for _, observer := range observers {
		subject.Attach(observer)
	}

	// Iniciar la escucha de websockets en una goroutine
	go subject.StartListening()

	// Mantener el programa en ejecución hasta que se presione Enter
	fmt.Println("Press Enter to exit...")
	reader.ReadString('\n')
}
