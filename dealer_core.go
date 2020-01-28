package main

import (
	"fmt"
	"net/http"
	"strconv"

	tda "github.com/xcz423844553/td_ameritrade_client_golang"
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
func (dealer *Dealer2) DownloadAllOptionChainAndStock(client *http.Client, symbol string) {
	daoOptionData := new(DaoOptionData)
	// daoStockData := new(DaoStockData)
	// daoStockHist := new(DaoStockHist)

	optionChain := tda.GetOptionChain(client, clientID, symbol)
	//store the option data into database
	for _, opt := range optionChain.CallMap {
		option := new(ApiTdaOptionWrapper)
		option.BuildWrapper(opt)
		insertErr := daoOptionData.InsertOrUpdateOptionData(option)
		if insertErr != nil {
			fmt.Println("Error encounted while inserting option chain data for " + symbol + " on date " + strconv.FormatInt(option.GetExpiration(), 10))
			fmt.Println(insertErr)
			continue
		}
		if option.Volume > option.OpenInterest*10 && option.Volume > 1200 && float32(option.Volume)*option.LastPrice > 10000 {
			fmt.Printf("%s has %v volume vs. %v open interest.\r\n\r\n", option.ContractSymbol, option.Volume, option.OpenInterest)
		}
	}
	for _, opt := range optionChain.PutMap {
		option := new(ApiTdaOptionWrapper)
		option.BuildWrapper(opt)
		insertErr := daoOptionData.InsertOrUpdateOptionData(option)
		if insertErr != nil {
			fmt.Println("Error encounted while inserting option chain data for " + symbol + " on date " + strconv.FormatInt(option.GetExpiration(), 10))
			fmt.Println(insertErr)
			continue
		}
		if option.Volume > option.OpenInterest*10 && option.Volume > 1200 && float32(option.Volume)*option.LastPrice > 10000 {
			fmt.Printf("%s has %v volume vs. %v open interest.\r\n\r\n", option.ContractSymbol, option.Volume, option.OpenInterest)
		}
	}

	//store the stock data into database
	// insertErr = daoStockData.InsertOrUpdateStockData(&stock)
	// if insertErr != nil {
	// 	fmt.Println("Error encounted while inserting stock quote data for " + symbol)
	// 	fmt.Println(insertErr)
	// }
	// insertErr = daoStockHist.InsertOrUpdateStockHist(&stock, GetTimeInYYYYMMDD64())
	// if insertErr != nil {
	// 	fmt.Println("Error encounted while inserting stock hist data for " + symbol)
	// 	fmt.Println(insertErr)
	// }
}

//DownloadAllOptionChainAndStock downloads the option chain and stock quote and insert all data into db_stock_data, db_option_data, db_stock_hist
// func (dealer *Dealer2) DownloadAllOptionChainAndStock(symbol string) {
// 	daoOptionData := new(DaoOptionData)
// 	daoStockData := new(DaoStockData)
// 	daoStockHist := new(DaoStockHist)
// 	api := new(YahooApi)
// 	var options []ApiYahooOption
// 	var apiErr error
// 	_, stock, expDates, apiErr := api.GetOptionChainStockAndExpDate(symbol, 0)
// 	if apiErr != nil {
// 		fmt.Println(apiErr.Error())
// 		return
// 	}
// 	for _, expDate := range expDates {
// 		if expDate <= ConvertTime64InUnix(GetTimeInYYYYMMDD64())+ValidOptionExpirationFromNow {
// 			newOptions, optionErr := api.GetOptionChain(symbol, expDate)
// 			if optionErr != nil {
// 				fmt.Println("Error encounted while getting option chain data for " + symbol + " on date " + strconv.FormatInt(expDate, 10))
// 				fmt.Println(optionErr)
// 				continue
// 			}
// 			options = append(options, newOptions...)
// 		}
// 	}
// 	//store the option data into database
// 	var insertErr error
// 	for _, option := range options {
// 		if option.GetVolume() != 0 {
// 			insertErr = daoOptionData.InsertOrUpdateOptionData(&option)
// 			if insertErr != nil {
// 				fmt.Println("Error encounted while inserting option chain data for " + symbol + " on date " + strconv.FormatInt(option.GetExpiration(), 10))
// 				fmt.Println(insertErr)
// 				continue
// 			}

// 			//comment
// 			if (option.GetStrike() > stock.GetMarketClose()*1.05 && option.GetOptionType() == "C") || (option.GetStrike() < stock.GetMarketClose()*0.95 && option.GetOptionType() == "P") {
// 				if (option.GetBid()+option.GetAsk())/2*float32(option.GetVolume()) > 2000 && option.GetVolume() > 2*option.GetOpenInterest() {
// 					if option.GetExpiration() > 20200126 {
// 						fmt.Printf("TEST: %s @ $%v has volume %v\r\n", symbol, option.GetStrike(), option.GetVolume())
// 					}
// 				}
// 			}
// 			//comment

// 		}
// 	}
// 	//store the stock data into database
// 	insertErr = daoStockData.InsertOrUpdateStockData(&stock)
// 	if insertErr != nil {
// 		fmt.Println("Error encounted while inserting stock quote data for " + symbol)
// 		fmt.Println(insertErr)
// 	}
// 	insertErr = daoStockHist.InsertOrUpdateStockHist(&stock, GetTimeInYYYYMMDD64())
// 	if insertErr != nil {
// 		fmt.Println("Error encounted while inserting stock hist data for " + symbol)
// 		fmt.Println(insertErr)
// 	}
// }
