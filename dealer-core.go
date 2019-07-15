package main

import (
	"fmt"
	"time"
)

type Dealer struct {
}

func (dealer *Dealer) GetOptionAndStockData(symbol string) {
	tblLogError := new(TblLogError)
	tblOptionDate := new(TblOptionData)
	tblStockData := new(TblStockData)
	tblStockHist := new(TblStockHist)
	tblStockReport := new(TblStockReport)
	yahooApiManager := new(YahooAPIManager)
	//orbit := new(Orbit)
	isEtf := false
	PrintMsgInConsole(MSGSYSTEM, LOGTYPE_DEALER, "Run option and stock data for symbol "+symbol)
	options, stock, _, yahooApiErr := yahooApiManager.GetOptionsAndStockDataBySymbol(symbol)
	if yahooApiErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_DEALER+" "+symbol, yahooApiErr.Error())
		tblLogError.InsertLogError(LOGTYPE_YAHOO_API_MANAGER, yahooApiErr.Error())
		return
	}
	//store the option data into database
	var insertErr error
	for _, option := range options {
		insertErr = tblOptionDate.InsertOrUpdateOptionData(option, isEtf)
		if insertErr != nil {
			PrintMsgInConsole(MSGERROR, LOGTYPE_DEALER+" "+symbol, insertErr.Error())
			tblLogError.InsertLogError(LOGTYPE_DB_OPTION_DATA, insertErr.Error())
			continue
		}
	}
	//store the stock data into database
	insertErr = tblStockData.InsertOrUpdateStockData(stock, isEtf)
	if insertErr != nil {
		tblLogError.InsertLogError(LOGTYPE_DB_STOCK_DATA, insertErr.Error())
	}
	insertErr = tblStockHist.InsertOrUpdateStockData(stock, GetTimeInYYYYMMDD())
	if insertErr != nil {
		tblLogError.InsertLogError(LOGTYPE_DB_STOCK_DATA, insertErr.Error())
	}

	ma60Range := 60
	ma120Range := 120
	histList, histErr := tblStockHist.SelectLastStockHistByCountAndBeforeDate(symbol, ma120Range, GetTimeInYYYYMMDD()+1)
	if histErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_DEALER+" "+symbol, insertErr.Error())
		tblLogError.InsertLogError(LOGTYPE_DB_STOCK_REPORT, insertErr.Error())
		return
	}
	if len(histList) < ma120Range {
		PrintMsgInConsole(MSGERROR, LOGTYPE_DEALER+" "+symbol, "Not enough hist data for MA120")
		return
	}
	var ma60Total float32
	var ma120Total float32
	for i := 0; i < ma60Range; i++ {
		var histAverage float32
		if histList[i].MarketOpen == 0 {
			histAverage = histList[i].MarketClose
		} else {
			histAverage = (histList[i].MarketHigh + histList[i].MarketLow + 2*histList[i].MarketClose) / 4
		}
		ma60Total += histAverage
	}
	for i := 0; i < ma120Range; i++ {
		var histAverage float32
		if histList[i].MarketOpen == 0 {
			histAverage = histList[i].MarketClose
		} else {
			histAverage = (histList[i].MarketHigh + histList[i].MarketLow + 2*histList[i].MarketClose) / 4
		}
		ma120Total += histAverage
	}

	report := RowStockReport{
		Symbol: stock.Symbol,
		Date:   GetTimeInYYYYMMDD(),
		MA60:   ma60Total / float32(ma60Range),
		MA120:  ma120Total / float32(ma120Range),
	}
	insertErr = tblStockReport.InsertOrUpdateStockData(report)
	if insertErr != nil {
		tblLogError.InsertLogError(LOGTYPE_DB_STOCK_REPORT, insertErr.Error())
	}

	//go orbit.runOptionReportForSymbol(symbol, GetTimeInYYYYMMDD(), false)
	PrintMsgInConsole(MSGSYSTEM, LOGTYPE_DEALER, "Complete option and stock data for symbol "+symbol)
	fmt.Println("ended-service @ " + time.Now().String())
}

func (dealer Dealer) GetOptionAndStockDataFromYahoo(symbols []string) {
	tblLogError := new(TblLogError)
	tblOptionDate := new(TblOptionData)
	tblStockDate := new(TblStockData)
	yahooApiManager := new(YahooAPIManager)
	//orbit := new(Orbit)
	isEtf := false
	for _, symbol := range symbols {
		PrintMsgInConsole(MSGSYSTEM, LOGTYPE_DEALER, "Run option and stock data for symbol "+symbol)
		options, stock, _, yahooApiErr := yahooApiManager.GetOptionsAndStockDataBySymbol(symbol)
		if yahooApiErr != nil {
			tblLogError.InsertLogError(LOGTYPE_YAHOO_API_MANAGER, yahooApiErr.Error())
			continue
		}
		//store the option data into database
		var insertErr error
		for _, option := range options {
			insertErr = tblOptionDate.InsertOrUpdateOptionData(option, isEtf)
			if insertErr != nil {
				tblLogError.InsertLogError(LOGTYPE_DB_OPTION_DATA, insertErr.Error())
				continue
			}
		}
		//store the stock data into database
		insertErr = tblStockDate.InsertOrUpdateStockData(stock, isEtf)
		if insertErr != nil {
			tblLogError.InsertLogError(LOGTYPE_DB_STOCK_DATA, insertErr.Error())
			continue
		}
		//go orbit.runOptionReportForSymbol(symbol, GetTimeInYYYYMMDD(), false)
		PrintMsgInConsole(MSGSYSTEM, LOGTYPE_DEALER, "Complete option and stock data for symbol "+symbol)
	}
	fmt.Println("ended-service @ " + time.Now().String())
}

func (dealer Dealer) GetOptionAndEtfDataFromYahoo(symbols []string) {
	tblLogError := new(TblLogError)
	tblOptionDate := new(TblOptionData)
	tblStockDate := new(TblStockData)
	yahooApiManager := new(YahooAPIManager)
	orbit := new(Orbit)
	isEtf := true
	for _, symbol := range symbols {
		PrintMsgInConsole(MSGSYSTEM, LOGTYPE_DEALER, "Run option and etf data for symbol "+symbol)
		options, stock, _, yahooApiErr := yahooApiManager.GetOptionsAndStockDataBySymbol(symbol)
		if yahooApiErr != nil {
			tblLogError.InsertLogError(LOGTYPE_YAHOO_API_MANAGER, yahooApiErr.Error())
			continue
		}
		//store the option data into database
		var insertErr error
		for _, option := range options {
			insertErr = tblOptionDate.InsertOrUpdateOptionData(option, isEtf)
			if insertErr != nil {
				tblLogError.InsertLogError(LOGTYPE_DB_OPTION_DATA_ETF, insertErr.Error())
				continue
			}
		}
		//store the stock data into database
		insertErr = tblStockDate.InsertOrUpdateStockData(stock, isEtf)
		if insertErr != nil {
			tblLogError.InsertLogError(LOGTYPE_DB_STOCK_DATA_ETF, insertErr.Error())
			continue
		}
		go orbit.runOptionReportForSymbol(symbol, GetTimeInYYYYMMDD(), true)
		PrintMsgInConsole(MSGSYSTEM, LOGTYPE_DEALER, "Complete option and etf data for symbol "+symbol)
	}
}
