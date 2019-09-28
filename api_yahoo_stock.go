package main

//ApiYahooQuote is a struct of the json object from yahoo stock
type ApiYahooQuote struct {
	Symbol                            string  `json:"symbol"`
	RegularMarketChange               float32 `json:"regularMarketChange"`
	RegularMarketOpen                 float32 `json:"regularMarketOpen"`
	RegularMarketDayHigh              float32 `json:"regularMarketDayHigh"`
	RegularMarketDayLow               float32 `json:"regularMarketDayLow"`
	RegularMarketVolume               int64   `json:"regularMarketVolume"`
	RegularMarketChangePercent        float32 `json:"regularMarketChangePercent"`
	RegularMarketPreviousClose        float32 `json:"regularMarketPreviousClose"`
	RegularMarketPrice                float32 `json:"regularMarketPrice"`
	RegularMarketTime                 int64   `json:"regularMarketTime"`
	EarningsTimestamp                 int64   `json:"earningsTimestamp"`
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
	BidSize                           int64   `json:"bidSize"`
	AskSize                           int64   `json:"askSize"`
	AverageDailyVolume3Month          int64   `json:"averageDailyVolume3Month"`
	AverageDailyVolume10Day           int64   `json:"averageDailyVolume10Day"`
	FiftyTwoWeekLowChange             float32 `json:"fiftyTwoWeekLowChange"`
	FiftyTwoWeekLowChangePercent      float32 `json:"fiftyTwoWeekLowChangePercent"`
	FiftyTwoWeekHighChange            float32 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekHighChangePercent     float32 `json:"fiftyTwoWeekHighChangePercent"`
	FiftyTwoWeekLow                   float32 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh                  float32 `json:"fiftyTwoWeekHigh"`
}

//isMarketOpen returns true if the market state is REGULAR
func (quote *ApiYahooQuote) isMarketOpen() bool {
	return quote.MarketState == "REGULAR"
}

//isMarketPreOpen returns true if the market state is PRE
func (quote *ApiYahooQuote) isMarketPreOpen() bool {
	return quote.MarketState == "PRE"
}

//Following functions are required for IStockData Interface

//GetSymbol returns the symbol
func (quote *ApiYahooQuote) GetSymbol() string {
	return quote.Symbol
}

//GetRegularMarketChange returns the price change during regular market hours
func (quote *ApiYahooQuote) GetRegularMarketChange() float32 {
	return quote.RegularMarketChange
}

//GetRegularMarketOpen returns the price at market open
func (quote *ApiYahooQuote) GetRegularMarketOpen() float32 {
	return quote.RegularMarketOpen
}

//GetRegularMarketDayHigh returns the highest price during regular market hours
func (quote *ApiYahooQuote) GetRegularMarketDayHigh() float32 {
	return quote.RegularMarketDayHigh
}

//GetRegularMarketDayLow returns the lowest price during regular market hours
func (quote *ApiYahooQuote) GetRegularMarketDayLow() float32 {
	return quote.RegularMarketDayLow
}

//GetRegularMarketVolume returns the volume during regular market hours
func (quote *ApiYahooQuote) GetRegularMarketVolume() int64 {
	return quote.RegularMarketVolume
}

//GetRegularMarketChangePercent returns the price percent change during regular market hours
func (quote *ApiYahooQuote) GetRegularMarketChangePercent() float32 {
	return quote.RegularMarketChangePercent
}

//GetRegularMarketPreviousClose returns the close price of regular market on the previous day
func (quote *ApiYahooQuote) GetRegularMarketPreviousClose() float32 {
	return quote.RegularMarketPreviousClose
}

//GetRegularMarketPrice returns the price of regular market
func (quote *ApiYahooQuote) GetRegularMarketPrice() float32 {
	return quote.RegularMarketPrice
}

//GetRegularMarketTime returns the unix time of regular market
func (quote *ApiYahooQuote) GetRegularMarketTime() int64 {
	return quote.RegularMarketTime
}

//GetEarningsTimestamp returns the unix timestamp of the earning
func (quote *ApiYahooQuote) GetEarningsTimestamp() int64 {
	return ConvertUTCUnixTimeInYYYYMMDD(quote.EarningsTimestamp)
}

//GetFiftyDayAverage returns the average price of past 50 days
func (quote *ApiYahooQuote) GetFiftyDayAverage() float32 {
	return quote.FiftyDayAverage
}

//GetFiftyDayAverageChange returns the average price change of past 50 days
func (quote *ApiYahooQuote) GetFiftyDayAverageChange() float32 {
	return quote.FiftyDayAverageChange
}

//GetFiftyDayAverageChangePercent returns the average price percent change of past 50 days
func (quote *ApiYahooQuote) GetFiftyDayAverageChangePercent() float32 {
	return quote.FiftyDayAverageChangePercent
}

