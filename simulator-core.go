package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
)

type Simulator struct {
}

type Transaction struct {
	BuyDate       int
	SellDate      int
	BuyPrice      float32
	SellPrice     float32
	ChangePrecent float32
	isOnHand      bool
	NetChange     float32
}

var wg sync.WaitGroup
var mut sync.Mutex

func (smlt *Simulator) RunSimulator() {
	var totalProfit float32
	symbolChan := smlt.ProducerSymbol()
	for i := 0; i < CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY; i++ {
		wg.Add(1)
		go smlt.ConsumerSymbol(symbolChan, 20181001, 20190526, &totalProfit)
	}
	wg.Wait()
	fmt.Println("Total Net Profit Percent is:", int(totalProfit*100))
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

func (smlt *Simulator) ProducerSymbol() <-chan string {
	outChan := make(chan string, CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY)
	tblSymbol := new(TblSymbol)
	tblLogError := new(TblLogError)
	//GET SYMBOL LIST FROM db_symbol
	// symbols, symbolSelectErr := tblSymbol.SelectSymbolByTrader(TRADER_NASDAQ)
	symbols, symbolSelectErr := tblSymbol.SelectSymbolByFilter()
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

func (smlt *Simulator) ConsumerSymbol(inChan <-chan string, startDate int, endDate int, totalProfit *float32) {
	for symbol := range inChan {
		smlt.SimulateInDateRange3(symbol, startDate, endDate, totalProfit)
	}
	wg.Done()
}

func (smlt *Simulator) SimulateInDateRange(symbol string, startDate int, endDate int, totalProfit *float32) {
	fmt.Println("Start simulating ", symbol)
	//Read data and build data struct
	file, fileOpenErr := os.Open("./daily_charts/" + symbol + ".csv")
	if fileOpenErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_IEX_API_MANAGER+" "+symbol, fileOpenErr.Error())
		return
	}
	defer file.Close()
	lines, csvReaderErr := csv.NewReader(file).ReadAll()
	if csvReaderErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_IEX_API_MANAGER+" "+symbol, csvReaderErr.Error())
		return
	}
	var dayCharts []IexDayChart
	for _, line := range lines {
		newDate := ConvertTimeInYYYYMMDD(line[0])
		if newDate < startDate {
			continue
		}
		if newDate > endDate {
			break
		}
		chart := new(IexDayChart).BuildFromCsv(line[0], line[1], line[2], line[3], line[4], line[5], line[6], line[7], line[8], line[9], line[10])
		dayCharts = append(dayCharts, *chart)
	}
	//Iterate through dates and add transactions to book
	var increaseCount int
	var transactions []Transaction
	var onHand Transaction
	for index, chart := range dayCharts {
		if onHand.isOnHand {
			if chart.High > onHand.BuyPrice*(1+PROFIT_PERCENT) {
				onHand.SellDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.SellPrice = onHand.BuyPrice * (1 + PROFIT_PERCENT)
				onHand.ChangePrecent = (onHand.SellPrice - onHand.BuyPrice) / onHand.BuyPrice
				onHand.NetChange = onHand.SellPrice - onHand.BuyPrice
				transactions = append(transactions, onHand)
				onHand = Transaction{}
			} else if chart.Low < onHand.BuyPrice*(1-LOSS_PERCENT) {
				onHand.SellDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.SellPrice = onHand.BuyPrice * (1 - LOSS_PERCENT)
				onHand.ChangePrecent = (onHand.SellPrice - onHand.BuyPrice) / onHand.BuyPrice
				onHand.NetChange = onHand.SellPrice - onHand.BuyPrice
				transactions = append(transactions, onHand)
				onHand = Transaction{}
			}
			increaseCount = 0
			continue
		}
		if increaseCount == 0 {
			increaseCount = 1
		} else {
			if chart.Vwap > dayCharts[index-1].Vwap {
				increaseCount++
			} else {
				increaseCount = 0
			}
		}
		if increaseCount >= CONTINOUS_PRICE_INCREASE_COUNT && chart.Vwap < dayCharts[index-CONTINOUS_PRICE_INCREASE_COUNT+1].Vwap*1.1 {
			if !onHand.isOnHand {
				onHand.isOnHand = true
				onHand.BuyDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.BuyPrice = chart.Close
				continue
			}
		}
	}
	if onHand.isOnHand {
		// lastChart := dayCharts[len(dayCharts)-1]
		// onHand.SellPrice = lastChart.Close
		// onHand.ChangePrecent = (onHand.SellPrice - onHand.BuyPrice) / onHand.BuyPrice
		// onHand.NetChange = onHand.SellPrice - onHand.BuyPrice
		transactions = append(transactions, onHand)
	}
	//Record transactions
	fmt.Println(transactions)
	var totalNetChange float32
	var totalChangePercent float32
	for _, trans := range transactions {
		totalNetChange += trans.NetChange
		totalChangePercent += trans.ChangePrecent
	}
	mut.Lock()
	*totalProfit += totalChangePercent
	if len(transactions) > 0 {
		fmt.Println(symbol, " Change Percent is: ", int(totalChangePercent*100))
	}
	mut.Unlock()
}

