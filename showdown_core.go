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