//GetTwoHundredDayAverage returns the average price of past 200 days
func (quote *ApiYahooQuote) GetTwoHundredDayAverage() float32 {
	return quote.TwoHundredDayAverage
}

//GetTwoHundredDayAverageChange returns the average price change of past 200 days
func (quote *ApiYahooQuote) GetTwoHundredDayAverageChange() float32 {
	return quote.TwoHundredDayAverageChange
}

//GetTwoHundredDayAverageChangePercent returns the average price percent change of past 200 days
func (quote *ApiYahooQuote) GetTwoHundredDayAverageChangePercent() float32 {
	return quote.TwoHundredDayAverageChangePercent
}

//IsTradeable returns whether the stock is tradeable
func (quote *ApiYahooQuote) IsTradeable() bool {
	return quote.Tradeable
}

//GetMarketState returns the state of market
func (quote *ApiYahooQuote) GetMarketState() string {
	return quote.MarketState
}

//GetPostMarketChangePercent returns the price percent change during post market
func (quote *ApiYahooQuote) GetPostMarketChangePercent() float32 {
	return quote.PostMarketChangePercent
}

//GetPostMarketTime returns the unix time of post market
func (quote *ApiYahooQuote) GetPostMarketTime() int64 {
	return quote.PostMarketTime
}

//GetPostMarketPrice returns the price during post market
func (quote *ApiYahooQuote) GetPostMarketPrice() float32 {
	return quote.PostMarketPrice
}

//GetPostMarketChange returns the price change during post market
func (quote *ApiYahooQuote) GetPostMarketChange() float32 {
	return quote.PostMarketChange
}

//GetBid returns the bid
func (quote *ApiYahooQuote) GetBid() float32 {
	return quote.Bid
}

//GetAsk returns the ask
func (quote *ApiYahooQuote) GetAsk() float32 {
	return quote.Ask
}

//GetBidSize returns the bid size
func (quote *ApiYahooQuote) GetBidSize() int64 {
	return quote.BidSize
}

//GetAskSize returns the ask size
func (quote *ApiYahooQuote) GetAskSize() int64 {
	return quote.AskSize
}

//GetAverageDailyVolume3Month returns the average daily volume in the past 3 months
func (quote *ApiYahooQuote) GetAverageDailyVolume3Month() int64 {
	return quote.AverageDailyVolume3Month
}

//GetAverageDailyVolume10Day returns the average daily volume in the past 10 days
func (quote *ApiYahooQuote) GetAverageDailyVolume10Day() int64 {
	return quote.AverageDailyVolume10Day
}

//GetFiftyTwoWeekLowChange returns the change of lowest price in the past 52 weeks
func (quote *ApiYahooQuote) GetFiftyTwoWeekLowChange() float32 {
	return quote.FiftyTwoWeekLowChange
}

//GetFiftyTwoWeekLowChangePercent returns the percent change of lowest price in the past 52 weeks
func (quote *ApiYahooQuote) GetFiftyTwoWeekLowChangePercent() float32 {
	return quote.FiftyTwoWeekLowChangePercent
}

//GetFiftyTwoWeekHighChange returns the change of highest price in the past 52 weeks
func (quote *ApiYahooQuote) GetFiftyTwoWeekHighChange() float32 {
	return quote.FiftyTwoWeekHighChange
}

//GetFiftyTwoWeekHighChangePercent returns the percent change of highest price in the past 52 weeks
func (quote *ApiYahooQuote) GetFiftyTwoWeekHighChangePercent() float32 {
	return quote.FiftyTwoWeekHighChangePercent
}

//GetFiftyTwoWeekLow returns the lowest price in the past 52 weeks
func (quote *ApiYahooQuote) GetFiftyTwoWeekLow() float32 {
	return quote.FiftyTwoWeekLow
}

//GetFiftyTwoWeekHigh returns the highest price in the past 52 weeks
func (quote *ApiYahooQuote) GetFiftyTwoWeekHigh() float32 {
	return quote.FiftyTwoWeekHigh
}

//Following functions are required for IStockHist Interface

//GetMarketOpen returns the highest price in the past 52 weeks
func (quote *ApiYahooQuote) GetMarketOpen() float32 {
	return quote.RegularMarketOpen
}

//GetMarketHigh returns the highest price in the past 52 weeks
func (quote *ApiYahooQuote) GetMarketHigh() float32 {
	return quote.RegularMarketDayHigh
}

//GetMarketLow returns the highest price in the past 52 weeks
func (quote *ApiYahooQuote) GetMarketLow() float32 {
	return quote.RegularMarketDayLow
}

//GetMarketClose returns the highest price in the past 52 weeks
func (quote *ApiYahooQuote) GetMarketClose() float32 {
	return quote.RegularMarketPrice
}

//GetVolume returns the highest price in the past 52 weeks
func (quote *ApiYahooQuote) GetVolume() int64 {
	return quote.RegularMarketVolume
}
