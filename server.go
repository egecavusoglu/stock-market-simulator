package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


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