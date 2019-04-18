package main

import (
	//"fmt"
	//"encoding/json"
)

type Stock struct {
	Symbol string `json:"symbol"`
	Sharesoutstanding int64 `json:"sharesOutstanding"`
	Bookvalue float32 `json:"bookValue"`
	Fiftydayaverage float32 `json:"fiftyDayAverage"`
	Fiftydayaveragechange float32 `json:"fiftyDayAverageChange"`
	Fiftydayaveragechangepercent float32 `json:"fiftyDayAverageChangePercent"`
	Twohundreddayaverage float32 `json:"twoHundredDayAverage"`
	Twohundreddayaveragechange float32 `json:"twoHundredDayAverageChange"`
	Forwardpe float32 `json:"forwardPE"`
	Regularmarketprice float32 `json:"regularMarketPrice"`
	Regularmarkettime int64 `json:"regularMarketTime"`
	Regularmarketchange float32 `json:"regularMarketChange"`
	Regularmarketopen float32 `json:"regularMarketOpen"`
	Regularmarketdayhigh float32 `json:"regularMarketDayHigh"`
	Regularmarketdaylow float32 `json:"regularMarketDayLow"`
	Regularmarketvolume int64 `json:"regularMarketVolume"`
	Epstrailingtwelvemonths float32 `json:"epsTrailingTwelveMonths"`
	Postmarketchangepercent float32 `json:"postMarketChangePercent"`
	Postmarkettime int64 `json:"postMarketTime"`
	Postmarketprice float32 `json:"postMarketPrice"`
	Postmarketchange float32 `json:"postMarketChange"`
	Regularmarketchangepercent float32 `json:"regularMarketChangePercent"`
	Regularmarketdayrange string `json:"regularMarketDayRange"`
	Regularmarketpreviousclose float32 `json:"regularMarketPreviousClose"`
	Bid float32 `json:"bid"`
	Ask float32 `json:"ask"`
	Bidsize int64 `json:"bidSize"`
	Asksize int64 `json:"askSize"`
	Averagedailyvolume3Month int64 `json:"averageDailyVolume3Month"`
	Averagedailyvolume10Day int64 `json:"averageDailyVolume10Day"`
	Fiftytwoweeklowchange float32 `json:"fiftyTwoWeekLowChange"`
	Fiftytwoweeklowchangepercent float32 `json:"fiftyTwoWeekLowChangePercent"`
	Fiftytwoweekrange string `json:"fiftyTwoWeekRange"`
	Fiftytwoweekhighchange float32 `json:"fiftyTwoWeekHighChange"`
	Fiftytwoweekhighchangepercent float32 `json:"fiftyTwoWeekHighChangePercent"`
	Fiftytwoweeklow float32 `json:"fiftyTwoWeekLow"`
	Fiftytwoweekhigh float32 `json:"fiftyTwoWeekHigh"`
	Dividenddate int64 `json:"dividendDate"`
	Earningstimestamp int64 `json:"earningsTimestamp"`
	Earningstimestampstart int64 `json:"earningsTimestampStart"`
	Earningstimestampend int64 `json:"earningsTimestampEnd"`
	Trailingannualdividendrate float32 `json:"trailingAnnualDividendRate"`
	Trailingpe float32 `json:"trailingPE"`
	Trailingannualdividendyield float32 `json:"trailingAnnualDividendYield"`
	Marketcap int64 `json:"marketCap"`
	Twohundreddayaveragechangepercent float32 `json:"twoHundredDayAverageChangePercent"`
	Epsforward float32 `json:"epsForward"`
	Pricetobook float32 `json:"priceToBook"`
	Sourceinterval int `json:"sourceInterval"`
}