package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Shuffler is a struct for recovering the historical stock data
type Shuffler2 struct {
}

//InitHistoricalStockDataFromYahoo initiate the historical stock data based on start date and end date of user inputs
func (sf *Shuffler2) InitHistoricalStockDataFromYahoo() {
	fmt.Println("The program is initiating the historical stock data. Type 'y' to continue. Type any other chars to skip.")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	if strings.Compare("y", text) == 0 {
		fmt.Print("-> Start Date in YYYYMMDD: ")
		startDate, _ := reader.ReadString('\n')
		startDate = strings.Replace(startDate, "\n", "", -1)
		fmt.Print("-> End Date in YYYYMMDD: ")
		endDate, _ := reader.ReadString('\n')
		endDate = strings.Replace(endDate, "\n", "", -1)
		fmt.Println("Initiating the historical stock data.")
		sDate, _ := strconv.Atoi(startDate)
		eDate, _ := strconv.Atoi(endDate)
		sf.RecoverHistoricalStockDataFromYahoo(sDate, eDate)
		fmt.Println("Finished initiating the historical stock data from .", startDate, " to ", endDate)
	} else {
		fmt.Println("Historical stock data initiation skipped.")
	}
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
