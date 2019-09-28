package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"
)

type YahooApi struct {
}

const (
	//YahooClientTimeout is the timeout when connecting to yahoo finance
	YahooClientTimeout = 10 * time.Second
	//URLOptionAndStockFromYahoo is the base url of yahoo finance
	URLOptionAndStockFromYahoo = "https://query1.finance.yahoo.com/v7/finance/options/"
)

//GetOptionAndStockURL returns the url of yahoo finance
func (api *YahooApi) GetOptionAndStockURL(symbol string) string {
	return URLOptionAndStockFromYahoo + symbol
}

//GetOptionURL returns the url of yahoo finance to get the options on a given expiration date
func (api *YahooApi) GetOptionURL(symbol string, exp int64) string {
	if exp <= 0 {
		return api.GetOptionAndStockURL(symbol)
	}
	return URLOptionAndStockFromYahoo + symbol + "?date=" + strconv.FormatInt(exp, 10)
}

//GetExpirationDate returns the expiration dates of the options from yahoo finance
func (api *YahooApi) GetExpirationDate(symbol string) ([]int64, error) {
	var expDates []int64
	var expDateErr error
	url := api.GetOptionAndStockURL(symbol)
	resp, connError := http.Get(url)
	if connError != nil {
		return expDates, connError
	}
	defer resp.Body.Close()
	body, parseError := ioutil.ReadAll(resp.Body)
	if parseError != nil {
		return expDates, parseError
	}
	var y ApiYahooResponse
	if jsonError := json.Unmarshal(body, &y); jsonError != nil {
		return expDates, jsonError
	}
	expDates, expDateErr = y.GetExpirationDate()
	if expDateErr != nil {
		return expDates, expDateErr
	}
	return expDates, nil
}

//GetQuote returns the quote from yahoo finance
func (api *YahooApi) GetQuote(symbol string) (ApiYahooQuote, error) {
	var stock ApiYahooQuote
	var stockErr error
	url := api.GetOptionAndStockURL(symbol)
	resp, connError := http.Get(url)
	if connError != nil {
		return stock, connError
	}
	defer resp.Body.Close()
	body, parseError := ioutil.ReadAll(resp.Body)
	if parseError != nil {
		return stock, parseError
	}
	var y ApiYahooResponse
	if jsonError := json.Unmarshal(body, &y); jsonError != nil {
		return stock, jsonError
	}
	stock, stockErr = y.GetQuote()
	if stockErr != nil {
		return stock, stockErr
	}
	return stock, nil
}

//GetOptionChain returns an array of options of a given expiration date; if the expiration date is 0, return the upcoming date
func (api *YahooApi) GetOptionChain(symbol string, expDate int64) ([]ApiYahooOption, error) {
	var options []ApiYahooOption
	var optionErr error
	url := api.GetOptionURL(symbol, expDate)
	resp, connError := http.Get(url)
	if connError != nil {
		return options, connError
	}
	defer resp.Body.Close()
	body, parseError := ioutil.ReadAll(resp.Body)
	if parseError != nil {
		return options, parseError
	}
	var y ApiYahooResponse
	if jsonError := json.Unmarshal(body, &y); jsonError != nil {
		return options, jsonError
	}
	options, optionErr = y.GetOptionChain()
	if optionErr != nil {
		return options, optionErr
	}
	return options, nil
}

