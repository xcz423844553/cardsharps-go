package main

import (
	"fmt"
	"math"
	"sort"
)

//Monitor is a library to monitor the market status
type Monitor struct {
}

//StockMonitor is to store the buying opportunity found by the monitor
type StockMonitor struct {
	stockChances map[string]*StockChance
}

//StockChance represents a buying opportunity found by the monitor
type StockChance struct {
	Symbol    string
	Price     float32
	BollMid   float32
	BollUpper float32
	BollLower float32
	KcMid     float32
	KcUpper   float32
	KcLower   float32
	Trend     string
	Status    string
}

const (
	PARAM_CROSS_RANGE          = 5   //Range to determine Boll-Kc trend
	PARAM_LOOKBACK_RANGE       = 20  //Range to look back to calculate EMA
	PARAM_MACD_EMA_RANGE_SHORT = 12  //Range of shorter EMA in MACD
	PARAM_MACD_EMA_RANGE_LONG  = 26  //Range of longer EMA in MACD
	PARAM_MACD_RANGE           = 9   //Range of MACD
	PARAM_DATA_RANGE           = 365 //Lookback range of data
	PARAM_MACD_MONITOR_RANGE   = 10  //Range of MACH to look back when evaluating the trend
)

//Init is to initiate the fields of StockMonitor
//Return: StockMonitor - the pointer of stock monitor
func (stockMonitor *StockMonitor) Init() *StockMonitor {
	stockMonitor.stockChances = make(map[string]*StockChance)
	return stockMonitor
}

//MonitorAllStock monitors all the stock
//Param: isMarketOpen - If market is open, the real time stock data will be included to calculate the statistics; Otherwise, the real time stock data is excluded; type bool
//Return: void
func (monitor *Monitor) MonitorAllStock() {
	bd := new(Board)
	stockMonitor := new(StockMonitor).Init()
	execution := func(symbol string) {
		monitor.MonitorStock(symbol, stockMonitor)
	}
	callback := func() {
		//Send email
		emailBody := ""
		emailBody += fmt.Sprintf("%10s%10s%10s%10s%10s%10s%10s%10s\r\n", "STATUS", "Trend", "Symbol", "Price", "BMID", "BUP", "KMID", "KUP")
		for _, chance := range stockMonitor.stockChances {
			emailBody += fmt.Sprintf("%10s", chance.Status)
			emailBody += fmt.Sprintf("%10s", chance.Trend)
			emailBody += fmt.Sprintf("%10s", chance.Symbol)
			emailBody += fmt.Sprintf("%10.2f", chance.Price)
			emailBody += fmt.Sprintf("%10.2f", chance.BollMid)
			emailBody += fmt.Sprintf("%10.2f", chance.BollUpper)
			emailBody += fmt.Sprintf("%10.2f", chance.KcMid)
			emailBody += fmt.Sprintf("%10.2f", chance.KcUpper)
			emailBody += "\r\n"
		}
		email := Email{
			senderId: EMAIL_SENDER,
			toIds:    []string{EMAIL_RECEIVER},
			subject:  "CardSharps Monitor Report",
			body:     emailBody,
			password: EMAIL_PASSWORD,
		}
		email.sendEmail()
		PrintMsgInConsole(MSGSYSTEM, LOGTYPE_SHUFFLER, "Completed one round of stock monitoring.")
	}
	bd.SymbolGame(SYMBOLTAG_STOCKSTAR, execution, callback)
}

