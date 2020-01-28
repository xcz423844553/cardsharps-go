package main

import (
	"strings"

	tda "github.com/xcz423844553/td_ameritrade_client_golang"
)

//ApiTdaOptionWrapper is a struct of the json object from yahoo option
type ApiTdaOptionWrapper struct {
	ContractSymbol    string  `json:"contractSymbol"`
	Symbol            string  `json:"symbol"`
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
	PutCall           string  `json:"putCall"` //enum[PUT, CALL]
}

func (option *ApiTdaOptionWrapper) BuildWrapper(opt tda.Option) {
	option.ContractSymbol = opt.Symbol
	option.Symbol = opt.Symbol
	option.Strike = float32(opt.StrikePrice)
	option.LastPrice = float32(opt.Last)
	option.PercentChange = float32(opt.PercentChange)
	option.Volume = opt.TotalVolume
	option.OpenInterest = int64(opt.OpenInterest)
	option.Bid = float32(opt.Bid)
	option.Ask = float32(opt.Ask)
	option.Expiration = ConvertUTCUnixTimeInYYYYMMDD(opt.ExpirationDate)
	option.ImpliedVolatility = float32(opt.Volatility)
	option.LastTradeDate = ConvertUTCUnixTimeInYYYYMMDD(opt.TradeTimeInLong)
	option.PutCall = opt.PutCall
}

//GetContractSymbol returns the contract symbol
func (option *ApiTdaOptionWrapper) GetContractSymbol() string {
	return option.ContractSymbol
}

//GetSymbol returns the symbol
func (option *ApiTdaOptionWrapper) GetSymbol() string {
	return strings.Split(option.Symbol, "_")[0]
}

//GetOptionType returns the option type
func (option *ApiTdaOptionWrapper) GetOptionType() string {
	return string([]rune(option.PutCall)[0])
}

//GetStrike returns the strike price
func (option *ApiTdaOptionWrapper) GetStrike() float32 {
	return option.Strike
}

//GetLastPrice returns the option price of last day
func (option *ApiTdaOptionWrapper) GetLastPrice() float32 {
	return option.LastPrice
}

//GetVolume returns the volume
func (option *ApiTdaOptionWrapper) GetVolume() int64 {
	return option.Volume
}

//GetOpenInterest returns the open interest
func (option *ApiTdaOptionWrapper) GetOpenInterest() int64 {
	return option.OpenInterest
}

//GetImpliedVolatility returns the IV
func (option *ApiTdaOptionWrapper) GetImpliedVolatility() float32 {
	return option.ImpliedVolatility
}

//GetPercentChange returns the percent change of the price
func (option *ApiTdaOptionWrapper) GetPercentChange() float32 {
	return option.PercentChange
}

//GetBid returns the price of the bid
func (option *ApiTdaOptionWrapper) GetBid() float32 {
	return option.Bid
}

//GetAsk returns the price of the ask
func (option *ApiTdaOptionWrapper) GetAsk() float32 {
	return option.Ask
}

//GetExpiration converts the expiration date in unix timestamp to YYYYMMDD format
func (option *ApiTdaOptionWrapper) GetExpiration() int64 {
	return ConvertUTCUnixTimeInYYYYMMDD(option.Expiration)
}

//GetLastTradeDate converts the last trade date in unix timestamp to YYYYMMDD format
func (option *ApiTdaOptionWrapper) GetLastTradeDate() int64 {
	return ConvertUTCUnixTimeInYYYYMMDD(option.LastTradeDate)
}
