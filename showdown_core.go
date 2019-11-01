package main

import "fmt"

type Showdown2 struct {
}

func (sd *Showdown2) runTheShow() error {
	board := new(Board2)
	dealer := new(Dealer2)
	executeFunc := func(symbol string) {
		dealer.DownloadAllOptionChainAndStock(symbol)
	}
	callback := func() {
		fmt.Println("Showdown is completed.")
	}
	boardErr := board.StartGame(board.GetSymbolTagAll(), executeFunc, callback)
	return boardErr
}

func (sd *Showdown2) isMarketOpen() bool {
	quote, err := new(YahooApi).GetQuote("SPY")
	if err != nil {
		return false
	}
	return quote.isMarketOpen()
}

func (sd *Showdown2) isMarketPreOpenOrOpen() bool {
	quote, err := new(YahooApi).GetQuote("SPY")
	if err != nil {
		return false
	}
	return quote.isMarketPreOpen()
}
