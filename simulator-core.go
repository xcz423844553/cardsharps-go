package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math"
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
	fmt.Println("Total Net Profit is:", int(totalProfit))
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
	// symbols, symbolSelectErr := tblSymbol.SelectSymbolByTrader(TRADER_RUSSELL)
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

func (smlt *Simulator) ConsumerSymbol(inChan <-chan string, startDate int, endDate int, totalProfit *float32) {
	for symbol := range inChan {
		smlt.SimulateInDateRange4(symbol, startDate, endDate, totalProfit)
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

func (smlt *Simulator) SimulateInDateRange4(symbol string, startDate int, endDate int, totalProfit *float32) {
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
		//Date	Open	High	Low	Close	Volume	UnadjustedVolume	Change	ChangePercent	Vwap	ChangeOverTime
		chart := new(IexDayChart).BuildFromCsv(line[0], line[1], line[2], line[3], line[4], line[5], line[6], line[7], line[8], line[9], line[10])
		dayCharts = append(dayCharts, *chart)
	}
	//Iterate through dates and add transactions to book
	var transactions []Transaction
	var leftOver []Transaction
	var onHand Transaction
	maRange := 20
	sigmaWidth := 2
	atrRange := 10
	var bollUpperArray []float32
	var bollLowerArray []float32
	var kcUpperArray []float32
	var kcLowerArray []float32
	for index, chart := range dayCharts {
		if index < 2*maRange {
			continue
		}

		//Calculate BOLL and KC
		ma, sd, ema, atr, statErr := smlt.getBasicStat(dayCharts, index, maRange, atrRange)
		if statErr != nil {
			PrintMsgInConsole(MSGERROR, LOGTYPE_IEX_API_MANAGER+" "+symbol, statErr.Error())
			return
		}
		//fmt.Printf("Symbol: %s Date: %v MA: %.2f SD: %.2f EMA: %.2f ATR: %.2f\r\n", symbol, chart.Date, ma, sd, ema, atr)

		bollMid := ma
		bollUpper := bollMid + float32(sigmaWidth)*sd
		bollLower := bollMid - float32(sigmaWidth)*sd
		kcMid := ema
		kcUpper := kcMid + float32(sigmaWidth)*atr
		kcLower := kcMid - float32(sigmaWidth)*atr
		bollUpperArray = append(bollUpperArray, bollUpper)
		bollLowerArray = append(bollLowerArray, bollLower)
		kcUpperArray = append(kcUpperArray, kcUpper)
		kcLowerArray = append(kcLowerArray, kcLower)
		//fmt.Printf("Symbol: %s Date: %v BMID: %.2f BUPPER: %.2f BLOWER: %.2f KCMID: %.2f KCUPPER: %.2f KCLOWER: %.2f\r\n", symbol, chart.Date, bollMid, bollUpper, bollLower, kcMid, kcUpper, kcLower)

		if onHand.isOnHand {
			if chart.High > onHand.BuyPrice*(1+0.1) {
				onHand.SellDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.SellPrice = onHand.BuyPrice * (1 + 0.1)
				onHand.ChangePrecent = (onHand.SellPrice - onHand.BuyPrice) / onHand.BuyPrice
				onHand.NetChange = onHand.SellPrice - onHand.BuyPrice
				fmt.Printf("Sell %s %v %.2f %.2f %.2f perc\r\n", symbol, chart.Date, onHand.BuyPrice, onHand.SellPrice, onHand.ChangePrecent*100)
				transactions = append(transactions, onHand)
				onHand = Transaction{}
			} else if chart.Low < bollMid {
				onHand.SellDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.SellPrice = bollMid
				onHand.ChangePrecent = (onHand.SellPrice - onHand.BuyPrice) / onHand.BuyPrice
				onHand.NetChange = onHand.SellPrice - onHand.BuyPrice
				fmt.Printf("Sell %s %v %.2f %.2f %.2f perc\r\n", symbol, chart.Date, onHand.BuyPrice, onHand.SellPrice, onHand.ChangePrecent*100)
				transactions = append(transactions, onHand)
				onHand = Transaction{}
			}
			continue
		}

		if len(bollUpperArray)-10 < 0 {
			continue
		}
		if !onHand.isOnHand {
			if bollUpper > kcUpper && bollUpperArray[len(bollUpperArray)-10] < kcUpperArray[len(kcUpperArray)-10] && chart.High > bollUpper {
				onHand.isOnHand = true
				onHand.BuyDate = ConvertTimeInYYYYMMDD(chart.Date)
				onHand.BuyPrice = bollUpper
				fmt.Printf("Buy %s %v %.2f\r\n", symbol, chart.Date, onHand.BuyPrice)
			}
			//short position
			// else if bollLower < kcLower && bollLowerArray[len(bollLowerArray)-6] > kcLowerArray[len(kcLowerArray)-6] && chart.Low < bollLower {
			// 	onHand.isOnHand = true
			// 	onHand.BuyDate = ConvertTimeInYYYYMMDD(chart.Date)
			// 	onHand.BuyPrice = chart.Close
			// 	// fmt.Printf("%v %.1f %.1f %.1f %.1f\r\n", chart.Date, ma120Average, ma120Average, dayCharts[index-1].Close, chart.Close)
			// }
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
		totalNetChange += trans.ChangePrecent * 1000
		totalChangePercent += trans.ChangePrecent
	}
	mut.Lock()
	*totalProfit += totalNetChange
	if len(leftOver) > 0 {
		fmt.Println(symbol, "LeftOver", leftOver)
		fmt.Println(symbol, " Change Percent is: \r\n", int(totalChangePercent*100))
	}
	mut.Unlock()
}

func (smlt *Simulator) getBasicStat(charts []IexDayChart, targetIndex int, indexRange int, atrRange int) (float32, float32, float32, float32, error) {
	var average float64
	var sd float64
	var atr float64
	if targetIndex < 2*indexRange {
		return float32(0), float32(average), float32(sd), float32(atr), errors.New(fmt.Sprintf("Not enough data to get basic stat. %v / %v \r\n", targetIndex, indexRange))
	}
	for i := indexRange; i > 0; i-- {
		var curIndex = targetIndex - i
		average += float64((charts[curIndex].High + charts[curIndex].Low + 2*charts[curIndex].Close) / 4)
	}
	average = average / float64(indexRange)
	for i := indexRange; i > 0; i-- {
		var curIndex = targetIndex - i
		sdElm := float64((charts[curIndex].High+charts[curIndex].Low+2*charts[curIndex].Close)/4) - average
		sd += math.Pow(sdElm, 2)
	}
	sd = math.Sqrt(sd / float64(indexRange))

	var emaPrev float64
	var ema float64
	var multiplier float64
	multiplier = 2 / float64(indexRange+1)
	for i := indexRange; i > 0; i-- {
		var curIndex = targetIndex - i
		if i == indexRange {
			for j := indexRange; j > 0; j-- {
				emaPrev += float64((charts[curIndex-j].High + charts[curIndex-j].Low + 2*charts[curIndex-j].Close) / 4)
			}
			emaPrev = emaPrev / float64(indexRange)
		} else {
			emaPrev = ema
		}
		ema = float64((charts[curIndex].High+charts[curIndex].Low+2*charts[curIndex].Close)/4)*multiplier + emaPrev*(1-multiplier)
	}

	for i := atrRange; i > 0; i-- {
		var curIndex = targetIndex - i
		trElm1 := float64(charts[curIndex].High - charts[curIndex].Low)
		trElm2 := math.Abs(float64(charts[curIndex].High - charts[curIndex-1].Close))
		trElm3 := math.Abs(float64(charts[curIndex].Low - charts[curIndex-1].Close))
		tr := math.Max(math.Max(trElm1, trElm2), trElm3)
		atr += tr
	}
	atr = atr / float64(atrRange)
	return float32(average), float32(sd), float32(ema), float32(atr), nil
}
