package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type YahooApiManager struct {
}

func (manager YahooApiManager) GetOptionsAndStockDataBySymbol(symbol string) ([]YahooOption, YahooQuote, []int64, error) {
	_, stock, expDates, err1 := manager.GetOptionsAndStockDataBySymbolAndExpDate(symbol, -1)
	var options []YahooOption
	if err1 != nil {
		return options, stock, expDates, errors.New("Failed to get expiration dates for " + symbol)
	}
	for _, expDate := range expDates {
		newOption, _, _, err2 := manager.GetOptionsAndStockDataBySymbolAndExpDate(symbol, expDate)
		if err2 != nil {
			fmt.Println("Failed to get option and stock data for " + symbol + " on date " + strconv.FormatInt(expDate, 10))
			continue
		}
		options = append(options, newOption...)
	}
	return options, stock, expDates, nil
}

func (manager YahooApiManager) GetOptionsAndStockDataBySymbolAndExpDate(symbol string, expDate int64) ([]YahooOption, YahooQuote, []int64, error) {
	var options []YahooOption
	var stock YahooQuote
	var expDates []int64
	var optionErr error
	var stockErr error
	var expDateErr error
	url := URL_OPTION + symbol
	if expDate > 0 {
		url += "?date=" + strconv.FormatInt(expDate, 10)
	}
	resp, connError := http.Get(url)
	if connError != nil {
		fmt.Println(connError)
		return options, stock, expDates, connError
	}
	defer resp.Body.Close()
	body, parseError := ioutil.ReadAll(resp.Body)
	if parseError != nil {
		fmt.Println(parseError)
		return options, stock, expDates, parseError
	}
	var y YahooResponse
	if jsonError := json.Unmarshal(body, &y); jsonError != nil {
		fmt.Println(jsonError)
		return options, stock, expDates, parseError
	}
	options, optionErr = y.GetOptions()
	stock, stockErr = y.GetQuote()
	expDates, expDateErr = y.GetExpirationDates()
	if optionErr != nil {
		return options, stock, expDates, optionErr
	}
	if stockErr != nil {
		return options, stock, expDates, stockErr
	}
	if expDateErr != nil {
		return options, stock, expDates, expDateErr
	}
	return options, stock, expDates, nil
}
