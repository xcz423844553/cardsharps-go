package main

import (
	"fmt"
	"strconv"
)

const (
	//ValidOptionExpirationFromNow determines the seconds from now which expiration date is valid to download to db_option
	ValidOptionExpirationFromNow = 31 * 24 * 3600
)

//Dealer downloads data and inserts into db at the end of each trading day
type Dealer2 struct {
}

func (dealer *Dealer2) isMarketOpen() bool {
	api := new(YahooApi)
	quote, apiErr := api.GetQuote("SPY")
	if apiErr != nil {
		return false
	}
	return quote.isMarketOpen()
}

func (dealer *Dealer2) isMarketPreOpenOrOpen() bool {
	api := new(YahooApi)
	quote, apiErr := api.GetQuote("SPY")
	if apiErr != nil {
		return false
	}
	return quote.isMarketPreOpen() || quote.isMarketOpen()
}

//DownloadAllOptionChainAndStock downloads the option chain and stock quote and insert all data into db_stock_data, db_option_data, db_stock_hist
func (dealer *Dealer2) DownloadAllOptionChainAndStock(symbol string) {
	daoOptionData := new(DaoOptionData)
	daoStockData := new(DaoStockData)
	daoStockHist := new(DaoStockHist)
	api := new(YahooApi)
	var options []ApiYahooOption
	var apiErr error
	_, stock, expDates, apiErr := api.GetOptionChainStockAndExpDate(symbol, 0)
	if apiErr != nil {
		fmt.Println(apiErr.Error())
		return
	}
	for _, expDate := range expDates {
		if expDate <= ConvertTime64InUnix(GetTimeInYYYYMMDD64())+ValidOptionExpirationFromNow {
			newOptions, optionErr := api.GetOptionChain(symbol, expDate)
			if optionErr != nil {
				fmt.Println("Error encounted while getting option chain data for " + symbol + " on date " + strconv.FormatInt(expDate, 10))
				fmt.Println(optionErr)
				continue
			}
			options = append(options, newOptions...)
		}
	}
	//store the option data into database
	var insertErr error
	for _, option := range options {
		if option.GetVolume() != 0 && option.GetOpenInterest() != 0 {
			insertErr = daoOptionData.InsertOrUpdateOptionData(&option)
			if insertErr != nil {
				fmt.Println("Error encounted while inserting option chain data for " + symbol + " on date " + strconv.FormatInt(option.GetExpiration(), 10))
				fmt.Println(insertErr)
				continue
			}
		}
	}
	//store the stock data into database
	insertErr = daoStockData.InsertOrUpdateStockData(&stock)
	if insertErr != nil {
		fmt.Println("Error encounted while inserting stock quote data for " + symbol)
		fmt.Println(insertErr)
	}
	insertErr = daoStockHist.InsertOrUpdateStockHist(&stock, GetTimeInYYYYMMDD64())
	if insertErr != nil {
		fmt.Println("Error encounted while inserting stock hist data for " + symbol)
		fmt.Println(insertErr)
	}
}