const (
	CONTINOUS_PRICE_INCREASE_COUNT = 3
	PROFIT_PERCENT                 = 0.03
	LOSS_PERCENT                   = 0.05
)

func (smlt *Simulator) SimulateInDateRange2(symbol string, startDate int, endDate int, totalProfit *float32) {
	fmt.Println("Start simulating ", symbol)
	//Read data and build data struct
	file, fileOpenErr := os.Open("./daily_charts/" + symbol + ".csv")
	if fileOpenErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_IEX_API_MANAGER+" "+symbol, fileOpenErr.Error())
		return
	}
	defer file.Close()
	lines, csvReaderErr := csv.NewReader(file).ReadAll()
	if csvReaderErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_IEX_API_MANAGER+" "+symbol, csvReaderErr.Error())
		return
	}
	var dayCharts []IexDayChart
	for _, line := range lines {
		newDate := ConvertTimeInYYYYMMDD(line[0])
		if newDate < startDate {
			continue
		}
		if newDate > endDate {
			break
		}
		chart := new(IexDayChart).BuildFromCsv(line[0], line[1], line[2], line[3], line[4], line[5], line[6], line[7], line[8], line[9], line[10])
		dayCharts = append(dayCharts, *chart)
	}
	//Iterate through dates and add transactions to book
	var transactions []Transaction
	var onHand Transaction
	var fiftyDayCloseTotal float32
	movingAverageRange := 50
	previousDayRange := 10
	previousDayRangePercent := 0.7
	for index, chart := range dayCharts {
		if index < movingAverageRange {
			fiftyDayCloseTotal += chart.Close //
			continue
		}
		if onHand.isOnHand {
			if chart.High > onHand.BuyPrice*(1+PROFIT_PERCENT) {
				onHand.SellDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.SellPrice = onHand.BuyPrice * (1 + PROFIT_PERCENT)
				onHand.ChangePrecent = (onHand.SellPrice - onHand.BuyPrice) / onHand.BuyPrice
				onHand.NetChange = onHand.SellPrice - onHand.BuyPrice
				transactions = append(transactions, onHand)
				onHand = Transaction{}
			} else if chart.Low < onHand.BuyPrice*(1-LOSS_PERCENT) {
				onHand.SellDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.SellPrice = onHand.BuyPrice * (1 - LOSS_PERCENT)
				onHand.ChangePrecent = (onHand.SellPrice - onHand.BuyPrice) / onHand.BuyPrice
				onHand.NetChange = onHand.SellPrice - onHand.BuyPrice
				transactions = append(transactions, onHand)
				onHand = Transaction{}
			}
			continue
		}
		fiftyDayCloseTotal += chart.Close                               //
		fiftyDayCloseTotal -= dayCharts[index-movingAverageRange].Close //
		var fiftyDayVwapAverage = fiftyDayCloseTotal / float32(movingAverageRange)
		// fmt.Printf("%v %.1f %.1f %.1f\r\n", chart.Date, fiftyDayVwapAverage, dayCharts[index-1].Close, chart.Close)
		if fiftyDayVwapAverage > dayCharts[index-1].Close && fiftyDayVwapAverage < chart.Close {
			if !onHand.isOnHand {
				var count int
				for i := 0; i < previousDayRange; i++ {
					if dayCharts[index-i].Close < fiftyDayVwapAverage {
						count++
					}
				}
				if count > int(float32(previousDayRange)*float32(previousDayRangePercent)) {
					onHand.isOnHand = true
					onHand.BuyDate = ConvertTimeInYYYYMMDD(chart.Date)
					onHand.BuyPrice = chart.Close
				}
				continue
			}
		}
	}
	if onHand.isOnHand {
		// lastChart := dayCharts[len(dayCharts)-1]
		// onHand.SellPrice = lastChart.Close
		// onHand.ChangePrecent = (onHand.SellPrice - onHand.BuyPrice) / onHand.BuyPrice
		// onHand.NetChange = onHand.SellPrice - onHand.BuyPrice
		transactions = append(transactions, onHand)
	}
	//Record transactions
	fmt.Println(transactions)
	var totalNetChange float32
	var totalChangePercent float32
	for _, trans := range transactions {
		totalNetChange += trans.NetChange
		totalChangePercent += trans.ChangePrecent
	}
	mut.Lock()
	*totalProfit += totalChangePercent
	if len(transactions) > 0 {
		fmt.Println(symbol, " Change Percent is: ", int(totalChangePercent*100))
	}
	mut.Unlock()
}

