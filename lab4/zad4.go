package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	ID           int
	CustomerName string
	Items        []string
	TotalAmount  float64
}

type ProcessResult struct {
	OrderID      int
	CustomerName string
	Success      bool
	ProcessTime  time.Duration
	Error        error
}


const (
	numOrders       = 20
	numWorkers      = 5
	maxRetries      = 3
	orderIntervalMs = 300
)

var customerNames = []string{"Anna", "Bartek", "Celina", "Damian", "Emilia"}
var itemList = []string{"Laptop", "Myszka", "Monitor", "Klawiatura", "Słuchawki"}


func generateOrders(orderCh chan<- Order) {
	for i := 1; i <= numOrders; i++ {
		time.Sleep(time.Duration(rand.Intn(orderIntervalMs)) * time.Millisecond)
		order := Order{
			ID:           i,
			CustomerName: customerNames[rand.Intn(len(customerNames))],
			Items:        []string{itemList[rand.Intn(len(itemList))]},
			TotalAmount:  float64(rand.Intn(900)+100) + rand.Float64(),
		}
		fmt.Printf("Wygenerowano zamówienie: %+v\n", order)
		orderCh <- order
	}
	close(orderCh)
}

func processOrder(order Order) ProcessResult {
	start := time.Now()
	sleepDuration := time.Duration(rand.Intn(1000)+500) * time.Millisecond
	time.Sleep(sleepDuration)

	success := rand.Float64() > 0.2 
	var err error
	if !success {
		err = fmt.Errorf("błąd przetwarzania zamówienia %d", order.ID)
	}
	return ProcessResult{
		OrderID:      order.ID,
		CustomerName: order.CustomerName,
		Success:      success,
		ProcessTime:  time.Since(start),
		Error:        err,
	}
}

func worker(id int, orderCh <-chan Order, resultCh chan<- ProcessResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range orderCh {
		var result ProcessResult
		for attempt := 1; attempt <= maxRetries; attempt++ {
			result = processOrder(order)
			if result.Success {
				break
			}
			fmt.Printf("Ponowna próba [%d] dla zamówienia %d\n", attempt, order.ID)
		}
		resultCh <- result
	}
}


func main() {
	rand.Seed(time.Now().UnixNano())

	orderCh := make(chan Order)
	resultCh := make(chan ProcessResult)

	var wg sync.WaitGroup

	go generateOrders(orderCh)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, orderCh, resultCh, &wg)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var total, success, failure int
	for result := range resultCh {
		total++
		if result.Success {
			success++
			fmt.Printf("Sukces [%d] %s (%v)\n", result.OrderID, result.CustomerName, result.ProcessTime)
		} else {
			failure++
			fmt.Printf("Niepowodzenie [%d] %s: %v\n", result.OrderID, result.CustomerName, result.Error)
		}
	}

	fmt.Println("\nStatystyki:")
	fmt.Printf("Łącznie zamówień: %d\n", total)
	fmt.Printf("Udane: %d (%.2f%%)\n", success, float64(success)/float64(total)*100)
	fmt.Printf("Nieudane: %d (%.2f%%)\n", failure, float64(failure)/float64(total)*100)
}