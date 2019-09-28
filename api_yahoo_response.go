package main

import (
	"errors"
)

//ApiYahooResponse is the struct of response from yahoo finance
type ApiYahooResponse struct {
	OptionChain struct {
		Results []struct {
			UnderlyingSymbol string                `json:"underlyingSymbol"`
			ExpirationDates  []int64               `json:"expirationDates"`
			Strikes          []float32             `json:"strikes"`
			HasMiniOptions   bool                  `json:"hasMiniOptions"`
			Quote            ApiYahooQuote         `json:"quote"`
			OptionsArray     []ApiYahooOptionArray `json:"options"`
		} `json:"result"`
		Error string `json:"error"`
	} `json:"optionChain"`
}

//ApiYahooOptionArray is the struct of the option array in the response from yahoo finance
type ApiYahooOptionArray struct {
	ExpirationDate int64            `json:"expirationDate"`
	Calls          []ApiYahooOption `json:"calls"`
	Puts           []ApiYahooOption `json:"puts"`
}

//isEmptyQuote returns true if there is no such quote
func (resp *ApiYahooResponse) isEmptyQuote() bool {
	resultsArray := resp.OptionChain.Results
	return len(resultsArray) == 0
}

//isEmptyOption returns true if there is no option chain
func (resp *ApiYahooResponse) isEmptyOption() bool {
	if resp.isEmptyQuote() {
		return true
	}
	return len(resp.OptionChain.Results[0].OptionsArray) == 0
}

//GetQuote returns the quote in response from yahoo finance
func (resp *ApiYahooResponse) GetQuote() (ApiYahooQuote, error) {
	var quote ApiYahooQuote
	if resp.isEmptyQuote() {
		return quote, errors.New("Quote response is empty")
	}
	quote = resp.OptionChain.Results[0].Quote
	return quote, nil
}

//GetOptions returns the option chain in response from yahoo finance
func (resp *ApiYahooResponse) GetOptionChain() ([]ApiYahooOption, error) {
	var options []ApiYahooOption
	if resp.isEmptyOption() {
		return options, errors.New("Option response is empty")
	}
	callArray := resp.OptionChain.Results[0].OptionsArray[0].Calls
	putArray := resp.OptionChain.Results[0].OptionsArray[0].Puts
	options = append(callArray, putArray...)
	return options, nil
}

//GetExpirationDate returns the expiration dates of the options in response from yahoo finance
func (resp *ApiYahooResponse) GetExpirationDate() ([]int64, error) {
	var expDates []int64
	if resp.isEmptyOption() {
		return expDates, errors.New("Expiration date response is empty")
	}
	numExpDates := len(resp.OptionChain.Results[0].ExpirationDates)
	expDates = make([]int64, numExpDates)
	for i := 0; i < numExpDates; i++ {
		expDates[i] = resp.OptionChain.Results[0].ExpirationDates[i]
	}
	return expDates, nil
}