//GetOptionChainStockAndExpDate returns the option chain, stock quote, and expiration dates of a given symbol at a given expiration date
//If expDate is 0, return the data of the upcoming date
func (api *YahooApi) GetOptionChainStockAndExpDate(symbol string, expDate int64) ([]ApiYahooOption, ApiYahooQuote, []int64, error) {
	var options []ApiYahooOption
	var stock ApiYahooQuote
	var expDates []int64
	var optionErr error
	var stockErr error
	var expDateErr error
	url := api.GetOptionURL(symbol, expDate)
	resp, connError := http.Get(url)
	if connError != nil {
		return options, stock, expDates, connError
	}
	defer resp.Body.Close()
	body, parseError := ioutil.ReadAll(resp.Body)
	if parseError != nil {
		return options, stock, expDates, parseError
	}
	var y ApiYahooResponse
	if jsonError := json.Unmarshal(body, &y); jsonError != nil {
		return options, stock, expDates, jsonError
	}
	options, optionErr = y.GetOptionChain()
	stock, stockErr = y.GetQuote()
	expDates, expDateErr = y.GetExpirationDate()
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

//GetAllOptionChainStockAndExpDate returns the option chains of all expiration dates, stock quote, and expiration dates of a given symbol
func (api *YahooApi) GetAllOptionChainStockAndExpDate(symbol string) ([]ApiYahooOption, ApiYahooQuote, []int64, error) {
	var options []ApiYahooOption
	_, stock, expDates, quoteErr := api.GetOptionChainStockAndExpDate(symbol, 0)
	if quoteErr != nil {
		return options, stock, expDates, quoteErr
	}
	for _, expDate := range expDates {
		newOptions, optionErr := api.GetOptionChain(symbol, expDate)
		if optionErr != nil {
			fmt.Println("Error encounted while getting option chain data for " + symbol + " on date " + strconv.FormatInt(expDate, 10))
			fmt.Println(optionErr)
			continue
		}
		options = append(options, newOptions...)
	}
	return options, stock, expDates, nil
}

//GetStockHist downloads the historical stock data from finance.yahoo.com
//PARAM: symbol - the symbol of the stock to be downloaded
//PARAM: start date - start date of the stock data
//PARAM: end date - end date of the stock data
//Return: yahoo stock history - an array of stock history data
func (api *YahooApi) GetStockHist(symbol string, startDate int, endDate int) ([]ApiYahooStockHist, error) {
	var stockHistArray []ApiYahooStockHist
	start := ConvertTimeInUnix(startDate)
	end := ConvertTimeInUnix(endDate)
	//Get Crumb
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: YahooClientTimeout,
		Jar:     jar,
	}
	yahooReq, yahooReqErr := http.NewRequest("GET", "https://finance.yahoo.com", nil)
	if yahooReqErr != nil {
		return stockHistArray, yahooReqErr
	}
	yahooReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; U; Linux i686) Gecko/20071127 Firefox/2.0.0.11")
	resp, _ := client.Do(yahooReq)

	//Test Crumb
	crumbReq, crumbErr := http.NewRequest("GET", "https://query1.finance.yahoo.com/v1/test/getcrumb", nil)
	if crumbErr != nil {
		return stockHistArray, crumbErr
	}
	crumbReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; U; Linux i686) Gecko/20071127 Firefox/2.0.0.11")
	resp, _ = client.Do(crumbReq)

	//Read Crumb
	reader := csv.NewReader(resp.Body)
	crumb, crumbErr := reader.Read()
	if crumbErr != nil {
		return stockHistArray, crumbErr
	}

	//Download stock historical data
	url := fmt.Sprintf(URL_STOCK_HIST, symbol, start, end, crumb[0])
	resp, clientErr := client.Get(url)
	if clientErr != nil {
		return stockHistArray, clientErr
	}
	defer resp.Body.Close()

	//Parse stock historical data and store in the array
	var csvData [][]string
	reader = csv.NewReader(resp.Body)
	reader.LazyQuotes = true
	csvData, readerErr := reader.ReadAll()
	if readerErr != nil {
		return stockHistArray, readerErr
	}
	for idx, row := range csvData {
		if idx == 0 {
			continue
		}
		colDate := ConvertTimeInYYYYMMDD(row[0])
		colOpen, _ := strconv.ParseFloat(row[1], 32)
		colHigh, _ := strconv.ParseFloat(row[2], 32)
		colLow, _ := strconv.ParseFloat(row[3], 32)
		colClose, _ := strconv.ParseFloat(row[5], 32)
		colVolume, _ := strconv.Atoi(row[6])
		stockHist := ApiYahooStockHist{
			Symbol: symbol,
			Date:   int64(colDate),
			Open:   float32(colOpen),
			High:   float32(colHigh),
			Low:    float32(colLow),
			Close:  float32(colClose),
			Volume: int64(colVolume),
		}
		stockHistArray = append(stockHistArray, stockHist)
	}
	return stockHistArray, nil
}
