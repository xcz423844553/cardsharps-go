package main

import (
	"errors"
	"fmt"
)

type YahooResponse struct {
	OptionChain struct {
		Results []struct {
			UnderlyingSymbol string             `json:"underlyingSymbol"`
			ExpirationDates  []int64            `json:"expirationDates"`
			Strikes          []float32          `json:"strikes"`
			HasMiniOptions   bool               `json:"hasMiniOptions"`
			Quote            YahooQuote         `json:"quote"`
			OptionsArray     []YahooOptionArray `json:"options"`
		} `json:"result"`
		Error string `json:"error"`
	} `json:"optionChain"`
}

type YahooQuote struct {
	Symbol                            string  `json:"symbol"`
	RegularMarketChange               float32 `json:"regularMarketChange"`
	RegularMarketOpen                 float32 `json:"regularMarketOpen"`
	RegularMarketDayHigh              float32 `json:"regularMarketDayHigh"`
	RegularMarketDayLow               float32 `json:"regularMarketDayLow"`
	RegularMarketVolume               int     `json:"regularMarketVolume"`
	RegularMarketChangePercent        float32 `json:"regularMarketChangePercent"`
	RegularMarketPreviousClose        float32 `json:"regularMarketPreviousClose"`
	RegularMarketPrice                float32 `json:"regularMarketPrice"`
	RegularMarketTime                 int64   `json:"regularMarketTime"`
	EarningsTimestamp                 int     `json:"earningsTimestamp"`
	FiftyDayAverage                   float32 `json:"fiftyDayAverage"`
	FiftyDayAverageChange             float32 `json:"fiftyDayAverageChange"`
	FiftyDayAverageChangePercent      float32 `json:"fiftyDayAverageChangePercent"`
	TwoHundredDayAverage              float32 `json:"twoHundredDayAverage"`
	TwoHundredDayAverageChange        float32 `json:"twoHundredDayAverageChange"`
	TwoHundredDayAverageChangePercent float32 `json:"twoHundredDayAverageChangePercent"`
	Tradeable                         bool    `json:"tradeable"`
	MarketState                       string  `json:"marketState"` //CLOSED
	PostMarketChangePercent           float32 `json:"postMarketChangePercent"`
	PostMarketTime                    int64   `json:"postMarketTime"`
	PostMarketPrice                   float32 `json:"postMarketPrice"`
	PostMarketChange                  float32 `json:"postMarketChange"`
	Bid                               float32 `json:"bid"`
	Ask                               float32 `json:"ask"`
	BidSize                           int     `json:"bidSize"`
	AskSize                           int     `json:"askSize"`
	AverageDailyVolume3Month          int     `json:"averageDailyVolume3Month"`
	AverageDailyVolume10Day           int     `json:"averageDailyVolume10Day"`
	FiftyTwoWeekLowChange             float32 `json:"fiftyTwoWeekLowChange"`
	FiftyTwoWeekLowChangePercent      float32 `json:"fiftyTwoWeekLowChangePercent"`
	FiftyTwoWeekHighChange            float32 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekHighChangePercent     float32 `json:"fiftyTwoWeekHighChangePercent"`
	FiftyTwoWeekLow                   float32 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh                  float32 `json:"fiftyTwoWeekHigh"`
}

type YahooOptionArray struct {
	ExpirationDate int64         `json:"expirationDate"`
	Calls          []YahooOption `json:"calls"`
	Puts           []YahooOption `json:"puts"`
}

type YahooOption struct {
	ContractSymbol    string  `json:"contractSymbol"`
	Strike            float32 `json:"strike"`
	LastPrice         float32 `json:"lastPrice"`
	PriceChange       float32 `json:"change"`
	PercentChange     float32 `json:"percentChange"`
	Volume            int     `json:"volume"`
	OpenInterest      int     `json:"openInterest"`
	Bid               float32 `json:"bid"`
	Ask               float32 `json:"ask"`
	Expiration        int64   `json:"expiration"`
	ImpliedVolatility float32 `json:"impliedVolatility"`
	InTheMoney        bool    `json:"inTheMoney"`
}

func (resp YahooResponse) isEmptyResult() bool {
	resultsArray := resp.OptionChain.Results
	return len(resultsArray) == 0
}

func (resp YahooResponse) GetQuote() (YahooQuote, error) {
	var quote YahooQuote
	if resp.isEmptyResult() {
		return quote, errors.New("Quote response is empty")
	}
	quote = resp.OptionChain.Results[0].Quote
	return quote, nil
}

func (resp YahooResponse) GetOptions() ([]YahooOption, error) {
	var options []YahooOption
	if resp.isEmptyResult() {
		return options, errors.New("Option response is empty")
	}
	callArray := resp.OptionChain.Results[0].OptionsArray[0].Calls
	putArray := resp.OptionChain.Results[0].OptionsArray[0].Puts
	options = append(callArray, putArray...)
	return options, nil
}

func (resp YahooResponse) GetExpirationDates() ([]int64, error) {
	//only reads the closest exp date
	expDates := make([]int64, 4)
	if resp.isEmptyResult() {
		return expDates, errors.New("Expiration date response is empty")
	}
	expDates[0] = resp.OptionChain.Results[0].ExpirationDates[0]
	expDates[1] = resp.OptionChain.Results[0].ExpirationDates[1]
	expDates[2] = resp.OptionChain.Results[0].ExpirationDates[2]
	expDates[3] = resp.OptionChain.Results[0].ExpirationDates[3]
	return expDates, nil
}

func (quote YahooQuote) isMarketOpen() bool {
	fmt.Println(quote.MarketState)
	return quote.MarketState == "REGULAR"
}

// func (resp YahooResponse) GetOptionByDefaultFilter() ([]YahooOption, error) {
// 	f := NewOptionFilter(maxOptionPercent float32, minOptionPercent float32,
// 		maxOpenInterest int64, minOpenInterest int64, maxVolume int64, minVolume int64,
// 		maxExpirationDate int64, minExpirationDate int64)
// 	var options []YahooOption
// 	if resp.isEmptyResult() {
// 		return options, errors.New("Option response is empty")
// 	}
// 	callArray := resp.OptionChain.Results[0].OptionsArray[0].Calls
// 	putArray := resp.OptionChain.Results[0].OptionsArray[0].Puts
// 	for _, c := range callArray {

// 	}
// 	for _, p := range putArray {

// 	}
// 	options = append(callArray, putArray...)
// 	return options, nil
// }

func (option YahooOption) GetSymbol() string {
	str := option.ContractSymbol
	return str[0:(len(str) - 15)]
}

func (option YahooOption) GetOptionType() string {
	str := option.ContractSymbol
	return str[(len(str) - 9):(len(str) - 8)]
}