//MonitorStock monitors one stock
//Param: symbol - stock symbol/quote
//Param: StockMonitor - pointer to the struct storing the status of monitoring
//Return: void
func (monitor *Monitor) MonitorStock(symbol string, stockMonitor *StockMonitor) {

	sharper := new(Sharper)
	yahooAPIManager := new(YahooAPIManager)
	stm := *stockMonitor

	//Get real time stock data
	stock, yahooAPIErr := yahooAPIManager.GetStockDataBySymbol(symbol)
	if yahooAPIErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_MONITOR, "YahooAPI Error: "+yahooAPIErr.Error())
		return
	}

	//Get historical stock data (reversed, 0 is previous day)
	histList, histErr := new(TblStockHist).SelectLastStockHistByCountAndBeforeDate(symbol, PARAM_DATA_RANGE, GetTimeInYYYYMMDD())
	if histErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_MONITOR, "Stock Hist Error: "+histErr.Error())
		return
	}
	if len(histList) < PARAM_CROSS_RANGE+PARAM_LOOKBACK_RANGE*2 {
		PrintMsgInConsole(MSGERROR, LOGTYPE_MONITOR, "Stock Hist Error: "+"Not enough Stock Hist for "+symbol)
		return
	}

	//Prepend real time stock to the head of historical stock data
	stockHist := RowStockHist{
		Symbol:      stock.Symbol,
		Date:        GetTimeInYYYYMMDD(),
		MarketOpen:  stock.RegularMarketOpen,
		MarketHigh:  stock.RegularMarketDayHigh,
		MarketLow:   stock.RegularMarketDayLow,
		MarketClose: stock.RegularMarketPrice,
		Volume:      stock.RegularMarketVolume,
	}
	histList = append([]RowStockHist{stockHist}, histList...)

	var closePriceList []float64

	for i := 0; i < len(histList); i++ {
		closePriceList = append(closePriceList, float64(histList[i].MarketClose))
	}

	//Calculate the parameters of Bollinger Band and Keltner Channel (reversed, index 0 is closest trading day)
	var bollMidArray []float32
	var bollUpperArray []float32
	var bollLowerArray []float32
	var kcMidArray []float32
	var kcUpperArray []float32
	var kcLowerArray []float32
	var macdArray []float64
	var diffArray []float64
	var macdErr error
	macdArray, diffArray, _, macdErr = sharper.CalcMACDStat(closePriceList, PARAM_MACD_EMA_RANGE_SHORT, PARAM_MACD_EMA_RANGE_LONG, PARAM_MACD_RANGE)
	if macdErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_MONITOR, "Calculate MACD Stat Error for "+symbol+": "+macdErr.Error())
		return
	}
	for i := 0; i < PARAM_CROSS_RANGE; i++ {
		bollMid, bollUpper, bollLower, kcMid, kcUpper, kcLower := sharper.CalcBollKcStat(histList, i, PARAM_CROSS_RANGE, PARAM_LOOKBACK_RANGE)
		bollMidArray = append(bollMidArray, bollMid)
		bollUpperArray = append(bollUpperArray, bollUpper)
		bollLowerArray = append(bollLowerArray, bollLower)
		kcMidArray = append(kcMidArray, kcMid)
		kcUpperArray = append(kcUpperArray, kcUpper)
		kcLowerArray = append(kcLowerArray, kcLower)
	}

	stc := StockChance{
		Symbol:    stock.Symbol,
		Price:     stock.RegularMarketPrice,
		BollMid:   bollMidArray[0],
		BollUpper: bollUpperArray[0],
		BollLower: bollLowerArray[0],
		KcMid:     kcMidArray[0],
		KcUpper:   kcUpperArray[0],
		KcLower:   kcLowerArray[0],
		Trend:     "",
		Status:    "",
	}

	//Analyze results
	var daysBeforeMacdFlip int
	var daysBeforeMacdNotFlip int
	for i := 1; i < MinInt(len(macdArray), PARAM_MACD_MONITOR_RANGE+1); i++ {
		if macdArray[i]*macdArray[0] < 0 {
			daysBeforeMacdFlip++
		} else {
			daysBeforeMacdNotFlip++
		}
	}

	var smallerPercentChangeCount int
	var largerVolumeCount int
	currentPercentChange := math.Abs(float64(histList[0].MarketClose - histList[0].MarketOpen/histList[0].MarketOpen))
	currentVolume := histList[0].Volume
	var secondLargeVolume int
	for i := 1; i < len(histList) && i < 21; i++ {
		newPercentChange := math.Abs(float64(histList[i].MarketClose - histList[i].MarketOpen/histList[i].MarketOpen))
		newVolume := histList[i].Volume
		if currentPercentChange > newPercentChange {
			smallerPercentChangeCount++
		}
		if currentVolume < newVolume {
			largerVolumeCount++
		}
		if secondLargeVolume < newVolume {
			secondLargeVolume = newVolume
		}
	}

	//TODO
	if (macdArray[0] > 0 && macdArray[1] < 0) && daysBeforeMacdFlip >= 10 {
		if _, ok := stm.stockChances[stock.Symbol]; ok {
			stc.Status = "KEEP"
			stc.Trend = "HOT"
		} else {
			stc.Status = "NEW"
			stc.Trend = "HOT"
		}
		stm.stockChances[stock.Symbol] = &stc
		// } else if (macdArray[0] < 0 && macdArray[1] > 0) && daysBeforeMacdFlip >= 10 {
		// 	if _, ok := stm.stockChances[stock.Symbol]; ok {
		// 		stc.Status = "KEEP"
		// 		stc.Trend = "X-Down"
		// 	} else {
		// 		stc.Status = "NEW"
		// 		stc.Trend = "X-Down"
		// 	}
		// 	stm.stockChances[stock.Symbol] = &stc
	} else if diffArray[0] > 0 && diffArray[1] < 0 {
		if _, ok := stm.stockChances[stock.Symbol]; ok {
			stc.Status = "KEEP"
			stc.Trend = "CROSS"
		} else {
			stc.Status = "NEW"
			stc.Trend = "CROSS"
		}
		stm.stockChances[stock.Symbol] = &stc
		// } else if bollUpperArray[0] > kcUpperArray[0] && bollUpperArray[PARAM_CROSS_RANGE-1] < kcUpperArray[PARAM_CROSS_RANGE-1] && stock.RegularMarketPrice > bollUpperArray[0] {
		// 	if _, ok := stm.stockChances[stock.Symbol]; ok {
		// 		stc.Status = "KEEP"
		// 		stc.Trend = "UP"
		// 	} else {
		// 		stc.Status = "NEW"
		// 		stc.Trend = "UP"
		// 	}
		// 	stm.stockChances[stock.Symbol] = &stc
		// 	fmt.Printf("Long %s %.2f\r\n", stock.Symbol, stock.RegularMarketPrice)
		// } else if bollLowerArray[0] < kcLowerArray[0] && bollLowerArray[PARAM_CROSS_RANGE-1] > kcLowerArray[PARAM_CROSS_RANGE-1] && stock.RegularMarketPrice < bollLowerArray[0] {
		// 	if _, ok := stm.stockChances[stock.Symbol]; ok {
		// 		stc.Status = "KEEP"
		// 		stc.Trend = "DOWN"
		// 	} else {
		// 		stc.Status = "NEW"
		// 		stc.Trend = "DOWN"
		// 	}
		// 	stm.stockChances[stock.Symbol] = &stc
		// 	fmt.Printf("Short %s %.2f\r\n", stock.Symbol, stock.RegularMarketPrice)
		// } else if (kcUpperArray[0]-bollUpperArray[0]+bollLowerArray[0]-kcLowerArray[0])/(bollUpperArray[0]-bollLowerArray[0]) > 2 {
		// 	if _, ok := stm.stockChances[stock.Symbol]; ok {
		// 		stc.Status = "KEEP"
		// 		stc.Trend = "SQUEEZE"
		// 	} else {
		// 		stc.Status = "NEW"
		// 		stc.Trend = "SQUEEZE"
		// 	}
		// 	stm.stockChances[stock.Symbol] = &stc
		// 	fmt.Printf("Squeeze %s %.2f\r\n", stock.Symbol, stock.RegularMarketPrice)
		// } else if macdArray[0] > 0 && macdArray[1] < 0 && kcUpperArray[0]-kcLowerArray[0] > bollUpperArray[0]-bollLowerArray[0] {
		// 	if _, ok := stm.stockChances[stock.Symbol]; ok {
		// 		stc.Status = "KEEP"
		// 		stc.Trend = "CROSS"
		// 	} else {
		// 		stc.Status = "NEW"
		// 		stc.Trend = "CROSS"
		// 	}
		// 	stm.stockChances[stock.Symbol] = &stc
		// 	fmt.Printf("Cross-1 %s %.2f\r\n", stock.Symbol, stock.RegularMarketPrice)
	} else if smallerPercentChangeCount <= 3 && largerVolumeCount < 2 {
		if _, ok := stm.stockChances[stock.Symbol]; ok {
			stc.Status = "KEEP"
			stc.Trend = "EJECT"
		} else {
			stc.Status = "NEW"
			stc.Trend = "EJECT"
		}
		stm.stockChances[stock.Symbol] = &stc
	} else {
		if _, ok := stm.stockChances[stock.Symbol]; ok {
			stc.Status = "REMOVE"
			stm.stockChances[stock.Symbol] = &stc
		}
	}

	PrintMsgInConsole(MSGSYSTEM, LOGTYPE_SHUFFLER, "Completed one round of stock monitor for "+symbol)
}

