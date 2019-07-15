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
	symbols, symbolSelectErr := tblSymbol.SelectAllSymbol()
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
	_, quote, _, err := new(YahooAPIManager).GetOptionsAndStockDataBySymbol("QQQ")
	if err != nil {
		return false
	}
	return quote.isMarketOpen()
}

func isMarketPreOpenOrOpen() bool {
	_, quote, _, err := new(YahooAPIManager).GetOptionsAndStockDataBySymbol("QQQ")
	if err != nil {
		return false
	}
	return quote.isMarketPreOpen() || quote.isMarketOpen()
}

func (sd *Showdown) runShowdown2() {
	tblLogSystem := new(TblLogSystem)
	tblLogSystem.InsertLogSystem(LOGTYPE_SHOWDOWN, "Showdown Core Started")
	symbolChan := sd.ProducerSymbol()
	for i := 0; i < CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY; i++ {
		go sd.ConsumerSymbol(symbolChan)
	}
	return

	//TODO UPDATED WITH ETF SYMBOLS
	// symbols = append(symbols, "SPY")
	// symbols = append(symbols, "DIA")
	// symbols = append(symbols, "QQQ")
	// queue := new(ShowdownQueue)
	// queue.SetSymbols(symbols)
	// for i := 0; i < len(symbols); i = i + 30 {
	// 	queue.waitGroup.Add(30)
	// 	for j := 0; j < 30; j = j + 1 {
	// 		go func(q *ShowdownQueue) {
	// 			dealer := new(Dealer)
	// 			dealer.GetOptionAndStockDataFromYahoo([]string{q.GetNextSymbol()})
	// 			queue.waitGroup.Done()
	// 		}(queue)
	// 	}
	// 	queue.waitGroup.Wait()
	// }
	// tblLogSystem.InsertLogSystem(LOGTYPE_SHOWDOWN, "Showdown Core Finished")
}

func (sd *Showdown) ProducerSymbol() <-chan string {
	outChan := make(chan string, CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY)
	tblSymbol := new(TblSymbol)
	tblLogError := new(TblLogError)
	//GET SYMBOL LIST FROM db_symbol
	symbols, symbolSelectErr := tblSymbol.SelectAllSymbol()
	if symbolSelectErr != nil {
		tblLogError.InsertLogError(LOGTYPE_DB_SYMBOL, symbolSelectErr.Error())
	}
	// symbols = []string{"AAPL"}
	go func() {
		for _, symbol := range symbols {
			outChan <- symbol
		}
		defer close(outChan)
	}()
	return outChan
}

func (sd *Showdown) ConsumerSymbol(inChan <-chan string) {
	dealer := new(Dealer)
	for symbol := range inChan {
		dealer.GetOptionAndStockData(symbol)
	}
}
