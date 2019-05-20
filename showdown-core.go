package main

type Showdown struct {
}

func (sd Showdown) runShowdown() {
	tblSymbol := new(TblSymbol)
	tblLogSystem := new(TblLogSystem)
	tblLogError := new(TblLogError)
	dealer := new(Dealer)
	tblLogSystem.InsertLogSystem(LOGTYPE_SHOWDOWN, "Showdown Core Started")
	//GET SYMBOL LIST FROM db_symbol
	symbols, symbolSelectErr := tblSymbol.SelectSymbolByFilter()
	if symbolSelectErr != nil {
		tblLogError.InsertLogError(LOGTYPE_DB_SYMBOL, symbolSelectErr.Error())
	}
	//TODO UPDATED WITH ETF SYMBOLS
	symbols = append(symbols, "SPY")
	symbols = append(symbols, "DIA")
	symbols = append(symbols, "QQQ")
	for i := 0; i < len(symbols); i = i + MAX_NUM_SYMBOL_ON_EACH_GOROUTINE {
		go dealer.GetOptionAndStockDataFromYahoo(symbols[i:MinInt(i+MAX_NUM_SYMBOL_ON_EACH_GOROUTINE, len(symbols))])
	}
	tblLogSystem.InsertLogSystem(LOGTYPE_SHOWDOWN, "Showdown Core Finished")
}

func isMarketOpen() bool {
	_, quote, _, err := new(YahooApiManager).GetOptionsAndStockDataBySymbol("SPY")
	if err != nil {
		return false
	}
	return quote.isMarketOpen()
}

func (sd Showdown) runShowdown2() {
	tblSymbol := new(TblSymbol)
	tblLogSystem := new(TblLogSystem)
	tblLogError := new(TblLogError)
	tblLogSystem.InsertLogSystem(LOGTYPE_SHOWDOWN, "Showdown Core Started")
	//GET SYMBOL LIST FROM db_symbol
	symbols, symbolSelectErr := tblSymbol.SelectSymbolByFilter()
	if symbolSelectErr != nil {
		tblLogError.InsertLogError(LOGTYPE_DB_SYMBOL, symbolSelectErr.Error())
	}
	//TODO UPDATED WITH ETF SYMBOLS
	symbols = append(symbols, "SPY")
	symbols = append(symbols, "DIA")
	symbols = append(symbols, "QQQ")
	queue := new(ShowdownQueue)
	queue.SetSymbols(symbols)
	for i := 0; i < len(symbols); i = i + 30 {
		queue.waitGroup.Add(30)
		for j := 0; j < 30; j = j + 1 {
			go func(q *ShowdownQueue) {
				dealer := new(Dealer)
				dealer.GetOptionAndStockDataFromYahoo([]string{q.GetNextSymbol()})
				queue.waitGroup.Done()
			}(queue)
		}
		queue.waitGroup.Wait()
	}
	tblLogSystem.InsertLogSystem(LOGTYPE_SHOWDOWN, "Showdown Core Finished")
}
