package main

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
