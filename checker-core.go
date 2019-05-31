package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Checker struct {
	dragonTail        map[string]*DragonTail
	addedDragonTail   map[string]*struct{}
	removedDragonTail map[string]*struct{}
	shield            map[string]*Shield
}

//Dragon Tail is the stock that has the dragon tail form to look at
type DragonTail struct {
	Symbol       string
	CurrentPrice float32
	MA60         float32
	MA120        float32
	Trend        string
}

//Shield is the stock that lift up or down quickly
type Shield struct {
	Symbol        string
	CurrentPrice  float32
	PercentChange float32
}

func (checker *Checker) InitChecker() *Checker {
	checker.dragonTail = make(map[string]*DragonTail)
	checker.addedDragonTail = make(map[string]*struct{})
	checker.removedDragonTail = make(map[string]*struct{})
	checker.shield = make(map[string]*Shield)
	return checker
}

func (checker Checker) runChecker() {
	yahooApiManager := new(YahooApiManager)
	tblLogError := new(TblLogError)
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	go func() {
		symbols := []string{"MSFT", "AMZN", "AAPL", "GOOGL", "GOOG", "FB", "JPM", "JNJ", "XOM", "V", "WMT", "BAC", "PG", "MA"}
		volumes := make([]int, len(symbols))
		totalVolumeInMinute := 0
		averageVolumes := []int{}
		for ; true; <-ticker.C {
			if !isMarketOpen() {
				fmt.Println("Market not open")
			} else {
				fmt.Println("Tick at ", time.Now())
				totalVolumeInMinute = 0
				for index, symbol := range symbols {
					_, quote, _, yahooErr := yahooApiManager.GetOptionsAndStockDataBySymbol(symbol)
					if yahooErr != nil {
						tblLogError.InsertLogError(LOGTYPE_CHECKER, yahooErr.Error())
					}
					fmt.Println(quote.Symbol, " ", quote.RegularMarketVolume)
					totalVolumeInMinute += quote.RegularMarketVolume - volumes[index]
					volumes[index] = quote.RegularMarketVolume
				}
				if len(averageVolumes) < 30 {
					averageVolumes = append(averageVolumes, totalVolumeInMinute)
				} else {
					averageVol := AverageInt(averageVolumes)
					if totalVolumeInMinute >= 2*averageVol {
						volList := ""
						for _, vol := range averageVolumes {
							volList += strconv.Itoa(vol) + "\r\n"
						}
						email := Email{
							senderId: EMAIL_SENDER,
							toIds:    []string{EMAIL_RECEIVER},
							subject:  "CardSharps Notification",
							body:     "High Volume Alert @ " + time.Now().String() + "\r\n" + volList,
							password: EMAIL_PASSWORD,
						}
						email.sendEmail()
					}
					averageVolumes = append(averageVolumes, totalVolumeInMinute)
					averageVolumes = averageVolumes[1:]
				}
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":9999", nil))
	// time.Sleep(8 * time.Hour)
}

var checkerWaitGroup sync.WaitGroup
var checkerMut sync.Mutex

func (checker Checker) runChecker2() {
	yahooApiManager := new(YahooApiManager)
	tblLogError := new(TblLogError)
	ticker := time.NewTicker(TICKER_VOLUME_CHECKER)
	defer ticker.Stop()
	go func() {
		symbols := []string{"MSFT", "AMZN", "AAPL", "GOOGL", "GOOG", "FB", "JPM", "JNJ", "XOM", "V", "WMT", "BAC", "PG", "MA"}
		//"CSCO", "PFE", "DIS", "VZ", "T", "CVX", "UNH", "HD", "KO", "MRK", "INTC", "WFC", "ORCL", "CMCSA", "PEP", "NFLX", "MCD", "C"}
		volumes := make([]int, len(symbols))
		totalVolumeInMinute := 0
		averageVolumes := []int{}
		for ; true; <-ticker.C {
			if !isMarketOpen() && !BYPASS_MARKET_STATUS {
				fmt.Println("Market not open yet")
			} else {
				fmt.Println("Tick at ", time.Now())
				totalVolumeInMinute = 0
				for index, symbol := range symbols {
					checkerWaitGroup.Add(1)
					go func(volumes []int, totalVolumeInMinute *int, index int, symbol string) {
						_, quote, _, yahooErr := yahooApiManager.GetOptionsAndStockDataBySymbol(symbol)
						if yahooErr != nil {
							tblLogError.InsertLogError(LOGTYPE_CHECKER, yahooErr.Error())
						}
						fmt.Println("(" + strconv.Itoa(index) + ")" + symbol + " " + strconv.Itoa(quote.RegularMarketVolume))
						checkerMut.Lock()
						*totalVolumeInMinute += quote.RegularMarketVolume - volumes[index]
						checkerMut.Unlock()
						volumes[index] = quote.RegularMarketVolume
						checkerWaitGroup.Done()
					}(volumes, &totalVolumeInMinute, index, symbol)
				}
				checkerWaitGroup.Wait()
				if len(averageVolumes) < MIN_LENGTH_OF_MINUTE_VOLUME_CHECKER {
					averageVolumes = append(averageVolumes, totalVolumeInMinute)
				} else {
					averageVol := AverageInt(averageVolumes)
					if totalVolumeInMinute >= MULTI_THRESHOLD_VOLUME_CHECKER*averageVol {
						volList := ""
						ratio := fmt.Sprintf("%.1f", float32(totalVolumeInMinute)/float32(averageVol))
						for index, vol := range averageVolumes {
							volList += "(" + strconv.Itoa(index) + ") " + strconv.Itoa(vol) + "\r\n"
						}
						volList += "(0) " + strconv.Itoa(totalVolumeInMinute) + " @ " + ratio + "\r\n"
						email := Email{
							senderId: EMAIL_SENDER,
							toIds:    []string{EMAIL_RECEIVER},
							subject:  "CardSharps Notification",
							body:     "High Volume Alert (" + ratio + ") @ " + time.Now().Format("15:04 Mon Jan _2 2006") + "\r\n" + volList,
							password: EMAIL_PASSWORD,
						}
						email.sendEmail()
					}
					averageVolumes = append(averageVolumes, totalVolumeInMinute)
					if len(averageVolumes) > MAX_LENGTH_OF_MINUTE_VOLUME_CHECKER {
						averageVolumes = averageVolumes[1:]
					}
				}
				fmt.Println(averageVolumes)
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":9999", nil))
}

func (checker Checker) runChecker3() {
	yahooApiManager := new(YahooApiManager)
	tblLogError := new(TblLogError)
	ticker := time.NewTicker(TICKER_MONEY_CHECKER)
	defer ticker.Stop()
	go func() {
		symbols := []string{"MSFT", "AMZN", "AAPL", "GOOGL", "GOOG", "FB", "JPM", "JNJ", "XOM", "V", "WMT", "BAC", "PG", "MA",
			"CSCO", "PFE", "DIS", "VZ", "T", "CVX", "UNH", "HD", "KO", "MRK", "INTC", "WFC", "ORCL", "CMCSA", "PEP", "NFLX", "MCD", "C",
			"ADBE", "PM", "ABT", "NKE", "PYPL", "UNP", "HON", "IBM", "CRM", "ABBV", "MDT", "UTX", "ACN", "LLY", "COST", "AVGO", "TMO",
			"AMGN", "TXN", "AXP", "MO", "LMT", "MMM", "SBUX", "NEE", "DHR", "QCOM", "NVDA", "AMT", "LOW", "GE", "UPS", "CHTR", "GILD",
			"USB", "BMY", "MDLZ", "MS", "GS", "COP", "CAT", "ADP"}
		volumes := make([]int64, len(symbols))
		symbolVolumes := make([]int, len(symbols))
		symbolPrices := make([]float32, len(symbols))
		totalMoneyInMinute := int64(0)
		averageMoneyArray := []int64{}
		for ; true; <-ticker.C {
			if !isMarketOpen() && !BYPASS_MARKET_STATUS {
				fmt.Println("Market not open yet")
			} else {
				fmt.Println("Tick at ", time.Now())
				totalMoneyInMinute = 0
				for index, symbol := range symbols {
					checkerWaitGroup.Add(1)
					go func(volumes []int64, symbolVolumes []int, symbolPrices []float32, totalMoneyInMinute *int64, index int, symbol string) {
						_, quote, _, yahooErr := yahooApiManager.GetOptionsAndStockDataBySymbol(symbol)
						if yahooErr != nil {
							tblLogError.InsertLogError(LOGTYPE_CHECKER, yahooErr.Error())
						}
						fmt.Println("(" + strconv.Itoa(index) + ")" + symbol + " " + strconv.Itoa(quote.RegularMarketVolume))
						checkerMut.Lock()
						*totalMoneyInMinute += int64(float32(int64(quote.RegularMarketVolume)-volumes[index]) * quote.RegularMarketPrice)
						checkerMut.Unlock()
						volumes[index] = int64(quote.RegularMarketVolume)
						symbolVolumes[index] = quote.RegularMarketVolume
						symbolPrices[index] = quote.RegularMarketPrice
						checkerWaitGroup.Done()
					}(volumes, symbolVolumes, symbolPrices, &totalMoneyInMinute, index, symbol)
				}
				checkerWaitGroup.Wait()
				//after the volume and price are updated, insert the data into database
				go func(symbolVolumes []int, symbolPrices []float32) {
					new(TblStockChecker1).InsertStockCheckerData(GetTimeInYYYYMMDD(), GetTimeInt(), symbolVolumes, symbolPrices)
				}(symbolVolumes[:], symbolPrices[:])
				if len(averageMoneyArray) < MIN_LENGTH_OF_MINUTE_MONEY_CHECKER {
					averageMoneyArray = append(averageMoneyArray, totalMoneyInMinute)
				} else {
					averageMoney := AverageInt64(averageMoneyArray)
					if totalMoneyInMinute >= int64(MULTI_THRESHOLD_MONEY_CHECKER*float32(averageMoney)) {
						volList := ""
						ratio := fmt.Sprintf("%.1f", float32(totalMoneyInMinute)/float32(averageMoney))
						for index, vol := range averageMoneyArray {
							volList += "(" + strconv.Itoa(index) + ") " + strconv.FormatInt(vol, 10) + "\r\n"
						}
						volList += "(0) " + strconv.FormatInt(totalMoneyInMinute, 10) + " @ " + ratio + "\r\n"
						email := Email{
							senderId: EMAIL_SENDER,
							toIds:    []string{EMAIL_RECEIVER},
							subject:  "CardSharps Transaction Notification",
							body:     "High Transaction Alert (" + ratio + ") @ " + time.Now().Format("15:04 Mon Jan _2 2006") + "\r\n" + volList,
							password: EMAIL_PASSWORD,
						}
						email.sendEmail()
					}
					averageMoneyArray = append(averageMoneyArray, totalMoneyInMinute)
					if len(averageMoneyArray) > MAX_LENGTH_OF_MINUTE_MONEY_CHECKER {
						averageMoneyArray = averageMoneyArray[1:]
					}
				}
				fmt.Println(averageMoneyArray)
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":7777", nil))
}

var spyCheckerWaitGroup sync.WaitGroup
var spyCheckerMut sync.Mutex

func (checker Checker) runSpyChecker(symbol string) {
	yahooApiManager := new(YahooApiManager)
	tblLogError := new(TblLogError)
	ticker := time.NewTicker(TICKER_VOLUME_CHECKER)
	defer ticker.Stop()
	go func() {
		_, _, expDates, expErr := yahooApiManager.GetOptionsAndStockDataBySymbol(symbol)
		if expErr != nil {
			PrintMsgInConsole(MSGERROR, LOGTYPE_SPY_CHECKER, expErr.Error())
			return
		}
		volumeMap := make(map[string]int)
		totalVolumeInMinute := 0
		totalCallVolumeInMinute := 0
		totalPutVolumeInMinute := 0
		averageVolumeArray := []int{}
		for ; true; <-ticker.C {
			if !isMarketOpen() && !BYPASS_MARKET_STATUS {
				fmt.Println("Market not open yet")
			} else {
				fmt.Println("Tick at ", time.Now())
				totalVolumeInMinute = 0
				totalCallVolumeInMinute = 0
				totalPutVolumeInMinute = 0
				for index, expDate := range expDates {
					if expDate > (time.Now().Unix() + UNIX_TWO_WEEK) {
						continue
					}
					spyCheckerWaitGroup.Add(1)
					go func(volumeMap map[string]int, totalVolumeInMinute *int, totalCallVolumeInMinute *int, totalPutVolumeInMinute *int, index int, symbol string, expDate int64) {
						options, _, _, yahooErr := yahooApiManager.GetOptionsAndStockDataBySymbolAndExpDate(symbol, expDate)
						if yahooErr != nil {
							tblLogError.InsertLogError(LOGTYPE_SPY_CHECKER, yahooErr.Error())
						}
						spyCheckerMut.Lock()
						for _, option := range options {
							vol, _ := volumeMap[option.ContractSymbol]
							volumeMap[option.ContractSymbol] = option.Volume
							*totalVolumeInMinute += option.Volume - vol
							if ([]rune(option.ContractSymbol))[9] == 'C' {
								*totalCallVolumeInMinute += option.Volume - vol
							} else {
								*totalPutVolumeInMinute += option.Volume - vol
							}
							fmt.Println("(" + strconv.Itoa(index) + ")" + option.ContractSymbol + " " + strconv.Itoa(*totalVolumeInMinute) + " " + strconv.Itoa(*totalCallVolumeInMinute) + " " + strconv.Itoa(*totalPutVolumeInMinute))
						}
						spyCheckerMut.Unlock()
						spyCheckerWaitGroup.Done()
					}(volumeMap, &totalVolumeInMinute, &totalCallVolumeInMinute, &totalPutVolumeInMinute, index, symbol, expDate)
				}
				spyCheckerWaitGroup.Wait()
				if len(averageVolumeArray) < MIN_LENGTH_OF_MINUTE_MONEY_CHECKER {
					averageVolumeArray = append(averageVolumeArray, totalVolumeInMinute)
				} else {
					averageVolume := AverageInt(averageVolumeArray)
					var ratio1 float32
					var ratio2 float32
					if totalPutVolumeInMinute != 0 {
						ratio1 = float32(totalCallVolumeInMinute) / float32(totalPutVolumeInMinute)
					}
					if totalCallVolumeInMinute != 0 {
						ratio2 = float32(totalPutVolumeInMinute) / float32(totalCallVolumeInMinute)
					}
					if totalVolumeInMinute >= averageVolume && (ratio1 > 2 || ratio2 > 2) { //MULTI_THRESHOLD_VOLUME_CHECKER*
						volList := ""
						//ratio := fmt.Sprintf("%.1f", float32(totalVolumeInMinute)/float32(averageVolume))
						for index, vol := range averageVolumeArray {
							volList += "(" + strconv.Itoa(index) + ") " + strconv.Itoa(vol) + "\r\n"
						}
						volList += "(-) " + strconv.Itoa(totalVolumeInMinute) + " @ Call:" + fmt.Sprintf("%.1f", ratio1) + " Put:" + fmt.Sprintf("%.1f", ratio2) + "\r\n"
						email := Email{
							senderId: EMAIL_SENDER,
							toIds:    []string{EMAIL_RECEIVER},
							subject:  "SPY Volume (Call/Put:" + fmt.Sprintf("%.1f", ratio1) + " Put/Call:" + fmt.Sprintf("%.1f", ratio2) + ") Call: " + strconv.Itoa(totalCallVolumeInMinute) + " Put: " + strconv.Itoa(totalPutVolumeInMinute),
							body:     "High SPY Volume Alert (Call/Put:" + fmt.Sprintf("%.1f", ratio1) + " Put/Call:" + fmt.Sprintf("%.1f", ratio2) + ") @ " + time.Now().Format("15:04 Mon Jan _2 2006") + "\r\n" + volList,
							password: EMAIL_PASSWORD,
						}
						email.sendEmail()
					}
					averageVolumeArray = append(averageVolumeArray, totalVolumeInMinute)
					if len(averageVolumeArray) > MAX_LENGTH_OF_MINUTE_VOLUME_CHECKER {
						averageVolumeArray = averageVolumeArray[1:]
					}
				}
				fmt.Println(averageVolumeArray)
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":8888", nil))
}

var checker4WaitGroup sync.WaitGroup
var checker4Mut sync.Mutex

func (checker *Checker) runChecker4() {
	tblLogSystem := new(TblLogSystem)
	tblLogSystem.InsertLogSystem(LOGTYPE_CHECKER, "Checker 4 Started")
	fmt.Println("Checker Runs")
	symbolChan := checker.Producer()
	for i := 0; i < CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY; i++ {
		checker4WaitGroup.Add(1)
		go checker.Consumer(symbolChan)
	}
	checker4WaitGroup.Wait()
	//Send email
	emailBody := ""
	emailBody += fmt.Sprintf("%-10s%10s%10s\r\n", "Symbol", "Price", "Percent")
	for _, shd := range checker.shield {
		emailBody += fmt.Sprintf("%-10s", shd.Symbol)
		emailBody += fmt.Sprintf("%10.2f", shd.CurrentPrice)
		emailBody += fmt.Sprintf("%10.2f", shd.PercentChange)
		emailBody += "\r\n"
	}
	emailBody += fmt.Sprintf("***********************************************************************************\r\n")
	emailBody += fmt.Sprintf("\r\n%-10s%10s%10s%10s%10s%10s\r\n", "Symbol", "Price", "Trend", "MA60", "MA120", "Note")
	for _, dt := range checker.dragonTail {
		emailBody += fmt.Sprintf("%-10s", dt.Symbol)
		emailBody += fmt.Sprintf("%10.2f", dt.CurrentPrice)
		emailBody += fmt.Sprintf("%10s", dt.Trend)
		emailBody += fmt.Sprintf("%10.2f", dt.MA60)
		emailBody += fmt.Sprintf("%10.2f", dt.MA120)
		if _, ok := checker.addedDragonTail[dt.Symbol]; ok {
			emailBody += fmt.Sprintf("%10s", "Added")
		} else if _, ok := checker.removedDragonTail[dt.Symbol]; ok {
			emailBody += fmt.Sprintf("%10s", "Removed")
		}
		emailBody += "\r\n"
	}
	email := Email{
		senderId: EMAIL_SENDER,
		toIds:    []string{EMAIL_RECEIVER},
		subject:  "CardSharps Checker Report",
		body:     emailBody,
		password: EMAIL_PASSWORD,
	}
	email.sendEmail()
	return
}

func (checker *Checker) Producer() <-chan string {
	outChan := make(chan string, CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY)
	tblSymbol := new(TblSymbol)
	tblLogError := new(TblLogError)
	//Clear the dragon tail (added and removed) before running
	checker.addedDragonTail = make(map[string]*struct{})
	for key, _ := range checker.removedDragonTail {
		delete(checker.dragonTail, key)
	}
	checker.removedDragonTail = make(map[string]*struct{})
	checker.shield = make(map[string]*Shield)
	//GET SYMBOL LIST FROM db_symbol
	symbols, symbolSelectErr := tblSymbol.SelectSymbolByFilter()
	if symbolSelectErr != nil {
		tblLogError.InsertLogError(LOGTYPE_DB_SYMBOL, symbolSelectErr.Error())
	}
	// symbols = []string{"CBAY", "AAPL", "PDD", "NIO"}
	go func() {
		for _, symbol := range symbols {
			outChan <- symbol
		}
		defer close(outChan)
	}()
	return outChan
}

// func (checker *Checker) Consumer(inChan <-chan string) {
// 	yahooApiManager := new(YahooApiManager)
// 	for symbol := range inChan {
// 		fmt.Println("Running checker for " + symbol)
// 		stock, yahooApiErr := yahooApiManager.GetStockDataBySymbol(symbol)
// 		if yahooApiErr != nil {
// 			PrintMsgInConsole(MSGERROR, LOGTYPE_CHECKER+" "+symbol, yahooApiErr.Error())
// 			continue
// 		}
// 		hist, histErr := new(TblStockHist).SelectLastStockHist(symbol)
// 		if histErr != nil {
// 			PrintMsgInConsole(MSGERROR, LOGTYPE_DB_STOCK_DATA, histErr.Error())
// 			continue
// 		}
// 		if hist.isEmpty() {
// 			PrintMsgInConsole(MSGERROR, LOGTYPE_DB_STOCK_DATA, "No Stock Hist for "+symbol)
// 			continue
// 		}
// 		report, reportErr := new(TblStockReport).SelectLastStockReport(symbol)
// 		if reportErr != nil {
// 			PrintMsgInConsole(MSGERROR, LOGTYPE_DB_STOCK_REPORT, reportErr.Error())
// 			continue
// 		}
// 		if report.isEmpty() {
// 			PrintMsgInConsole(MSGERROR, LOGTYPE_DB_STOCK_REPORT, "No Stock Report for "+symbol)
// 			continue
// 		}
// 		var histAverage float32
// 		if hist.MarketOpen == 0 {
// 			histAverage = hist.MarketClose
// 		} else {
// 			histAverage = (2*hist.MarketClose + hist.MarketHigh + hist.MarketLow) / 4
// 		}

// 		//If satify all the assumptions, send email to users
// 		if stock.RegularMarketPrice > MinFloat32(report.MA60, report.MA120) && histAverage < MinFloat32(report.MA60, report.MA120) {
// 			email := Email{
// 				senderId: EMAIL_SENDER,
// 				toIds:    []string{EMAIL_RECEIVER},
// 				subject:  fmt.Sprintf("Uptrend %v - %0.2f - %0.2f - %0.2f", symbol, stock.RegularMarketPrice, report.MA60, report.MA120),
// 				body:     time.Now().Format("15:04 Mon Jan _2 2006") + "\r\n" + symbol + "\r\n" + "Price: " + fmt.Sprintf("%0.2f", stock.RegularMarketPrice) + "\r\n" + "MA60: " + fmt.Sprintf("%0.2f", report.MA60) + "\r\n" + "MA120: " + fmt.Sprintf("%0.2f", report.MA120) + "\r\n",
// 				password: EMAIL_PASSWORD,
// 			}
// 			email.sendEmail()
// 		} else if stock.RegularMarketPrice < MaxFloat32(report.MA60, report.MA120) && histAverage > MaxFloat32(report.MA60, report.MA120) {
// 			email := Email{
// 				senderId: EMAIL_SENDER,
// 				toIds:    []string{EMAIL_RECEIVER},
// 				subject:  fmt.Sprintf("Downtrend %v - %0.2f - %0.2f - %0.2f", symbol, stock.RegularMarketPrice, report.MA60, report.MA120),
// 				body:     time.Now().Format("15:04 Mon Jan _2 2006") + "\r\n" + symbol + "\r\n" + "Price: " + fmt.Sprintf("%0.2f", stock.RegularMarketPrice) + "\r\n" + "MA60: " + fmt.Sprintf("%0.2f", report.MA60) + "\r\n" + "MA120: " + fmt.Sprintf("%0.2f", report.MA120) + "\r\n",
// 				password: EMAIL_PASSWORD,
// 			}
// 			email.sendEmail()
// 		}
// 	}
// }

func (checker *Checker) Consumer(inChan <-chan string) {
	yahooApiManager := new(YahooApiManager)
	for symbol := range inChan {
		fmt.Println("Running checker for " + symbol)
		stock, yahooApiErr := yahooApiManager.GetStockDataBySymbol(symbol)
		if yahooApiErr != nil {
			PrintMsgInConsole(MSGERROR, LOGTYPE_CHECKER+" "+symbol, yahooApiErr.Error())
			continue
		}
		histList, histErr := new(TblStockHist).SelectLastStockHistByCountAndBeforeDate(symbol, DRAGON_TAIL_LENGTH, GetTimeInYYYYMMDD())
		if histErr != nil {
			PrintMsgInConsole(MSGERROR, LOGTYPE_DB_STOCK_DATA, histErr.Error())
			continue
		}
		if len(histList) < DRAGON_TAIL_LENGTH {
			PrintMsgInConsole(MSGERROR, LOGTYPE_DB_STOCK_DATA, "Not enough Stock Hist for "+symbol)
			continue
		}
		reportList, reportErr := new(TblStockReport).SelectLastStockReportByCountAndBeforeDate(symbol, DRAGON_TAIL_LENGTH, GetTimeInYYYYMMDD())
		if reportErr != nil {
			PrintMsgInConsole(MSGERROR, LOGTYPE_DB_STOCK_REPORT, reportErr.Error())
			continue
		}
		if len(reportList) != len(histList) || len(reportList) == 0 {
			PrintMsgInConsole(MSGERROR, LOGTYPE_DB_STOCK_REPORT, "Not enough Stock Report for "+symbol)
			continue
		}
		// var histAvgList []float32
		var numHistAboveMaxMA int
		var numHistBelowMinMA int
		var lastHistAvg float32
		for i, h := range histList {
			histAvg := h.MarketClose
			if h.MarketOpen != 0 {
				histAvg = (2*h.MarketClose + h.MarketHigh + h.MarketLow) / 4
			}
			if i == 0 {
				lastHistAvg = histAvg
			}
			if histAvg >= MaxFloat32(reportList[i].MA60, reportList[i].MA120) {
				numHistAboveMaxMA++
			} else if histAvg <= MinFloat32(reportList[i].MA60, reportList[i].MA120) {
				numHistBelowMinMA++
			}
		}

		//Shield
		if stock.RegularMarketChangePercent > SHIELD_HEIGHT || stock.RegularMarketChangePercent < -SHIELD_HEIGHT {
			checker.shield[symbol] = new(Shield)
			checker.shield[symbol].Symbol = symbol
			checker.shield[symbol].CurrentPrice = stock.RegularMarketPrice
			checker.shield[symbol].PercentChange = stock.RegularMarketChangePercent
		}

		//Dragon Tail
		if stock.RegularMarketPrice > MinFloat32(reportList[0].MA60, reportList[0].MA120) &&
			stock.RegularMarketPrice < MaxFloat32(reportList[0].MA60, reportList[0].MA120) &&
			lastHistAvg < MinFloat32(reportList[0].MA60, reportList[0].MA120) &&
			MaxFloat32(reportList[0].MA60, reportList[0].MA120) >= MinFloat32(reportList[0].MA60, reportList[0].MA120)*DRAGON_TAIL_MA_GAP &&
			numHistBelowMinMA > DRAGON_TAIL_LENGTH*DRAGON_TAIL_LIMIT {
			_, ok := checker.dragonTail[symbol]
			if ok {
				checker.dragonTail[symbol].CurrentPrice = stock.RegularMarketPrice
				checker.dragonTail[symbol].Trend = "UP"
			} else {
				checker.dragonTail[symbol] = new(DragonTail)
				checker.dragonTail[symbol].Symbol = symbol
				checker.dragonTail[symbol].CurrentPrice = stock.RegularMarketPrice
				checker.dragonTail[symbol].MA60 = reportList[0].MA60
				checker.dragonTail[symbol].MA120 = reportList[0].MA120
				checker.dragonTail[symbol].Trend = "UP"
				checker.addedDragonTail[symbol] = new(struct{})
			}
		} else if stock.RegularMarketPrice > MinFloat32(reportList[0].MA60, reportList[0].MA120) &&
			stock.RegularMarketPrice < MaxFloat32(reportList[0].MA60, reportList[0].MA120) &&
			lastHistAvg > MaxFloat32(reportList[0].MA60, reportList[0].MA120) &&
			MaxFloat32(reportList[0].MA60, reportList[0].MA120) >= MinFloat32(reportList[0].MA60, reportList[0].MA120)*DRAGON_TAIL_MA_GAP &&
			numHistAboveMaxMA > DRAGON_TAIL_LENGTH*DRAGON_TAIL_LIMIT {
			_, ok := checker.dragonTail[symbol]
			if ok {
				checker.dragonTail[symbol].CurrentPrice = stock.RegularMarketPrice
				checker.dragonTail[symbol].Trend = "DOWN"
			} else {
				checker.dragonTail[symbol] = new(DragonTail)
				checker.dragonTail[symbol].Symbol = symbol
				checker.dragonTail[symbol].CurrentPrice = stock.RegularMarketPrice
				checker.dragonTail[symbol].MA60 = reportList[0].MA60
				checker.dragonTail[symbol].MA120 = reportList[0].MA120
				checker.dragonTail[symbol].Trend = "DOWN"
				checker.addedDragonTail[symbol] = new(struct{})
			}
		} else {
			_, ok := checker.dragonTail[symbol]
			if ok {
				checker.removedDragonTail[symbol] = new(struct{})
			}
		}
	}
	checker4WaitGroup.Done()
}
