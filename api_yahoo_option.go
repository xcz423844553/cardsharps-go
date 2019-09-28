package main

//ApiYahooOption is a struct of the json object from yahoo option
type ApiYahooOption struct {
	ContractSymbol    string  `json:"contractSymbol"`
	Strike            float32 `json:"strike"`
	LastPrice         float32 `json:"lastPrice"`
	PercentChange     float32 `json:"percentChange"`
	Volume            int64   `json:"volume"`
	OpenInterest      int64   `json:"openInterest"`
	Bid               float32 `json:"bid"`
	Ask               float32 `json:"ask"`
	Expiration        int64   `json:"expiration"`
	ImpliedVolatility float32 `json:"impliedVolatility"`
	LastTradeDate     int64   `json:"lastTradeDate"`
}

//GetContractSymbol returns the contract symbol
func (option *ApiYahooOption) GetContractSymbol() string {
	return option.ContractSymbol
}

//GetSymbol returns the symbol
func (option *ApiYahooOption) GetSymbol() string {
	str := option.ContractSymbol
	return str[0:(len(str) - 15)]
}

//GetOptionType returns the option type
func (option *ApiYahooOption) GetOptionType() string {
	str := option.ContractSymbol
	return str[(len(str) - 9):(len(str) - 8)]
}

//GetStrike returns the strike price
func (option *ApiYahooOption) GetStrike() float32 {
	return option.Strike
}

//GetLastPrice returns the option price of last day
func (option *ApiYahooOption) GetLastPrice() float32 {
	return option.LastPrice
}

//GetVolume returns the volume
func (option *ApiYahooOption) GetVolume() int64 {
	return option.Volume
}

//GetOpenInterest returns the open interest
func (option *ApiYahooOption) GetOpenInterest() int64 {
	return option.OpenInterest
}

//GetImpliedVolatility returns the IV
func (option *ApiYahooOption) GetImpliedVolatility() float32 {
	return option.ImpliedVolatility
}

//GetPercentChange returns the percent change of the price
func (option *ApiYahooOption) GetPercentChange() float32 {
	return option.PercentChange
}

//GetBid returns the price of the bid
func (option *ApiYahooOption) GetBid() float32 {
	return option.Bid
}

//GetAsk returns the price of the ask
func (option *ApiYahooOption) GetAsk() float32 {
	return option.Ask
}

//GetExpiration converts the expiration date in unix timestamp to YYYYMMDD format
func (option *ApiYahooOption) GetExpiration() int64 {
	return ConvertUTCUnixTimeInYYYYMMDD(option.Expiration)
}

//GetLastTradeDate converts the last trade date in unix timestamp to YYYYMMDD format
func (option *ApiYahooOption) GetLastTradeDate() int64 {
	return ConvertUTCUnixTimeInYYYYMMDD(option.LastTradeDate)
}
