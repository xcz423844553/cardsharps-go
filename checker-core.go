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
			"CSCO", "PFE", "DIS", "VZ", "T", "CVX", "UNH", "HD", "KO", "MRK", "INTC", "WFC", "ORCL", "CMCSA", "PEP", "NFLX", "MCD", "C"}
		volumes := make([]int, len(symbols))
		symbolVolumes := make([]int, len(symbols))
		symbolPrices := make([]float32, len(symbols))
		totalMoneyInMinute := 0
		averageMoneyArray := []int{}
		for ; true; <-ticker.C {
			if !isMarketOpen() && !BYPASS_MARKET_STATUS {
				fmt.Println("Market not open yet")
			} else {
				fmt.Println("Tick at ", time.Now())
				totalMoneyInMinute = 0
				for index, symbol := range symbols {
					checkerWaitGroup.Add(1)
					go func(volumes []int, symbolVolumes []int, symbolPrices []float32, totalMoneyInMinute *int, index int, symbol string) {
						_, quote, _, yahooErr := yahooApiManager.GetOptionsAndStockDataBySymbol(symbol)
						if yahooErr != nil {
							tblLogError.InsertLogError(LOGTYPE_CHECKER, yahooErr.Error())
						}
						fmt.Println("(" + strconv.Itoa(index) + ")" + symbol + " " + strconv.Itoa(quote.RegularMarketVolume))
						checkerMut.Lock()
						*totalMoneyInMinute += int(float32(quote.RegularMarketVolume-volumes[index]) * quote.RegularMarketPrice)
						checkerMut.Unlock()
						volumes[index] = quote.RegularMarketVolume
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
					averageMoney := AverageInt(averageMoneyArray)
					if totalMoneyInMinute >= MULTI_THRESHOLD_MONEY_CHECKER*averageMoney {
						volList := ""
						ratio := fmt.Sprintf("%.1f", float32(totalMoneyInMinute)/float32(averageMoney))
						for index, vol := range averageMoneyArray {
							volList += "(" + strconv.Itoa(index) + ") " + strconv.Itoa(vol) + "\r\n"
						}
						volList += "(0) " + strconv.Itoa(totalMoneyInMinute) + " @ " + ratio + "\r\n"
						email := Email{
							senderId: EMAIL_SENDER,
							toIds:    []string{EMAIL_RECEIVER},
							subject:  "CardSharps Notification",
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
		averageVolumeArray := []int{}
		for ; true; <-ticker.C {
			if !isMarketOpen() && !BYPASS_MARKET_STATUS {
				fmt.Println("Market not open yet")
			} else {
				fmt.Println("Tick at ", time.Now())
				totalVolumeInMinute = 0
				for index, expDate := range expDates {
					if expDate > (time.Now().Unix() + UNIX_TWO_WEEK) {
						continue
					}
					spyCheckerWaitGroup.Add(1)
					go func(volumeMap map[string]int, totalVolumeInMinute *int, index int, symbol string, expDate int64) {
						options, _, _, yahooErr := yahooApiManager.GetOptionsAndStockDataBySymbolAndExpDate(symbol, expDate)
						if yahooErr != nil {
							tblLogError.InsertLogError(LOGTYPE_SPY_CHECKER, yahooErr.Error())
						}
						spyCheckerMut.Lock()
						for _, option := range options {
							vol, _ := volumeMap[option.ContractSymbol]
							volumeMap[option.ContractSymbol] = option.Volume
							*totalVolumeInMinute += option.Volume - vol
							fmt.Println("(" + strconv.Itoa(index) + ")" + option.ContractSymbol + " " + strconv.Itoa(option.Volume))
						}
						spyCheckerMut.Unlock()
						spyCheckerWaitGroup.Done()
					}(volumeMap, &totalVolumeInMinute, index, symbol, expDate)
				}
				spyCheckerWaitGroup.Wait()
				if len(averageVolumeArray) < MIN_LENGTH_OF_MINUTE_MONEY_CHECKER {
					averageVolumeArray = append(averageVolumeArray, totalVolumeInMinute)
				} else {
					averageVolume := AverageInt(averageVolumeArray)
					if totalVolumeInMinute >= MULTI_THRESHOLD_VOLUME_CHECKER*averageVolume {
						volList := ""
						ratio := fmt.Sprintf("%.1f", float32(totalVolumeInMinute)/float32(averageVolume))
						for index, vol := range averageVolumeArray {
							volList += "(" + strconv.Itoa(index) + ") " + strconv.Itoa(vol) + "\r\n"
						}
						volList += "(0) " + strconv.Itoa(totalVolumeInMinute) + " @ " + ratio + "\r\n"
						email := Email{
							senderId: EMAIL_SENDER,
							toIds:    []string{EMAIL_RECEIVER},
							subject:  "CardSharps Notification",
							body:     "High SPY Volume Alert (" + ratio + ") @ " + time.Now().Format("15:04 Mon Jan _2 2006") + "\r\n" + volList,
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
