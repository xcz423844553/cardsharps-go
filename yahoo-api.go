package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"
)

type YahooAPIManager struct {
}

type YahooStockHist struct {
	Date   int     `json:"date"`
	Open   float32 `json:"open"`
	High   float32 `json:"high"`
	Low    float32 `json:"low"`
	Close  float32 `json:"close"`
	Volume int     `json:"volume"`
}

const (
	ClientTimeout = 10 * time.Second
)

func (manager *YahooAPIManager) GetOptionsAndStockDataBySymbol(symbol string) ([]YahooOption, YahooQuote, []int64, error) {
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

func (manager *YahooAPIManager) GetOptionsAndStockDataBySymbolAndExpDate(symbol string, expDate int64) ([]YahooOption, YahooQuote, []int64, error) {
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
		fmt.Println(string(body))
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
func (manager *YahooAPIManager) GetStockDataBySymbol(symbol string) (YahooQuote, error) {
	var stock YahooQuote
	var stockErr error
	url := URL_OPTION + symbol
	resp, connError := http.Get(url)
	if connError != nil {
		fmt.Println("ConnError: ", connError)
		return stock, connError
	}
	defer resp.Body.Close()
	body, parseError := ioutil.ReadAll(resp.Body)
	if parseError != nil {
		fmt.Println("ParseError: ", parseError)
		return stock, parseError
	}
	var y YahooResponse
	if jsonError := json.Unmarshal(body, &y); jsonError != nil {
		fmt.Println("JsonError: ", jsonError)
		return stock, parseError
	}
	stock, stockErr = y.GetQuote()
	if stockErr != nil {
		return stock, stockErr
	}
	return stock, nil
}

//GetStockHistoryFromYahoo downloads the historical stock data from finance.yahoo.com
//PARAM: symbol - the symbol of the stock to be downloaded
//PARAM: start date - start date of the stock data
//PARAM: end date - end date of the stock data
//Return: yahoo stock history - an array of stock history data
func (manager *YahooAPIManager) GetStockHistoryFromYahoo(symbol string, startDate int, endDate int) ([]YahooStockHist, error) {
	var stockHistArray []YahooStockHist

	start := ConvertTimeInUnix(startDate)
	end := ConvertTimeInUnix(endDate)

	//Get Crumb
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: ClientTimeout,
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
		PrintMsgInConsole(MSGERROR, LOGTYPE_YAHOO_API_MANAGER, "Crumb Error: "+crumbErr.Error())
		return stockHistArray, crumbErr
	}

	//Download stock historical data
	url := fmt.Sprintf(URL_STOCK_HIST, symbol, start, end, crumb[0])
	resp, clientErr := client.Get(url)
	if clientErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_YAHOO_API_MANAGER, "Client Error: "+clientErr.Error())
		return stockHistArray, clientErr
	}
	defer resp.Body.Close()

	//Parse stock historical data and store in the array
	var csvData [][]string
	reader = csv.NewReader(resp.Body)
	reader.LazyQuotes = true
	csvData, readerErr := reader.ReadAll()
	if readerErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_YAHOO_API_MANAGER, "Reader Error: "+readerErr.Error())
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
		stockHist := YahooStockHist{
			Date:   colDate,
			Open:   float32(colOpen),
			High:   float32(colHigh),
			Low:    float32(colLow),
			Close:  float32(colClose),
			Volume: colVolume,
		}
		stockHistArray = append(stockHistArray, stockHist)
	}
	return stockHistArray, nil
}
