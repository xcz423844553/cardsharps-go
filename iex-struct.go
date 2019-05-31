package main

import "strconv"

const (
	IexChartRange_1d  = "1d"
	IexChartRange_1m  = "1m"
	IexChartRange_3m  = "3m"
	IexChartRange_6m  = "6m"
	IexChartRange_ytd = "ytd"
	IexChartRange_1y  = "1y"
	IexChartRange_2y  = "2y"
	IexChartRange_5y  = "5y"
)

type IexChart struct {
	Date   string `json:"date"`
	Minute string `json:"minute"`
	// Label                string  `json:"label"`
	High                 float32 `json:"high"`
	Low                  float32 `json:"low"`
	Average              float32 `json:"average"`
	Volume               int     `json:"volume"`
	Notional             float32 `json:"notional"`
	NumberOfTrades       int     `json:"numberOfTrades"`
	MarketHigh           float32 `json:"marketHigh"`
	MarketLow            float32 `json:"marketLow"`
	MarketAverage        float32 `json:"marketAverage"`
	MarketVolume         int     `json:"marketVolume"`
	MarketNotional       float32 `json:"marketNotional"`
	MarketNumberOfTrades int     `json:"marketNumberOfTrades"`
	Open                 float32 `json:"open"`
	Close                float32 `json:"close"`
	MarketOpen           float32 `json:"marketOpen"`
	MarketClose          float32 `json:"marketClose"`
	ChangeOverTime       float32 `json:"changeOverTime"`
	MarketChangeOverTime float32 `json:"marketChangeOverTime"`
}

type IexDayChart struct {
	Date             string  `json:"date"`
	Open             float32 `json:"open"`
	High             float32 `json:"high"`
	Low              float32 `json:"low"`
	Close            float32 `json:"close"`
	Volume           int     `json:"volume"`
	UnadjustedVolume int     `json:"unadjustedVolume"`
	Change           float32 `json:"change"`
	ChangePercent    float32 `json:"changePercent"`
	Vwap             float32 `json:"vwap"`
	//Label            string  `json:"label"`
	ChangeOverTime float32 `json:"changeOverTime"`
}

func (dc *IexDayChart) BuildFromCsv(strs ...string) *IexDayChart {
	dc.Date = strs[0]
	open, _ := strconv.ParseFloat(strs[1], 32)
	high, _ := strconv.ParseFloat(strs[2], 32)
	low, _ := strconv.ParseFloat(strs[3], 32)
	close, _ := strconv.ParseFloat(strs[4], 32)
	volume, _ := strconv.Atoi(strs[5])
	unadjustedVolume, _ := strconv.Atoi(strs[6])
	change, _ := strconv.ParseFloat(strs[7], 32)
	changePercent, _ := strconv.ParseFloat(strs[8], 32)
	vwap, _ := strconv.ParseFloat(strs[9], 32)
	changeOverTime, _ := strconv.ParseFloat(strs[10], 32)
	dc.Open = float32(open)
	dc.High = float32(high)
	dc.Low = float32(low)
	dc.Close = float32(close)
	dc.Volume = volume
	dc.UnadjustedVolume = unadjustedVolume
	dc.Change = float32(change)
	dc.ChangePercent = float32(changePercent)
	dc.Vwap = float32(vwap)
	dc.ChangeOverTime = float32(changeOverTime)
	return dc
}
