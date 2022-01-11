package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var stocks map[string]*Stock

func main(){
	rand.Seed(time.Now().UnixNano())

	stocks = generateStocks(1) // generate the stock market tickers
	
	ticker := time.NewTicker(time.Second)
	go func (){
		for {
			select{
			case <- ticker.C:
				for _, stock := range(stocks){
					tick(stock, true)
				}
				fmt.Println()
			}
			
		}
		
	}()

	// Register server
	router := mux.NewRouter()
	router.HandleFunc("/stocks", getStocks)
	router.HandleFunc("/stock/{name}", getStock)
	log.Fatal(http.ListenAndServe(":8080", router))
	
}

// Types and enums

type Stock struct{
	Name string `json:"name"`
	Price float64 `json:"price"`
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
	initialPrice := generateInitialPrice()
	fmt.Println(initialPrice)
	s := Stock{name, initialPrice}
	return s
}

func generateStocks(num  int) map[string]*Stock {
	stocks := make(map[string]*Stock)	
	for i := 0; i < num; i ++ {
		stock := generateStock()
		stocks[stock.Name] = &stock
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

func generateInitialPrice() float64 {
	return rand.Float64() * 1300
}

func generateNextPrice(currentPrice float64 , direction Direction ) float64{
	const CHANGE_MAGNITUDE = 10
	random := rand.Float64()
	percentChanged := (random * CHANGE_MAGNITUDE - CHANGE_MAGNITUDE/2) / 100
	change := currentPrice * percentChanged
	return currentPrice + change
}

func updateStockPrice(stock *Stock){
	newPrice := generateNextPrice(stock.Price, up)
	stock.Price = newPrice
}

func tick(stock *Stock, verbose bool){
	updateStockPrice(stock)
	if (verbose){
		fmt.Println(stock.Name, stock.Price)
	}
}

func getStocks(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)

}

func getStock(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	ticker := vars["name"]
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	stock := stocks[ticker]
	if (stock != nil){
		json.NewEncoder(w).Encode(stock)
		return 
	}

	stock = stocks["$"+ticker]
	if (stock != nil){
		json.NewEncoder(w).Encode(stock)
		return
	} 
	
	fmt.Fprint(w, "Ticker not found")
	
	
	
}