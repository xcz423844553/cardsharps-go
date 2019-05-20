package main

import (
	"fmt"
	"time"
)

func runCore() {
	fmt.Println("runCore() starts")
	//initDb()
	// return
	//1. get a list of stock to monitor
	symbols, err111 := new(TblSymbol).SelectSymbolByFilter()
	if err111 != nil {
		fmt.Println(err111)
	}
	fmt.Println("started-loop @ " + time.Now().String())
	for i := 0; i < len(symbols); i = i + 100 {
		go runOptionAndStockData(symbols[i:MinInt(i+100, len(symbols))])
	}
	fmt.Println("ended-loop @ " + time.Now().String())
	// return
	symbols0 := []string{"SPY", "DIA", "QQQ"}

	//2. activate option and stock data getter
	go runOptionAndStockData(symbols0)
	// for range time.Tick(time.Minute * 3) {
	// 	if isMarketOpen() {
	// 		go runOptionAndStockData(symbols0)
	// 	} else {
	// 		//break
	// 		continue
	// 	}
	// }

	// ticker := time.NewTicker(1 * time.Second)
	// fmt.Println("Started at ", time.Now())
	// defer ticker.Stop()
	// go func() {
	// 	for ; true; <-ticker.C {
	// 		fmt.Println("Tick at ", time.Now())
	// 	}
	// }()
	// time.Sleep(10 * tim.Second)
	// fmt.Println("Stopped at ", time.Now())

	//run all stocks

	/*3. terminate if option data getter
	and stock data getter are both terminated
	*/
	fmt.Println("runCore() ends")
}

// func isMarketOpen() bool {
// 	_, quote, _, err := new(YahooApiManager).GetOptionsAndStockDataBySymbol("SPY")
// 	if err != nil {
// 		return false
// 	}
// 	return quote.isMarketOpen()
// }

func runOptionAndStockData(symbols []string) {
	for _, symbol := range symbols {
		fmt.Printf("Run option and stock data for symbol %s\n", symbol)
		options, stock, _, err := new(YahooApiManager).GetOptionsAndStockDataBySymbol(symbol)
		if err != nil {
			continue
		}
		//store the option data into database
		for _, option := range options {
			var err2 error
			if option.GetSymbol() == "SPY" {
				err2 = new(TblOptionData).InsertOrUpdateOptionData(option, true)
			} else {
				err2 = new(TblOptionData).InsertOrUpdateOptionData(option, false)
			}
			if err2 != nil {
				panic(err2)
			}
		}
		//store the stock data into database
		err = new(TblStockData).InsertOrUpdateStockData(stock, false)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Complete option and stock data for symbol %s\n", symbol)
	}

	//fmt.Println("Run option report")
	//runOptionReport(symbols)
}

func runOptionReport(symbols []string) {
	//get an array of closet expiration date
	expDates := make([]int64, len(symbols))
	for index, symbol := range symbols {
		_, _, exp, err := new(YahooApiManager).GetOptionsAndStockDataBySymbol(symbol)
		if err != nil {
			fmt.Println(err)
			continue
		}
		expDates[index] = exp[0]
	}
	//run report every minute
	//for range time.Tick(time.Minute * 1) {
	go runOptionReportDetail(symbols, expDates)
	//}
}

func runOptionReportDetail(symbols []string, expDates []int64) {
	for _, symbol := range symbols {
		new(Orbit).runOptionReportForSymbol(symbol, GetTimeInYYYYMMDD(), false)
	}
}