func (smlt *Simulator) SimulateInDateRange3(symbol string, startDate int, endDate int, totalProfit *float32) {
	//	fmt.Println("Start simulating ", symbol)
	//Read data and build data struct
	file, fileOpenErr := os.Open("./daily_charts/" + symbol + ".csv")
	if fileOpenErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_IEX_API_MANAGER+" "+symbol, fileOpenErr.Error())
		return
	}
	defer file.Close()
	lines, csvReaderErr := csv.NewReader(file).ReadAll()
	if csvReaderErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_IEX_API_MANAGER+" "+symbol, csvReaderErr.Error())
		return
	}
	var dayCharts []IexDayChart
	for _, line := range lines {
		newDate := ConvertTimeInYYYYMMDD(line[0])
		if newDate < startDate {
			continue
		}
		if newDate > endDate {
			break
		}
		chart := new(IexDayChart).BuildFromCsv(line[0], line[1], line[2], line[3], line[4], line[5], line[6], line[7], line[8], line[9], line[10])
		dayCharts = append(dayCharts, *chart)
	}
	//Iterate through dates and add transactions to book
	var transactions []Transaction
	var leftOver []Transaction
	var onHand Transaction
	var ma60Total float32
	var ma120Total float32
	ma60Range := 60
	ma120Range := 120
	for index, chart := range dayCharts {

		//Transfer data to db
		// fmt.Println("Transferring data for ", symbol)
		stock := YahooQuote{
			Symbol:               symbol,
			RegularMarketOpen:    chart.Open,
			RegularMarketDayHigh: chart.High,
			RegularMarketDayLow:  chart.Low,
			RegularMarketPrice:   chart.Close,
			RegularMarketVolume:  chart.Volume,
		}
		new(TblStockHist).InsertOrUpdateStockData(stock, ConvertTimeInYYYYMMDD(chart.Date))

		if onHand.isOnHand {
			if chart.High > MaxFloat32(ma60Total/float32(ma60Range), ma120Total/float32(ma120Range)) {
				onHand.SellDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.SellPrice = MaxFloat32(ma60Total/float32(ma60Range), ma120Total/float32(ma120Range))
				onHand.ChangePrecent = (onHand.SellPrice - onHand.BuyPrice) / onHand.BuyPrice
				onHand.NetChange = onHand.SellPrice - onHand.BuyPrice
				transactions = append(transactions, onHand)
				onHand = Transaction{}
			} else if chart.Low < MinFloat32(ma60Total/float32(ma60Range), ma120Total/float32(ma120Range)) {
				onHand.SellDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.SellPrice = MinFloat32(ma60Total/float32(ma60Range), ma120Total/float32(ma120Range)) // * (1 - LOSS_PERCENT)
				onHand.ChangePrecent = (onHand.SellPrice - onHand.BuyPrice) / onHand.BuyPrice
				onHand.NetChange = onHand.SellPrice - onHand.BuyPrice
				transactions = append(transactions, onHand)
				onHand = Transaction{}
			}
			continue
		}
		ma60Total += (chart.High + chart.Low + 2*chart.Close) / 4
		ma120Total += (chart.High + chart.Low + 2*chart.Close) / 4
		if index >= ma60Range {
			dayChart := dayCharts[index-ma60Range]
			ma60Total -= (dayChart.High + dayChart.Low + 2*dayChart.Close) / 4
		}
		if index >= ma120Range {
			dayChart := dayCharts[index-ma120Range]
			ma120Total -= (dayChart.High + dayChart.Low + 2*dayChart.Close) / 4
		}
		if index < ma60Range || index < ma120Range {
			continue
		}
		var ma60Average = ma60Total / float32(ma60Range)
		var ma120Average = ma120Total / float32(ma120Range)

		//Transfer data to Db
		report := RowStockReport{
			Symbol: symbol,
			Date:   ConvertTimeInYYYYMMDD(chart.Date),
			MA60:   ma60Average,
			MA120:  ma120Average,
		}
		new(TblStockReport).InsertOrUpdateStockData(report)

		// fmt.Printf("%v %.1f %.1f\r\n", chart.Date, ma60Average, ma120Average)
		if !onHand.isOnHand {
			if ma120Average/ma60Average > 1 && chart.Close > ma60Average && dayCharts[index-1].Close < ma60Average {
				onHand.isOnHand = true
				onHand.BuyDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.BuyPrice = chart.Close
				// fmt.Printf("%v %.1f %.1f %.1f %.1f\r\n", chart.Date, ma120Average, ma120Average, dayCharts[index-1].Close, chart.Close)
			} else if ma60Average/ma120Average > 1 && chart.Close > ma120Average && dayCharts[index-1].Close < ma120Average {
				onHand.isOnHand = true
				onHand.BuyDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.BuyPrice = chart.Close
				// fmt.Printf("%v %.1f %.1f %.1f %.1f\r\n", chart.Date, ma120Average, ma120Average, dayCharts[index-1].Close, chart.Close)
			}
		}
	}
	if onHand.isOnHand {
		// lastChart := dayCharts[len(dayCharts)-1]
		// onHand.SellPrice = lastChart.Close
		// onHand.ChangePrecent = (onHand.SellPrice - onHand.BuyPrice) / onHand.BuyPrice
		// onHand.NetChange = onHand.SellPrice - onHand.BuyPrice
		leftOver = append(leftOver, onHand)
	}
	//Record transactions
	var totalNetChange float32
	var totalChangePercent float32
	for _, trans := range transactions {
		totalNetChange += trans.NetChange
		totalChangePercent += trans.ChangePrecent
	}
	mut.Lock()
	*totalProfit += totalChangePercent
	if len(leftOver) > 0 {
		fmt.Println(symbol, "LeftOver", leftOver)
		//fmt.Println(symbol, " Change Percent is: ", int(totalChangePercent*100))
	}
	mut.Unlock()
}
