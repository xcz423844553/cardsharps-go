package main

import (
	"errors"
	"fmt"
	"time"
)

//Monitor is a library to monitor the market status
type Monitor2 struct {
}

//DailyReport stores the reports of daily monitoring
type DailyReport struct {
	reports      []Report
	optionReport []OptionReport
}

//Report stores all the information of a opportunity
type Report struct {
	Symbol    string
	Price     float32
	BollMid   float32
	BollUpper float32
	BollLower float32
	KcMid     float32
	KcUpper   float32
	KcLower   float32
	Note      string
}

//OptionReport stores the information of the option
type OptionReport struct {
	Symbol             string
	OptionType         string
	ExpirationDate     int64
	Strike             float32
	CurrentPrice       float32
	DoublePrice        []float32
	DoublePriceChange  []float32
	CurrentOptionPrice float32
	DoubleOptionPrice  float32
	Atr                float32
	AtrPercent         float32
	Ma                 float32
	Std                float32
	StdPercent         float32
}

//PrintInEmail returns the information in the report in a string to be sent via email
func (report *Report) PrintInEmail() string {
	return fmt.Sprintf("%s\r\n", report.Symbol)
}

const (
	//ParamCrossRange is Range to determine Boll-Kc trend
	ParamCrossRange = 5
	//ParamLookbackRange is Range to look back to calculate EMA
	ParamLookbackRange = 20
	//ParamMacdEmaRangeShort is Range of shorter EMA in MACD
	ParamMacdEmaRangeShort = 12
	//ParamMacdEmaRangeLong is Range of longer EMA in MACD
	ParamMacdEmaRangeLong = 26
	//ParamMacdRange is Range of MACD
	ParamMacdRange = 9
	//ParamDataRange is Lookback range of data
	ParamDataRange = 365
	//ParamMacdMonitorRange is Range of MACH to look back when evaluating the trend
	ParamMacdMonitorRange = 10
)

//RunDailyMonitor runs daily monitoring program
func (monitor *Monitor2) RunDailyMonitor(tag string) {
	board := new(Board2)
	dailyReport := new(DailyReport)
	executeFunc := func(symbol string) {
		//Add models here
		monitor.MonitorModel1(symbol, GetTimeInYYYYMMDD64(), dailyReport)
		// monitor.MonitorModel2(symbol, GetTimeInYYYYMMDD64(), dailyReport)
	}
	callback := func() {
		monitor.SendDailyMonitor(dailyReport)
	}
	board.StartGame(tag, executeFunc, callback)
}

//SendDailyMonitor sends the reports of daily monitoring program out
func (monitor *Monitor2) SendDailyMonitor(dailyReport *DailyReport) {
	fmt.Println("send email")
	email := Email{
		senderId: EMAIL_SENDER,
		toIds:    []string{EMAIL_RECEIVER},
		subject:  "CardSharps Daily Report",
		body:     "",
		password: EMAIL_PASSWORD,
	}
	email.SendEmailInTemplate()
	return

	// Symbol    string
	// Price     float32
	// BollMid   float32
	// BollUpper float32
	// BollLower float32
	// KcMid     float32
	// KcUpper   float32
	// KcLower   float32
	// Note      string

	// emailBody := ""
	// emailBody += fmt.Sprintf("%10s%10s%10s%10s%10s%10s%10s%10s\r\n", "STATUS", "Trend", "Symbol", "Price", "BMID", "BUP", "KMID", "KUP")
	// for _, report := range dailyReport.reports {
	// 	emailBody += report.PrintInEmail()
	// 	emailBody += "\r\n"
	// }
	// email := Email{
	// 	senderId: EMAIL_SENDER,
	// 	toIds:    []string{EMAIL_RECEIVER},
	// 	subject:  "CardSharps Daily Report",
	// 	body:     emailBody,
	// 	password: EMAIL_PASSWORD,
	// }
	// email.sendEmail()
	// fmt.Println("Daily Report is sent via email.")
}

