# Stock Market Simulator

Launch your stock exchange server and API in your local. Develop your trading algorithms with ease.

This binary runs a simulated stock market based on randomized numbers, updating each ticker every second and exposing the data through a local web server. See [usage](#usage) and [consuming](#consuming-the-api) sections to get started.

## Usage

1. Download the binary from the releases tab.
   - You may need to change security & privacy settings of your computer to allow running unverified binaries. In macOS, go to `System Preferences > Security & Privacy` to enable running the program.
2. Run executable.

### Program Flags

- `seed` (Int64): Initialize pseudo random sequence with a specific seed. Defaults to unix time.
- `verbose` (Bool): In verbose mode, tick events will be printed out to standart output. Defaults to false.
- `count` (Int): Number of stock tickers you would like to generate. Defaults to 3.
- `port` (Int): Port number http server is launched on. Defaults to 8080.

Example usage

```console
./stock-market-simulator -seed 2 -verbose=true -count 12 -port 8081
```

## Consuming the API

You can consume the stock market HTTP server via REST and websocket.

1. `GET /stocks` : Returns all stocks in the market
2. `GET /stocks/{symbol}` : Return specific ticker data.
3. `WS /stocks/live` : Open websocket connection to server.
   - Send empty message to get all stocks data.
   - Send ticker symbol (eg. `$ABC`) to get specific stock data.

## Project Roadmap:

- [x] CLI: Verbose flag to print tickers
- [x] CLI: Get number of tickers
- [x] Emit data on websocket server
- [x] Rest Endpoint: give individual stock data
- [x] Rest Endpoint: give all stock names
- [x] Configure CI/CD Release
- [ ] CLI: Get initial tickers from json (with initial ticker names and initial price) Stock {"$NAME", 120.0}
- [x] Write `How to use` docs, explain flags.
- [ ] Include contributing guide. (Env setup, dependencies, testing...)
- [ ] Deploy as Homebrew tap
