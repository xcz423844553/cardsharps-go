package main

import "fmt"

//Shuffler is a struct for recovering the historical stock data
type Shuffler2 struct {
}

//RecoverHistoricalStockDataFromYahoo downloads historical stock data from finance.yahoo.com and store the stock data in the database
//Param: start date - start date of the historical stock data in format of YYYYMMDD, type int
//Param: end date - end date of the historical stock data in format of YYYYMMDD, type int
//Return: void
func (sf *Shuffler2) RecoverHistoricalStockDataFromYahoo(startDate int, endDate int) {
	bd := new(Board2)
	execution := func(symbol string) {
		sf.RecoverHistoricalStockDataBySymbolFromYahoo(symbol, startDate, endDate)
	}
	callback := func() {
		fmt.Println("Completed the recovery of historical stock data.")
	}
	bd.StartGame(bd.GetSymbolTagAll(), execution, callback)
}

//RecoverHistoricalStockDataBySymbolFromYahoo downloads historical stock data of a given symbol from finance.yahoo.com and store the stock data in the database
//Param: symbol - stock symbol/quote
//Param: start date - start date of the historical stock data in format of YYYYMMDD, type int
//Param: end date - end date of the historical stock data in format of YYYYMMDD, type int
//Return: void
func (sf *Shuffler2) RecoverHistoricalStockDataBySymbolFromYahoo(symbol string, startDate int, endDate int) {
	fmt.Println("Started the recovery of historical stock data of " + symbol)
	yahooApi := new(YahooApi)
	stockHistArray, yahooErr := yahooApi.GetStockHist(symbol, startDate, endDate)
	if yahooErr != nil {
		fmt.Println("Yahoo API Error: " + yahooErr.Error())
		return
	}
	for _, stockHist := range stockHistArray {
		new(DaoStockHist).InsertOrUpdateStockHist(&stockHist, stockHist.GetDate())

	}
	fmt.Println("Completed the recovery of historical stock data of " + symbol)
}