//MonitorModel1 calculates the MACD, Bollinger Bands, and Keltner Channel
func (monitor *Monitor2) MonitorModel1(symbol string, currentDate int64, dailyReport *DailyReport) error {
	sharper := new(Sharper2)
	yahooApi := new(YahooApi)

	//Get real time stock data
	quote, yahooApiErr := yahooApi.GetQuote(symbol)
	if yahooApiErr != nil {
		return yahooApiErr
	}

	//Get historical stock data (reversed, 0 is previous day)
	histList, histErr := new(DaoStockHist).SelectLastNumberStockHistBeforeDate(symbol, ParamDataRange, currentDate)
	if histErr != nil {
		return histErr
	}
	if len(histList) < ParamCrossRange+ParamLookbackRange*2 {
		return errors.New("Stock Hist Error: " + "Not enough Stock Hist for " + symbol)
	}

	//Prepare real time stock to the head of historical stock data
	stockHist := RowStockHist2{
		Symbol:      quote.GetSymbol(),
		Date:        currentDate,
		MarketOpen:  quote.GetMarketOpen(),
		MarketHigh:  quote.GetMarketHigh(),
		MarketLow:   quote.GetMarketLow(),
		MarketClose: quote.GetMarketClose(),
		Volume:      quote.GetVolume(),
	}
	histList = append([]RowStockHist2{stockHist}, histList...)

	var closePriceList []float32

	for i := 0; i < len(histList); i++ {
		closePriceList = append(closePriceList, float32(histList[i].MarketClose))
	}

	//Calculate the parameters of Bollinger Band and Keltner Channel (reversed, index 0 is closest trading day)
	var bollMidArray []float32
	var bollUpperArray []float32
	var bollLowerArray []float32
	var kcMidArray []float32
	var kcUpperArray []float32
	var kcLowerArray []float32
	var macdArray []float32
	//var diffArray []float32
	var deaArray []float32
	var macdErr error
	macdArray, _, deaArray, macdErr = sharper.CalcMACDStat(closePriceList, ParamMacdEmaRangeShort, ParamMacdEmaRangeLong, ParamMacdRange)
	if macdErr != nil {
		return errors.New("Error encountered while calculating MACD for " + symbol + ": " + macdErr.Error())
	}
	for i := 0; i < ParamCrossRange; i++ {
		bollMid, bollUpper, bollLower, kcMid, kcUpper, kcLower := sharper.CalcBollKcStat(histList, int64(i), int64(ParamCrossRange), int64(ParamLookbackRange))
		bollMidArray = append(bollMidArray, bollMid)
		bollUpperArray = append(bollUpperArray, bollUpper)
		bollLowerArray = append(bollLowerArray, bollLower)
		kcMidArray = append(kcMidArray, kcMid)
		kcUpperArray = append(kcUpperArray, kcUpper)
		kcLowerArray = append(kcLowerArray, kcLower)
	}

	// stc := StockChance{
	// 	Symbol:    stock.Symbol,
	// 	Price:     stock.RegularMarketPrice,
	// 	BollMid:   bollMidArray[0],
	// 	BollUpper: bollUpperArray[0],
	// 	BollLower: bollLowerArray[0],
	// 	KcMid:     kcMidArray[0],
	// 	KcUpper:   kcUpperArray[0],
	// 	KcLower:   kcLowerArray[0],
	// 	Trend:     "",
	// 	Status:    "",
	// }

	//Analyze results
	var dayBeforeMacdFlip int
	var dayWhenMacdFlip int
	dayBeforeMacdFlip = MinInt(len(macdArray), ParamMacdMonitorRange+1)
	dayWhenMacdFlip = MinInt(len(macdArray), ParamMacdMonitorRange+1)
	for i := 0; i < MinInt(len(macdArray), ParamMacdMonitorRange+1); i++ {
		if (macdArray[i] >= 0 && macdArray[0] <= 0) || (macdArray[i] <= 0 && macdArray[0] >= 0) {
			dayWhenMacdFlip = i
			break
		}
	}
	for i := dayWhenMacdFlip; i < MinInt(len(macdArray), ParamMacdMonitorRange+1); i++ {
		if (macdArray[i] >= 0 && macdArray[0] >= 0) || (macdArray[i] <= 0 && macdArray[0] <= 0) {
			dayBeforeMacdFlip = i
			break
		}
	}

	var deaLargerCount int
	var deaLowerCount int
	for i := 0; i < len(deaArray); i++ {
		if deaArray[i] > deaArray[0] {
			deaLargerCount++
		} else {
			deaLowerCount++
		}
	}

	if dayWhenMacdFlip < 3 && dayBeforeMacdFlip-dayWhenMacdFlip > 5 && (deaLargerCount > deaLowerCount*4 || deaLargerCount*4 < deaLowerCount) {
		if macdArray[0] > 0 {
			fmt.Printf("%s is a potential up. %.2f\r\n", symbol, macdArray[0])
		} else {
			fmt.Printf("%s is a potential down. %.2f\r\n", symbol, macdArray[0])
		}
	}
	return nil
}

