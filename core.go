package main

import (
	"fmt"
)

func runCore() {
	fmt.Println("runCore() starts")
	initDb()
	//1. get a list of stock to monitor
	symbols := []string{"AAPL", "BA"}
	//2. activate option and stock data getter
	for _, symbol := range symbols {
		options, stock, err := GetOptionsAndStockDataBySymbol(symbol)
		if err != nil {
			panic(err)
		}
		for _, option := range options {
			err2 := InsertOrUpdateOptionData(option)
			if err2 != nil {
				panic(err2)
			}
		}
		err = InsertOrUpdateStockData(stock)
		if err != nil {
			panic(err)
		}
	}
	/*3. terminate if option data getter
	and stock data getter are both terminated
	*/
	fmt.Println("runCore() ends")
}
