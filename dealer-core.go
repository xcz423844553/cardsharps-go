package main

type Dealer struct {
}

func (dealer Dealer) GetOptionAndStockDataFromYahoo(symbols []string) {
	tblLogError := new(TblLogError)
	tblOptionDate := new(TblOptionData)
	tblStockDate := new(TblStockData)
	yahooApiManager := new(YahooApiManager)
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
}

func (dealer Dealer) GetOptionAndEtfDataFromYahoo(symbols []string) {
	tblLogError := new(TblLogError)
	tblOptionDate := new(TblOptionData)
	tblStockDate := new(TblStockData)
	yahooApiManager := new(YahooApiManager)
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