//MonitorModel2 calculates the options
func (monitor *Monitor2) MonitorModel2(symbol string, currentDate int64, dailyReport *DailyReport) error {
	sharper := new(Sharper2)
	yahooApi := new(YahooApi)

	//Get real time stock data and option data
	options, quote, _, yahooApiErr := yahooApi.GetAllOptionChainStockAndExpDate(symbol)
	if yahooApiErr != nil {
		return yahooApiErr
	}

	//Get historical stock data (reversed, 0 is previous day)
	histList, histErr := new(DaoStockHist).SelectLastNumberStockHistBeforeDate(symbol, ParamDataRange, currentDate)
	if histErr != nil {
		return histErr
	}
	if len(histList) < ParamCrossRange+ParamLookbackRange*2 {
		return errors.New("Stock Hist Error: " + "Not enough Stock Hist for " + symbol)
	}

	//Prepare real time stock to the head of historical stock data
	stockHist := RowStockHist2{
		Symbol:      quote.GetSymbol(),
		Date:        currentDate,
		MarketOpen:  quote.GetMarketOpen(),
		MarketHigh:  quote.GetMarketHigh(),
		MarketLow:   quote.GetMarketLow(),
		MarketClose: quote.GetMarketClose(),
		Volume:      quote.GetVolume(),
	}
	histList = append([]RowStockHist2{stockHist}, histList...)

	var closePriceList []float32

	for i := 0; i < len(histList); i++ {
		closePriceList = append(closePriceList, float32(histList[i].MarketClose))
	}

	ma, sd := sharper.CalcMAAndSDStat(histList, 0, ParamLookbackRange)
	atr := sharper.CalcATRStat(histList, 0, ParamLookbackRange)
	priceUpperBound := quote.GetRegularMarketPrice() + atr
	priceLowerBound := quote.GetRegularMarketPrice() - atr
	fmt.Printf("MA:%.2f   Std:%.2f   ATR:%.2f\r\n", ma, sd, atr)
	var targetMultiple float32 = 2.0
	var interestRate float32 = 0.02
	var dividendYield float32 = 0.0

	for _, option := range options {
		dateCur := time.Date(int(currentDate/10000), time.Month(int(currentDate%10000/100)), int(currentDate%1000000), 0, 0, 0, 0, time.UTC)
		dateExp := time.Date(int(option.GetExpiration()/10000), time.Month(int(option.GetExpiration()%10000/100)), int(option.GetExpiration()%1000000), 0, 0, 0, 0, time.UTC)
		currentDaysToExpire := int(dateExp.Sub(dateCur).Hours() / 24)
		var doublePrice []float32
		var doublePriceChange []float32
		nextDayDoublePrice := sharper.CalcBlackScholesTrick(targetMultiple, option.GetOptionType(), quote.GetRegularMarketPrice(), option.GetStrike(),
			option.GetImpliedVolatility(), interestRate, dividendYield, currentDaysToExpire, currentDaysToExpire-1)
		if nextDayDoublePrice > priceLowerBound && nextDayDoublePrice < priceUpperBound {
			for futureDaysToExpire := currentDaysToExpire; futureDaysToExpire > 0; futureDaysToExpire-- {
				dp := sharper.CalcBlackScholesTrick(targetMultiple, option.GetOptionType(), quote.GetRegularMarketPrice(), option.GetStrike(),
					option.GetImpliedVolatility(), interestRate, dividendYield, currentDaysToExpire, futureDaysToExpire)
				doublePrice = append(doublePrice, dp)
				doublePriceChange = append(doublePriceChange, (dp-quote.GetRegularMarketPrice())/quote.GetRegularMarketPrice())

			}
			rpt := OptionReport{
				Symbol:             quote.GetSymbol(),
				OptionType:         option.GetOptionType(),
				ExpirationDate:     option.GetExpiration(),
				Strike:             option.GetStrike(),
				CurrentPrice:       quote.GetRegularMarketPrice(),
				DoublePrice:        doublePrice,
				DoublePriceChange:  doublePriceChange,
				CurrentOptionPrice: option.GetAsk(),
				DoubleOptionPrice:  2 * option.GetAsk(),
				Atr:                atr,
				AtrPercent:         atr / quote.GetRegularMarketPrice() * 100,
				Ma:                 ma,
				Std:                sd,
				StdPercent:         sd / quote.GetRegularMarketPrice() * 100,
			}
			dailyReport.optionReport = append(dailyReport.optionReport, rpt)
		}
	}
	return nil
}

