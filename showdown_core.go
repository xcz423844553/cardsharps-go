package main

import (
	"fmt"
	"log"

	tda "github.com/xcz423844553/td_ameritrade_client_golang"
)

type Showdown2 struct {
}

func (sd *Showdown2) runTheShow() error {
	board := new(Board2)
	dealer := new(Dealer2)

	//initiate client to connect to TD Ameritrade
	client, _, _, err := tda.GetClient(authCode, clientID, refreshCode)
	if err != nil {
		log.Fatalln(err)
	}

	executeFunc := func(symbol string) {
		dealer.DownloadAllOptionChainAndStock(client, symbol)
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
