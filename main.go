package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// globals
var stocks map[string]*Stock
var upgrader = websocket.Upgrader{}
var timeTicker  *time.Ticker

func main(){
	writeWelcomeMessage()
	
	// Program flags
	seedPtr := flag.Int64("seed", time.Now().UnixNano(), "Initialize pseudo random sequence with a specific seed. Defaults to unix time.")
	verbosePtr := flag.Bool("verbose", false, "In verbose mode, tick events will be printed out to standart output. Defaults to false.")
	tickerCountPtr := flag.Int("count", 3, "Number of stock tickers you would like to generate. Defaults to 3.")
	portPtr := flag.Int("port", 8080, "Port number http server is launched on. Defaults to 8080.")
	flag.Parse()

	// Initial setup
	rand.Seed(*seedPtr)
	timeTicker = time.NewTicker(time.Second)
	stocks = generateStocks(*tickerCountPtr) // generate the stock market tickers
	
	// Handle tickers in a go routine each second
	go registerTicker(timeTicker, *verbosePtr)

	// Register server endpoints
	port := strconv.Itoa(*portPtr)
	router := mux.NewRouter()
	router.HandleFunc("/", getStocksHandler) // REST: Get all stocks in a REST endpoint
	router.HandleFunc("/stocks", getStocksHandler) // REST: Get all stocks in a REST endpoint
	router.HandleFunc("/stocks/live", publishStocksHandler) // WS: Get stocks data over websocket
	router.HandleFunc("/stocks/{name}", getStockHandler) // REST: Get specific stock data
	color.Green("Server is live on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
	
}


func writeWelcomeMessage(){
	fmt.Println(
	color.CyanString("\nWelcome to stock-market-simulator!"),
	`

HTTP Server endpoints:
	GET / : Get all stocks data
	GET /stocks : Get all stocks data
	GET /stocks/{symbol} : Get specific stock's data
	WEBSOCKET /stocks/live : Get stocks data over websocket.
	`)
}

