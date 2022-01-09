package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main(){

	stock := generateStock()
	
	for {
		updateStockPrice(&stock)
		fmt.Println(stock.name, stock.price)
		time.Sleep(time.Second)
	}
}

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

func generateStock() Stock{
	name := generateStockName()
	s := Stock{name, 42.1}
	return s
}

func generateStockName() string {
	return "$ABC"
}

func generateNextPrice(currentPrice float64 , direction Direction ) float64{
	rand.Seed(time.Now().UnixNano())
	random := rand.Float64()
	percentChanged := (random * 10 - 2.5) / 100
	change := currentPrice * percentChanged
	return currentPrice + change
}

func updateStockPrice(stock *Stock){
	newPrice := generateNextPrice(stock.price, up)
	stock.price = newPrice
}