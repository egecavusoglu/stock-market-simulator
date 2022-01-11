package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main(){
	rand.Seed(time.Now().UnixNano())

	stocks := generateStocks(10) // generate the stock market tickers
	
	for { // timer loop
		for _, stock := range(stocks){
		updateStockPrice(&stock)
		fmt.Println(stock.name, stock.price)
		}
		fmt.Println()
		time.Sleep(time.Second)
	}
}

// Types and enums
type Stock struct{
	name string
	price float64
	// priceHistory []float64
}

type Direction int
const  (
	up Direction = iota
	down
)

// Functions

func generateStock() Stock{
	name := generateStockName()
	s := Stock{name, 42.1}
	return s
}

func generateStocks(num  int) []Stock {	
	var stocks []Stock
	for i := 0; i < num; i ++ {
		stock := generateStock()
		stocks = append(stocks, stock)
	}
	return stocks
}

func generateRandomString(length int) string {
	var name string
	for i:= 0; i < length; i++ {
		char := 'A' + rune(rand.Intn(26))
		name = name + string(char)
	}
	return name
}

func generateStockName() string {
	return "$" + generateRandomString(3)
}

func generateNextPrice(currentPrice float64 , direction Direction ) float64{
	random := rand.Float64()
	percentChanged := (random * 10 - 2.5) / 100
	change := currentPrice * percentChanged
	return currentPrice + change
}

func updateStockPrice(stock *Stock){
	newPrice := generateNextPrice(stock.price, up)
	stock.price = newPrice
}