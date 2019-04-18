package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	//"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
)

type OptionFilter struct {
	MaxOptionPercent  float32 `json:"maxOptionPercent"`
	MinOptionPercent  float32 `json:"minOptionPercent"`
	MaxOpenInterest   int64   `json:"maxOpenInterest"`
	MinOpenInterest   int64   `json:"minOpenInterest"`
	MaxVolume         int64   `json:"maxVolume"`
	MinVolume         int64   `json:"minVolume"`
	MaxExpirationDate int64   `json:"maxExpirationDate"`
	MinExpirationDate int64   `json:"minExpirationDate"`
}

func GetOptionsAndStockDataBySymbol(symbol string) ([]YahooOption, YahooQuote, error) {
	return GetOptionsAndStockDataBySymbolAndExpDate(symbol, -1)
}

func GetOptionsAndStockDataBySymbolAndExpDate(symbol string, expDate int64) ([]YahooOption, YahooQuote, error) {
	url := URL_OPTION + symbol
	if expDate > 0 {
		url += "?dates=" + strconv.FormatInt(expDate, 10)
	}
	fmt.Println(url)
	resp, connError := http.Get(url)
	if connError != nil {
		log.Fatal(connError)
	}
	defer resp.Body.Close()
	body, parseError := ioutil.ReadAll(resp.Body)
	if parseError != nil {
		log.Fatal(parseError)
	}
	var y YahooResponse
	if jsonError := json.Unmarshal(body, &y); jsonError != nil {
		log.Fatal(jsonError)
	}
	option, optionErr := y.GetOption()
	stock, stockErr := y.GetQuote()
	if optionErr != nil {
		return option, stock, optionErr
	}
	if stockErr != nil {
		return option, stock, stockErr
	}
	return option, stock, nil
}

func IsInOptionFilter

func NewOptionFilter(maxOptionPercent float32, minOptionPercent float32,
	maxOpenInterest int64, minOpenInterest int64, maxVolume int64, minVolume int64,
	maxExpirationDate int64, minExpirationDate int64) OptionFilter {
	f := new(OptionFilter)
	if maxOptionPercent == 0.0 {
		f.MaxOptionPercent = 1.0
	} else {
		f.MaxOptionPercent = maxOptionPercent
	}
	if minOptionPercent == 0.0 {
		f.MinOptionPercent = -1.0
	} else {
		f.MinOptionPercent = minOptionPercent
	}
	if maxOpenInterest == 0 {
		f.MaxOpenInterest = math.MaxInt64
	} else {
		f.MaxOpenInterest = maxOpenInterest
	}
	if minOpenInterest == 0 {
		f.MinOpenInterest = math.MinInt64
	} else {
		f.MinOpenInterest = minOpenInterest
	}
	if maxVolume == 0 {
		f.MaxVolume = math.MaxInt64
	} else {
		f.MaxVolume = maxVolume
	}
	if minVolume == 0 {
		f.MinVolume = math.MinInt64
	} else {
		f.MinVolume = minVolume
	}
	if maxExpirationDate == 0 {
		f.MaxExpirationDate = math.MaxInt64
	} else {
		f.MaxExpirationDate = maxExpirationDate
	}
	if minExpirationDate == 0 {
		f.MinExpirationDate = math.MinInt64
	} else {
		f.MinExpirationDate = minExpirationDate
	}
	return f
}
