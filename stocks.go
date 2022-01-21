package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

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

func generateRandomStock() Stock{
	name := generateStockName()
	initialPrice := generateInitialPrice()
	s := Stock{name, initialPrice}
	return s
}

func generateStocks(num  int, filepath string) map[string]*Stock {

	// tickers are parsed from the file specified.
	if (filepath != ""){
		parsed := parseTickersFromJson(filepath)
		return parsed
	}

	stocks := make(map[string]*Stock)	
	for i := 0; i < num; i ++ {
		stock := generateRandomStock()
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
	if (verbose){
		fmt.Println(stock.Name, stock.Price)
	}
	updateStockPrice(stock)
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

func parseTickersFromJson(filepath string) map[string]*Stock{
	jsonFile, err := os.Open(filepath)
	
	if (err != nil){
		// return nil, errors.New("JSON file at `" + filepath + "` cannot be parsed.")
		panic("JSON file at `" + filepath + "` cannot be parsed.")
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var parsedStocks []Stock
	json.Unmarshal(byteValue, &parsedStocks)
	
	stockMap := make(map[string]*Stock)

	for _, stock := range(parsedStocks){
		s := Stock{Name: stock.Name, Price: stock.Price}
		stockMap[s.Name] = &s
	}
	
	return stockMap
}