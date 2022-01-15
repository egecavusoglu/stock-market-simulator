# Stock Market Simulator

Launch your stock exchange server and API in your local. Develop your trading algorithms with ease.

You can consume the stock market HTTP server via REST and websocket.

1. `GET /stocks` : Returns all stocks in the market
2. `GET /stocks/{symbol}` : Return specific ticker data.
3. `WS /stocks/live` : Open websocket connection to server.
   - Send empty message to get all stocks data.
   - Send ticker symbol (eg. `$ABC`) to get specific stock data.

Project Roadmap:

- [x] CLI: Verbose flag to print tickers
- [x] CLI: Get number of tickers
- [ ] CLI: Get initial tickers from json (with initial ticker names and initial price) Stock {"$NAME", 120.0}
- [x] Emit data on websocket server
- [x] Rest Endpoint: give individual stock data
- [x] Rest Endpoint: give all stock names
- [x] Configure CI/CD Release
- [ ] Deploy as Homebrew tap
