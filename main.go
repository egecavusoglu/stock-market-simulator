package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var stocks map[string]*Stock
var upgrader = websocket.Upgrader{}
var timeTicker  *time.Ticker

func main(){
	// Initial setup
	rand.Seed(time.Now().UnixNano())
	timeTicker = time.NewTicker(time.Second)
	stocks = generateStocks(4) // generate the stock market tickers
	
	
	// Handle tickers in a go routine each second
	go registerTicker(timeTicker, false)

	// Register server endpoints
	router := mux.NewRouter()
	router.HandleFunc("/stocks", getStocksHandler) // REST: Get all stocks in a REST endpoint
	router.HandleFunc("/stocks/live", publishStocksHandler) // WS: Publish stock data on websocket
	router.HandleFunc("/stocks/{name}", getStockHandler) // REST: Get specific stock data

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

func getStock(ticker string) (*Stock, error){
	ticker = strings.ToUpper(ticker)
	stock := stocks[ticker]
	if (stock != nil){
		return stock, nil
	}
	stock = stocks["$"+ticker]
	if (stock != nil){
		return stock, nil
	}
	return nil, errors.New("ticker $" + ticker + " not found")
}

func tick(stock *Stock, verbose bool){
	updateStockPrice(stock)
	if (verbose){
		fmt.Println(stock.Name, stock.Price)
	}
}

func registerTicker(ticker *time.Ticker, verbose bool){
	for ; true; <-ticker.C {
		for _, stock := range(stocks){
			tick(stock, verbose)
		}
		if (verbose){
			fmt.Println()
		}
	}	
}

func getStocksHandler(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)

}

func getStockHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	ticker := vars["name"]
	stock, err := getStock(ticker)
	if (err != nil){
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

func publishStocksHandler(w http.ResponseWriter, r *http.Request){
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		parsedMessage := string(message)
		
		if parsedMessage == "" { // send all stock data thru websocket
			encodedStocks, _ := json.Marshal(stocks)
			c.WriteMessage(mt, encodedStocks)
			continue
		}
		
		stock, stockErr := getStock(parsedMessage)
		if stockErr != nil {
			c.WriteMessage(mt, []byte(stockErr.Error()))
			continue	
		}
		encodedStock, _ := json.Marshal(stock)
		err = c.WriteMessage(mt, encodedStock)

		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}