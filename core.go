package main

import (
	"fmt"
	"log"
	"time"
)

func runCore() {
	fmt.Println("runCore() starts")
	initDb()
	// return
	//1. get a list of stock to monitor
	// symbols := []string{"AAPL", "BA"}
	// symbols := []string{"BA"}
	// symbols := []string{"UNH", "TRV", "MCD", "AXP", "MSFT", "CAT", "UTX", "IBM", "BA", "MMM", "DIS", "NKE", "VZ", "KO", "AAPL", "PG", "CSCO", "WMT", "INTC", "V", "CVX", "XOM", "HD", "JNJ", "JPM", "WBA", "DOW", "GS", "MRK", "PFE", "DIA", "SPY", "QQQ"}
	symbols1 := []string{"UNH", "TRV", "MCD", "AXP", "MSFT"}
	symbols2 := []string{"CAT", "UTX", "IBM", "BA", "MMM"}
	symbols3 := []string{"DIS", "NKE", "VZ", "KO", "AAPL"}
	symbols4 := []string{"PG", "CSCO", "WMT", "INTC", "V"}
	symbols5 := []string{"CVX", "XOM", "HD", "JNJ", "JPM"}
	symbols6 := []string{"WBA", "DOW", "GS", "MRK", "PFE"}
	symbols0 := []string{"DIA", "SPY", "QQQ"}
	//2. activate option and stock data getter
	for range time.Tick(time.Minute * 3) {
		go runOptionAndStockData(symbols1)
		go runOptionAndStockData(symbols2)
		go runOptionAndStockData(symbols3)
		go runOptionAndStockData(symbols4)
		go runOptionAndStockData(symbols5)
		go runOptionAndStockData(symbols6)
		go runOptionAndStockData(symbols0)
	}
	/*3. terminate if option data getter
	and stock data getter are both terminated
	*/
	fmt.Println("runCore() ends")
}

func runOptionAndStockData(symbols []string) {
	for _, symbol := range symbols {
		fmt.Printf("Run option and stock data for symbol %s\n", symbol)
		options, stock, _, err := new(YahooApiManager).GetOptionsAndStockDataBySymbol(symbol)
		if err != nil {
			continue
		}
		//store the option data into database
		for _, option := range options {
			err2 := new(TblOptionData).InsertOrUpdateOptionData(option)
			if err2 != nil {
				panic(err2)
			}
		}
		//store the stock data into database
		err = new(TblStockData).InsertOrUpdateStockData(stock)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Run option report")
	runOptionReport(symbols)
}

func runOptionReport(symbols []string) {
	//get an array of closet expiration date
	expDates := make([]int64, len(symbols))
	for index, symbol := range symbols {
		_, _, exp, err := new(YahooApiManager).GetOptionsAndStockDataBySymbol(symbol)
		if err != nil {
			log.Fatal(err)
		}
		expDates[index] = exp[0]
	}
	//run report every minute
	//for range time.Tick(time.Minute * 1) {
	go runOptionReportDetail(symbols, expDates)
	//}
}

func runOptionReportDetail(symbols []string, expDates []int64) {
	for index, symbol := range symbols {
		new(TblOptionReport).GenerateReport(symbol, expDates[index])
	}
}
