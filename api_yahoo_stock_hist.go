package main

//ApiYahooStockHist is a struct of the history data from yahoo finance
type ApiYahooStockHist struct {
	Symbol string  `json:"symbol"`
	Date   int64   `json:"date"`
	Open   float32 `json:"open"`
	High   float32 `json:"high"`
	Low    float32 `json:"low"`
	Close  float32 `json:"close"`
	Volume int64   `json:"volume"`
}

//GetDate returns the date of the history
func (hist *ApiYahooStockHist) GetDate() int64 {
	return hist.Date
}

//Following functions are the implementation of interface IStockHist

//GetSymbol returns the symbol of the stock
func (hist *ApiYahooStockHist) GetSymbol() string {
	return hist.Symbol
}

//GetMarketOpen returns the open price of the stock
func (hist *ApiYahooStockHist) GetMarketOpen() float32 {
	return hist.Open
}

//GetMarketHigh returns the highest price of the stock
func (hist *ApiYahooStockHist) GetMarketHigh() float32 {
	return hist.High
}

//GetMarketLow returns the lowest price of the stock
func (hist *ApiYahooStockHist) GetMarketLow() float32 {
	return hist.Low
}

//GetMarketClose returns the close price of the stock
func (hist *ApiYahooStockHist) GetMarketClose() float32 {
	return hist.Close
}

//GetVolume returns the volume of the stock
func (hist *ApiYahooStockHist) GetVolume() int64 {
	return hist.Volume
}
