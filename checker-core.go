package main

import (
	"errors"
	"fmt"
	"log"
	"math"
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
	bollKcChance      map[string]*BollKcChance
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

//BollKcChance is the stock recommended using Bollinger Band and Keltner Channel Strategy
type BollKcChance struct {
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

func (checker *Checker) InitChecker() *Checker {
	checker.dragonTail = make(map[string]*DragonTail)
	checker.addedDragonTail = make(map[string]*struct{})
	checker.removedDragonTail = make(map[string]*struct{})
	checker.shield = make(map[string]*Shield)
	checker.bollKcChance = make(map[string]*BollKcChance)
	return checker
}

func (checker Checker) runChecker() {
	yahooApiManager := new(YahooAPIManager)
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
	yahooApiManager := new(YahooAPIManager)
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
	yahooApiManager := new(YahooAPIManager)
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
	yahooApiManager := new(YahooAPIManager)
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
	symbols, symbolSelectErr := tblSymbol.SelectAllSymbol()
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

/// <summary>
/// This function runs the checking algorithm and filter out the stocks with dragon tail pattern.
/// </summary>
/// <param name="inChan"> Channel sending out the symbol</param>
/// <returns> no return </returns>
func (checker *Checker) Consumer(inChan <-chan string) {
	yahooApiManager := new(YahooAPIManager)
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
			MaxFloat32(reportList[0].MA60, reportList[0].MA120) >= MinFloat32(reportList[0].MA60, reportList[0].MA120)*(1+DRAGON_TAIL_MA_GAP) &&
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
			MaxFloat32(reportList[0].MA60, reportList[0].MA120) >= MinFloat32(reportList[0].MA60, reportList[0].MA120)*(1+DRAGON_TAIL_MA_GAP) &&
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

func (checker *Checker) Producer5() <-chan string {
	outChan := make(chan string, CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY)
	tblSymbol := new(TblSymbol)
	tblLogError := new(TblLogError)
	//Clear the dragon tail (added and removed) before running
	// checker.addedDragonTail = make(map[string]*struct{})
	// for key, _ := range checker.removedDragonTail {
	// 	delete(checker.dragonTail, key)
	// }
	// checker.removedDragonTail = make(map[string]*struct{})
	// checker.shield = make(map[string]*Shield)
	//GET SYMBOL LIST FROM db_symbol
	// symbols, symbolSelectErr := tblSymbol.SelectSymbolByTrader(TRADER_DOW)
	symbols, symbolSelectErr := tblSymbol.SelectAllSymbol()
	if symbolSelectErr != nil {
		tblLogError.InsertLogError(LOGTYPE_DB_SYMBOL, symbolSelectErr.Error())
	}
	// symbols = []string{"CBAY", "AAPL", "PDD", "NIO"}
	// symbols = []string{"AAPL"}
	go func() {
		for _, symbol := range symbols {
			outChan <- symbol
		}
		defer close(outChan)
	}()
	return outChan
}

var checker5WaitGroup sync.WaitGroup
var checker5Mut sync.Mutex

func (checker *Checker) runChecker5() {
	tblLogSystem := new(TblLogSystem)
	tblLogSystem.InsertLogSystem(LOGTYPE_CHECKER, "Checker 5 Started")
	fmt.Println("Checker Runs")
	symbolChan := checker.Producer5()
	for i := 0; i < CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY; i++ {
		checker5WaitGroup.Add(1)
		go checker.Consumer5(symbolChan)
	}
	checker5WaitGroup.Wait()
	//Send email
	emailBody := ""
	emailBody += fmt.Sprintf("%10s%10s%10s%10s%10s%10s%10s%10s\r\n", "STATUS", "Trend", "Symbol", "Price", "BMID", "BUP", "KMID", "KUP")
	for _, chance := range checker.bollKcChance {
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
		subject:  "CardSharps Checker Report",
		body:     emailBody,
		password: EMAIL_PASSWORD,
	}
	email.sendEmail()
	return
}

func (checker *Checker) Consumer5(inChan <-chan string) {
	PARAM_MARKET_OPEN := false
	PARAM_CROSS_RANGE := 5
	PARAM_LOOKBACK_RANGE := 20
	PARAM_MACD_EMA_RANGE_SHORT := 12
	PARAM_MACD_EMA_RANGE_LONG := 26
	PARAM_MACD_RANGE := 9
	yahooApiManager := new(YahooAPIManager)

	for symbol := range inChan {
		fmt.Println("Running checker for " + symbol)

		//Get real time stock data
		stock, yahooApiErr := yahooApiManager.GetStockDataBySymbol(symbol)
		if yahooApiErr != nil {
			PrintMsgInConsole(MSGERROR, LOGTYPE_CHECKER+" "+symbol, yahooApiErr.Error())
			continue
		}

		//Get historical stock data (reversed, 0 is previous day)
		histList, histErr := new(TblStockHist).SelectLastStockHistByCountAndBeforeDate(symbol, PARAM_CROSS_RANGE+PARAM_LOOKBACK_RANGE*2+PARAM_MACD_EMA_RANGE_LONG*2+PARAM_MACD_RANGE*2, GetTimeInYYYYMMDD()+1)
		if histErr != nil {
			PrintMsgInConsole(MSGERROR, LOGTYPE_DB_STOCK_DATA, histErr.Error())
			continue
		}
		if len(histList) < PARAM_CROSS_RANGE+PARAM_LOOKBACK_RANGE*2 {
			PrintMsgInConsole(MSGERROR, LOGTYPE_DB_STOCK_DATA, "Not enough Stock Hist for "+symbol)
			continue
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
		if PARAM_MARKET_OPEN {
			histList = append([]RowStockHist{stockHist}, histList...)
		}

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
		var macd float64
		var macdArray []float32
		var macdErr error
		for i := 0; i < PARAM_CROSS_RANGE; i++ {
			bollMid, bollUpper, bollLower, kcMid, kcUpper, kcLower := checker.CalcBollKcStat(histList, i, PARAM_CROSS_RANGE, PARAM_LOOKBACK_RANGE)
			bollMidArray = append(bollMidArray, bollMid)
			bollUpperArray = append(bollUpperArray, bollUpper)
			bollLowerArray = append(bollLowerArray, bollLower)
			kcMidArray = append(kcMidArray, kcMid)
			kcUpperArray = append(kcUpperArray, kcUpper)
			kcLowerArray = append(kcLowerArray, kcLower)
			macd, macdErr = checker.CalcMACDStat(closePriceList[i:], PARAM_MACD_EMA_RANGE_SHORT, PARAM_MACD_EMA_RANGE_LONG, PARAM_MACD_RANGE)
			if macdErr != nil {
				PrintMsgInConsole(MSGERROR, LOGTYPE_CHECKER, macdErr.Error())
				break
			}
			macdArray = append(macdArray, float32(macd))
		}
		if macdErr != nil {
			continue
		}

		bollKcChance := BollKcChance{
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

		if bollUpperArray[0] > kcUpperArray[0] && bollUpperArray[PARAM_CROSS_RANGE-1] < kcUpperArray[PARAM_CROSS_RANGE-1] && stock.RegularMarketPrice > bollUpperArray[0] {
			if _, ok := checker.bollKcChance[stock.Symbol]; ok {
				bollKcChance.Status = "KEEP"
				bollKcChance.Trend = "UP"
			} else {
				bollKcChance.Status = "NEW"
				bollKcChance.Trend = "UP"
			}
			checker.bollKcChance[stock.Symbol] = &bollKcChance
			fmt.Printf("Long %s %.2f\r\n", stock.Symbol, stock.RegularMarketPrice)
		} else if bollLowerArray[0] < kcLowerArray[0] && bollLowerArray[PARAM_CROSS_RANGE-1] > kcLowerArray[PARAM_CROSS_RANGE-1] && stock.RegularMarketPrice < bollLowerArray[0] {
			if _, ok := checker.bollKcChance[stock.Symbol]; ok {
				bollKcChance.Status = "KEEP"
				bollKcChance.Trend = "DOWN"
			} else {
				bollKcChance.Status = "NEW"
				bollKcChance.Trend = "DOWN"
			}
			checker.bollKcChance[stock.Symbol] = &bollKcChance
			fmt.Printf("Short %s %.2f\r\n", stock.Symbol, stock.RegularMarketPrice)
		} else if (kcUpperArray[0]-bollUpperArray[0]+bollLowerArray[0]-kcLowerArray[0])/(bollUpperArray[0]-bollLowerArray[0]) > 2 {
			if _, ok := checker.bollKcChance[stock.Symbol]; ok {
				bollKcChance.Status = "KEEP"
				bollKcChance.Trend = "SQUEEZE"
			} else {
				bollKcChance.Status = "NEW"
				bollKcChance.Trend = "SQUEEZE"
			}
			checker.bollKcChance[stock.Symbol] = &bollKcChance
			fmt.Printf("Squeeze %s %.2f\r\n", stock.Symbol, stock.RegularMarketPrice)
			// } else if bollUpperArray[0] < kcUpperArray[0] && stock.RegularMarketPrice > kcUpperArray[0] && stock.RegularMarketDayLow > bollLowerArray[0] {
			// 	if _, ok := checker.bollKcChance[stock.Symbol]; ok {
			// 		bollKcChance.Status = "KEEP"
			// 		bollKcChance.Trend = "TEST"
			// 	} else {
			// 		bollKcChance.Status = "NEW"
			// 		bollKcChance.Trend = "TEST"
			// 	}
			// 	checker.bollKcChance[stock.Symbol] = &bollKcChance
			// 	fmt.Printf("Squeeze %s %.2f\r\n", stock.Symbol, stock.RegularMarketPrice)
		} else if macdArray[0] > 0 && macdArray[1] < 0 {
			if _, ok := checker.bollKcChance[stock.Symbol]; ok {
				bollKcChance.Status = "KEEP"
				bollKcChance.Trend = "TEST"
			} else {
				bollKcChance.Status = "NEW"
				bollKcChance.Trend = "TEST"
			}
			checker.bollKcChance[stock.Symbol] = &bollKcChance
			fmt.Printf("Squeeze %s %.2f\r\n", stock.Symbol, stock.RegularMarketPrice)
		} else {
			if _, ok := checker.bollKcChance[stock.Symbol]; ok {
				bollKcChance.Status = "REMOVE"
				checker.bollKcChance[stock.Symbol] = &bollKcChance
			}
		}
	}
	checker5WaitGroup.Done()
}

func (checker *Checker) CalcMACDStat(data []float64, emaRangeShort int, emaRangeLong int, macdRange int) (float64, error) {
	//Minimum length of data is PARAM_MACD_RANGE * 2+PARAM_EMA_LONG_RANGE * 2
	PARAM_EMA_SHORT_RANGE := emaRangeShort
	PARAM_EMA_LONG_RANGE := emaRangeLong
	PARAM_MACD_RANGE := macdRange

	if len(data) < PARAM_MACD_RANGE*2+PARAM_EMA_LONG_RANGE*2 {
		return 0, errors.New("Not enough data to calculate MACD")
	}

	var emaShort float64 //Short Exponential Moving Average
	var emaLong float64  //Long Exponential Moving Average
	var diff []float64
	var dea float64
	var macd float64

	var emaShortErr error
	var emaLongErr error
	var deaErr error

	for i := 0; i < PARAM_MACD_RANGE*2; i++ {

		// fmt.Printf("MACD #: %d\r\n", i)
		emaShort, emaShortErr = checker.CalcEMAStat(data[i:], PARAM_EMA_SHORT_RANGE)
		if emaShortErr != nil {
			return macd, emaShortErr
		}
		// fmt.Printf("SHORT: %.2f\r\n", emaShort)
		emaLong, emaLongErr = checker.CalcEMAStat(data[i:], PARAM_EMA_LONG_RANGE)
		if emaLongErr != nil {
			return macd, emaLongErr
		}
		// fmt.Printf("LONG: %.2f\r\n", emaLong)
		diff = append(diff, emaShort-emaLong)
	}
	dea, deaErr = checker.CalcEMAStat(diff, PARAM_MACD_RANGE)
	if deaErr != nil {
		return macd, deaErr
	}
	macd = diff[0] - dea
	// fmt.Printf("DIFF value: %.2f\r\n", diff[0])
	// fmt.Printf("DEA value: %.2f\r\n", dea)
	return macd, nil
}

func (checker *Checker) CalcEMAStat(data []float64, emaRange int) (float64, error) {
	var ema float64
	if len(data) < emaRange || emaRange <= 0 {
		return ema, errors.New("Not enough data to calculate EMA")
	}
	PARAM_EMA_RANGE := emaRange
	EMA_MULTIPLIER := 2 / float64(PARAM_EMA_RANGE+1)

	var ma float64 //Moving Average

	//Calculate MA(EMA_RANGE)
	for i := PARAM_EMA_RANGE; i < PARAM_EMA_RANGE*2; i++ {
		ma += float64(data[i])
	}
	ma /= float64(PARAM_EMA_RANGE)

	//Calculate EMA(EMA_RANGE)
	var emaPrev float64
	emaPrev = ma
	for i := PARAM_EMA_RANGE - 1; i >= 0; i-- {
		ema = float64(data[i])*EMA_MULTIPLIER + emaPrev*(1-EMA_MULTIPLIER)
		emaPrev = ema
	}
	return ema, nil
}

func (checker *Checker) CalcBollKcStat(histList []RowStockHist, targetIndex int, crossRange int, lookbackRange int) (float32, float32, float32, float32, float32, float32) {
	PARAM_MA_RANGE := lookbackRange
	PARAM_EMA_PREV_RANGE := lookbackRange
	PARAM_ATR_RANGE := lookbackRange / 2
	PARAM_NUM_SIGMA := 2
	EMA_MULTIPLIER := 2 / float64(PARAM_MA_RANGE+1)

	var ma float64        //Moving Average
	var sd float64        //Standard Deviation
	var ema float64       //Exponential Moving Average
	var atr float64       //Average True Range
	var bollMid float64   //Mid Line of Bollinger Band
	var bollUpper float64 //Upper Bound of Bollinger Band
	var bollLower float64 //Lower Bound of Bollinger Band
	var kcMid float64     //Mid Line of Keltner Channel
	var kcUpper float64   //Upper Bound of Keltner CHannel
	var kcLower float64   //Lower Bound of Keltner Channel

	//Calculate MA(MA_RANGE), including the real time stock price
	for i := targetIndex; i < targetIndex+PARAM_MA_RANGE; i++ {
		ma += float64((histList[i].MarketHigh + histList[i].MarketLow + 2*histList[i].MarketClose) / 4)
	}
	ma /= float64(PARAM_MA_RANGE)

	//Calculate SIGMA(MA_RANGE), including the real time stock price
	for i := targetIndex; i < targetIndex+PARAM_MA_RANGE; i++ {
		sdElm := float64((histList[i].MarketHigh+histList[i].MarketLow+2*histList[i].MarketClose)/4) - ma
		sd += math.Pow(sdElm, 2)
	}
	sd = math.Sqrt(sd / float64(PARAM_MA_RANGE))

	//Calculate EMA(MA_RANGE, EMA_PREV_RANGE), including the real time stock price
	var emaPrev float64
	for i := targetIndex + PARAM_MA_RANGE; i < targetIndex+PARAM_MA_RANGE+PARAM_EMA_PREV_RANGE; i++ {
		emaPrev += float64((histList[i].MarketHigh + histList[i].MarketLow + 2*histList[i].MarketClose) / 4)
	}
	emaPrev /= float64(PARAM_EMA_PREV_RANGE)
	for i := targetIndex + PARAM_MA_RANGE - 1; i >= targetIndex; i-- {
		ema = float64((histList[i].MarketHigh+histList[i].MarketLow+2*histList[i].MarketClose)/4)*EMA_MULTIPLIER + emaPrev*(1-EMA_MULTIPLIER)
		emaPrev = ema
	}

	//Calculate ATR(ATR_RANGE), including the real time stock price
	for i := targetIndex; i < targetIndex+PARAM_ATR_RANGE; i++ {
		trElm1 := float64(histList[i].MarketHigh - histList[i].MarketLow)
		trElm2 := math.Abs(float64(histList[i].MarketHigh - histList[i+1].MarketClose))
		trElm3 := math.Abs(float64(histList[i].MarketLow - histList[i+1].MarketClose))
		tr := math.Max(math.Max(trElm1, trElm2), trElm3)
		atr += tr
	}
	atr = atr / float64(PARAM_ATR_RANGE)

	//Calculate Parameters of Bollinger Band
	bollMid = ma
	bollUpper = bollMid + float64(PARAM_NUM_SIGMA)*sd
	bollLower = bollMid - float64(PARAM_NUM_SIGMA)*sd
	//Calculate Parameters of Keltner Channel
	kcMid = ema
	kcUpper = kcMid + float64(PARAM_NUM_SIGMA)*atr
	kcLower = kcMid - float64(PARAM_NUM_SIGMA)*atr
	// fmt.Printf("Index:%v Symbol:%s CLOSE: %.2f MA:%.2f SD:%.2f EMA:%.2f ATR:%.2f BMID:%.2f BUP:%.2f BLO:%.2f KMID:%.2f KUP:%.2f KLO:%.2f\r\n", targetIndex, histList[targetIndex].Symbol, histList[targetIndex].MarketClose, ma, sd, ema, atr, bollMid, bollUpper, bollLower, kcMid, kcUpper, kcLower)
	return float32(bollMid), float32(bollUpper), float32(bollLower), float32(kcMid), float32(kcUpper), float32(kcLower)
}