// func (monitor *Monitor2) MonitorPCR(symbol string, date int64) {
// 	expDate, callVol, putVol, callOi, putOi, err := new(TblOptionData).SelectOptionDataVolumeBySymbolAndDate(symbol, date)
// 	if err != nil {
// 		panic(err)
// 	}
// 	sort.Slice(expDate, func(i, j int) bool { return expDate[i] < expDate[j] })

// 	//Send email
// 	emailBody := ""
// 	emailBody += fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s\r\n", "Date", "ExpDate", "CPR-Vol", "CPR-OI", "CallVol", "PutVol", "CallOI", "PutOI", "Vol Hint", "OI Hint")

// 	for _, exp := range expDate {
// 		cprVol := float64(callVol[exp]) / float64(putVol[exp])
// 		cprOi := float64(callOi[exp]) / float64(putOi[exp])
// 		volHint := ""
// 		oiHint := ""
// 		if cprVol < 1 {
// 			volHint = "Buyer Bear"
// 		} else {
// 			volHint = "Buyer Bull"
// 		}
// 		if cprOi < 1 {
// 			oiHint = "Holder Bull"
// 		} else {
// 			oiHint = "Holder Bear"
// 		}
// 		emailBody += fmt.Sprintf("%d", date)
// 		emailBody += fmt.Sprintf(" %d", exp)
// 		emailBody += fmt.Sprintf(" %.2f", cprVol)
// 		emailBody += fmt.Sprintf(" %.2f", cprOi)
// 		emailBody += fmt.Sprintf(" %d", callVol[exp])
// 		emailBody += fmt.Sprintf(" %d", putVol[exp])
// 		emailBody += fmt.Sprintf(" %d", callOi[exp])
// 		emailBody += fmt.Sprintf(" %d", putOi[exp])
// 		emailBody += fmt.Sprintf(" %s", volHint)
// 		emailBody += fmt.Sprintf(" %s", oiHint)
// 		emailBody += "\r\n"
// 	}

// 	email := Email{
// 		senderId: EMAIL_SENDER,
// 		toIds:    []string{EMAIL_RECEIVER},
// 		subject:  "CardSharps Option Monitor Report",
// 		body:     emailBody,
// 		password: EMAIL_PASSWORD,
// 	}
// 	email.sendEmail()

// 	return
// }