func (monitor *Monitor) MonitorPCR(symbol string, date int64) {
	expDate, callVol, putVol, callOi, putOi, err := new(TblOptionData).SelectOptionDataVolumeBySymbolAndDate(symbol, date)
	if err != nil {
		panic(err)
	}
	sort.Slice(expDate, func(i, j int) bool { return expDate[i] < expDate[j] })

	//Send email
	emailBody := ""
	emailBody += fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s\r\n", "Date", "ExpDate", "CPR-Vol", "CPR-OI", "CallVol", "PutVol", "CallOI", "PutOI", "Vol Hint", "OI Hint")

	for _, exp := range expDate {
		cprVol := float64(callVol[exp]) / float64(putVol[exp])
		cprOi := float64(callOi[exp]) / float64(putOi[exp])
		volHint := ""
		oiHint := ""
		if cprVol < 1 {
			volHint = "Buyer Bear"
		} else {
			volHint = "Buyer Bull"
		}
		if cprOi < 1 {
			oiHint = "Holder Bull"
		} else {
			oiHint = "Holder Bear"
		}
		emailBody += fmt.Sprintf("%d", date)
		emailBody += fmt.Sprintf(" %d", exp)
		emailBody += fmt.Sprintf(" %.2f", cprVol)
		emailBody += fmt.Sprintf(" %.2f", cprOi)
		emailBody += fmt.Sprintf(" %d", callVol[exp])
		emailBody += fmt.Sprintf(" %d", putVol[exp])
		emailBody += fmt.Sprintf(" %d", callOi[exp])
		emailBody += fmt.Sprintf(" %d", putOi[exp])
		emailBody += fmt.Sprintf(" %s", volHint)
		emailBody += fmt.Sprintf(" %s", oiHint)
		emailBody += "\r\n"
	}

	email := Email{
		senderId: EMAIL_SENDER,
		toIds:    []string{EMAIL_RECEIVER},
		subject:  "CardSharps Option Monitor Report",
		body:     emailBody,
		password: EMAIL_PASSWORD,
	}
	email.sendEmail()

	return
}
