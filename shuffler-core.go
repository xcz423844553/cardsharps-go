package main

import (
	"encoding/csv"
	"net/http"
	"strconv"
)

//Shuffler is a library for recovering the historical stock data
type Shuffler struct {
}

//RecoverHistoricalStockDataFromMacroTrends downloads historical stock data from MacroTrends.com and store the stock data in the database
//Param: start date - start date of the historical stock data in format of YYYYMMDD, type int
//Return: void
func (sf *Shuffler) RecoverHistoricalStockDataFromMacroTrends(startDate int) {
	bd := new(Board)
	execution := func(symbol string) {
		sf.RecoverHistoricalStockDataBySymbolFromMacroTrends(symbol, startDate)
	}
	callback := func() {
		PrintMsgInConsole(MSGSYSTEM, LOGTYPE_SHUFFLER, "Completed the recovery of historical stock data.")
	}
	bd.SymbolGame(SYMBOLTAG_ALL, execution, callback)
}

//RecoverHistoricalStockDataBySymbolFromMacroTrends downloads historical stock data of a given symbol from MacroTrends.com and store the stock data in database "stock_hist2"
//Param: symbol - stock symbol/quote
//Param: start date - start date of the historical stock data in format of YYYYMMDD, type int
//Return: void
func (sf *Shuffler) RecoverHistoricalStockDataBySymbolFromMacroTrends(symbol string, startDate int) {
	//Download csv file from MacroTrends
	url := URL_MACROTRENDS + symbol
	resp, downloadErr := http.Get(url)
	if downloadErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SHUFFLER, "Download Error: "+downloadErr.Error())
		return
	}
	defer resp.Body.Close()

	//Read csv content
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true
	data, readerErr := reader.ReadAll()
	if readerErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SHUFFLER, "Reader Error: "+readerErr.Error())
		return
	}

	//Iterate lines of the csv file, the real data start from line 16
	for idx, row := range data {
		if idx >= 15 {
			colDate := ConvertTimeInYYYYMMDD(row[0])
			colOpen, _ := strconv.ParseFloat(row[1], 32)
			colHigh, _ := strconv.ParseFloat(row[2], 32)
			colLow, _ := strconv.ParseFloat(row[3], 32)
			colClose, _ := strconv.ParseFloat(row[4], 32)
			colVolume, _ := strconv.Atoi(row[5])
			if colDate >= startDate {
				stock := YahooQuote{
					Symbol:               symbol,
					RegularMarketOpen:    float32(colOpen),
					RegularMarketDayHigh: float32(colHigh),
					RegularMarketDayLow:  float32(colLow),
					RegularMarketPrice:   float32(colClose),
					RegularMarketVolume:  colVolume,
				}
				new(TblStockHist).InsertOrUpdateStockData(stock, colDate)
			}
		}
	}
	PrintMsgInConsole(MSGSYSTEM, LOGTYPE_SHUFFLER, "Completed the recovery of historical stock data of "+symbol)
}

//RecoverHistoricalStockDataFromYahoo downloads historical stock data from finance.yahoo.com and store the stock data in the database
//Param: start date - start date of the historical stock data in format of YYYYMMDD, type int
//Param: end date - end date of the historical stock data in format of YYYYMMDD, type int
//Return: void
func (sf *Shuffler) RecoverHistoricalStockDataFromYahoo(startDate int, endDate int) {
	bd := new(Board)
	execution := func(symbol string) {
		sf.RecoverHistoricalStockDataBySymbolFromYahoo(symbol, startDate, endDate)
	}
	callback := func() {
		PrintMsgInConsole(MSGSYSTEM, LOGTYPE_SHUFFLER, "Completed the recovery of historical stock data.")
	}
	bd.SymbolGame(SYMBOLTAG_ALL, execution, callback)
}

//RecoverHistoricalStockDataBySymbolFromYahoo downloads historical stock data of a given symbol from finance.yahoo.com and store the stock data in the database
//Param: symbol - stock symbol/quote
//Param: start date - start date of the historical stock data in format of YYYYMMDD, type int
//Param: end date - end date of the historical stock data in format of YYYYMMDD, type int
//Return: void
func (sf *Shuffler) RecoverHistoricalStockDataBySymbolFromYahoo(symbol string, startDate int, endDate int) {

	PrintMsgInConsole(MSGSYSTEM, LOGTYPE_SHUFFLER, "Started the recovery of historical stock data of "+symbol)

	yahooAPIManager := new(YahooAPIManager)

	stockHistArray, yahooErr := yahooAPIManager.GetStockHistoryFromYahoo(symbol, startDate, endDate)

	if yahooErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SHUFFLER, "Yahoo API Error: "+yahooErr.Error())
		return
	}

	for _, stockHist := range stockHistArray {
		stock := YahooQuote{
			Symbol:               symbol,
			RegularMarketOpen:    stockHist.Open,
			RegularMarketDayHigh: stockHist.High,
			RegularMarketDayLow:  stockHist.Low,
			RegularMarketPrice:   stockHist.Close,
			RegularMarketVolume:  stockHist.Volume,
		}
		new(TblStockHist).InsertOrUpdateStockData(stock, stockHist.Date)

	}
	PrintMsgInConsole(MSGSYSTEM, LOGTYPE_SHUFFLER, "Completed the recovery of historical stock data of "+symbol)
}
