package main

import (
	"fmt"
	"testing"
)

func TestDbStockData(t *testing.T) {
	dao := new(DaoStockData)
	if dbDropErr := dao.DropTableIfExist(); dbDropErr != nil {
		t.Error(dbDropErr.Error())
	}
	if dbCreateErr := dao.CreateTableIfNotExist(); dbCreateErr != nil {
		t.Error(dbCreateErr.Error())
	}
	symbol := "AAPL"
	date := GetTimeInYYYYMMDD64()
	row := RowStockData{
		Symbol:                            symbol,
		RegularMarketChange:               0.80000305,
		RegularMarketOpen:                 219.96,
		RegularMarketDayHigh:              220.82,
		RegularMarketDayLow:               219.13,
		RegularMarketVolume:               16092987,
		RegularMarketChangePercent:        0.13140397,
		RegularMarketPreviousClose:        219.9,
		RegularMarketPrice:                220.7,
		RegularMarketTime:                 1568750401,
		EarningsTimestamp:                 1564531200,
		FiftyDayAverage:                   208.83485,
		FiftyDayAverageChange:             11.865143,
		FiftyDayAverageChangePercent:      0.05681591,
		TwoHundredDayAverage:              198.22174,
		TwoHundredDayAverageChange:        22.478256,
		TwoHundredDayAverageChangePercent: 0.11339955,
		Tradeable:                         true,
		MarketState:                       "CLOSE",
		PostMarketChangePercent:           0.13140397,
		PostMarketTime:                    1568764782,
		PostMarketPrice:                   220.99,
		PostMarketChange:                  0.29000854,
		Bid:                               0,
		Ask:                               0,
		BidSize:                           10,
		AskSize:                           8,
		AverageDailyVolume3Month:          26355010,
		AverageDailyVolume10Day:           30691957,
		FiftyTwoWeekLowChange:             78.7,
		FiftyTwoWeekLowChangePercent:      0.5542253,
		FiftyTwoWeekHighChange:            -12.770004,
		FiftyTwoWeekHighChangePercent:     -0.054696552,
		FiftyTwoWeekLow:                   142,
		FiftyTwoWeekHigh:                  233.47,
	}
	input := ApiYahooQuote{
		Symbol:                            symbol,
		RegularMarketChange:               0.80000305,
		RegularMarketOpen:                 219.96,
		RegularMarketDayHigh:              220.82,
		RegularMarketDayLow:               219.13,
		RegularMarketVolume:               16092987,
		RegularMarketChangePercent:        0.13140397,
		RegularMarketPreviousClose:        219.9,
		RegularMarketPrice:                220.7,
		RegularMarketTime:                 1568750401,
		EarningsTimestamp:                 1564531200,
		FiftyDayAverage:                   208.83485,
		FiftyDayAverageChange:             11.865143,
		FiftyDayAverageChangePercent:      0.05681591,
		TwoHundredDayAverage:              198.22174,
		TwoHundredDayAverageChange:        22.478256,
		TwoHundredDayAverageChangePercent: 0.11339955,
		Tradeable:                         true,
		MarketState:                       "CLOSE",
		PostMarketChangePercent:           0.13140397,
		PostMarketTime:                    1568764782,
		PostMarketPrice:                   220.99,
		PostMarketChange:                  0.29000854,
		Bid:                               0,
		Ask:                               0,
		BidSize:                           10,
		AskSize:                           8,
		AverageDailyVolume3Month:          26355010,
		AverageDailyVolume10Day:           30691957,
		FiftyTwoWeekLowChange:             78.7,
		FiftyTwoWeekLowChangePercent:      0.5542253,
		FiftyTwoWeekHighChange:            -12.770004,
		FiftyTwoWeekHighChangePercent:     -0.054696552,
		FiftyTwoWeekLow:                   142,
		FiftyTwoWeekHigh:                  233.47,
	}
	if dbInsertErr := dao.InsertOrUpdateStockData(&input); dbInsertErr != nil {
		t.Error(dbInsertErr.Error())
	}
	if dbUpdateErr := dao.InsertOrUpdateStockData(&input); dbUpdateErr != nil {
		t.Error(dbUpdateErr.Error())
	}

	stock, dbStockErr := dao.SelectStockDataBySymbolAndDate(symbol, date)
	if dbStockErr != nil {
		t.Error(dbStockErr.Error())
	}
	if row.Symbol != stock.Symbol ||
		fmt.Sprintf("%.2f", row.RegularMarketChange) != fmt.Sprintf("%.2f", stock.RegularMarketChange) ||
		fmt.Sprintf("%.2f", row.RegularMarketOpen) != fmt.Sprintf("%.2f", stock.RegularMarketOpen) ||
		fmt.Sprintf("%.2f", row.RegularMarketDayHigh) != fmt.Sprintf("%.2f", stock.RegularMarketDayHigh) ||
		fmt.Sprintf("%.2f", row.RegularMarketDayLow) != fmt.Sprintf("%.2f", stock.RegularMarketDayLow) ||
		row.RegularMarketVolume != stock.RegularMarketVolume ||
		fmt.Sprintf("%.2f", row.RegularMarketChangePercent) != fmt.Sprintf("%.2f", stock.RegularMarketChangePercent) ||
		fmt.Sprintf("%.2f", row.RegularMarketPreviousClose) != fmt.Sprintf("%.2f", stock.RegularMarketPreviousClose) ||
		fmt.Sprintf("%.2f", row.RegularMarketPrice) != fmt.Sprintf("%.2f", stock.RegularMarketPrice) ||
		row.RegularMarketTime != stock.RegularMarketTime ||
		ConvertUTCUnixTimeInYYYYMMDD(row.EarningsTimestamp) != stock.EarningsTimestamp ||
		fmt.Sprintf("%.2f", row.FiftyDayAverage) != fmt.Sprintf("%.2f", stock.FiftyDayAverage) ||
		fmt.Sprintf("%.2f", row.FiftyDayAverageChange) != fmt.Sprintf("%.2f", stock.FiftyDayAverageChange) ||
		fmt.Sprintf("%.2f", row.FiftyDayAverageChangePercent) != fmt.Sprintf("%.2f", stock.FiftyDayAverageChangePercent) ||
		fmt.Sprintf("%.2f", row.TwoHundredDayAverage) != fmt.Sprintf("%.2f", stock.TwoHundredDayAverage) ||
		fmt.Sprintf("%.2f", row.TwoHundredDayAverageChange) != fmt.Sprintf("%.2f", stock.TwoHundredDayAverageChange) ||
		fmt.Sprintf("%.2f", row.TwoHundredDayAverageChangePercent) != fmt.Sprintf("%.2f", stock.TwoHundredDayAverageChangePercent) ||
		fmt.Sprintf("%.2f", stock.PostMarketChangePercent) != fmt.Sprintf("%.2f", stock.PostMarketChangePercent) ||
		row.PostMarketTime != stock.PostMarketTime ||
		fmt.Sprintf("%.2f", row.PostMarketPrice) != fmt.Sprintf("%.2f", stock.PostMarketPrice) ||
		fmt.Sprintf("%.2f", row.PostMarketChange) != fmt.Sprintf("%.2f", stock.PostMarketChange) ||
		fmt.Sprintf("%.2f", row.Bid) != fmt.Sprintf("%.2f", stock.Bid) ||
		fmt.Sprintf("%.2f", row.Ask) != fmt.Sprintf("%.2f", stock.Ask) ||
		row.BidSize != stock.BidSize ||
		row.AskSize != stock.AskSize ||
		row.AverageDailyVolume3Month != stock.AverageDailyVolume3Month ||
		row.AverageDailyVolume10Day != stock.AverageDailyVolume10Day ||
		fmt.Sprintf("%.2f", row.FiftyTwoWeekLowChange) != fmt.Sprintf("%.2f", stock.FiftyTwoWeekLowChange) ||
		fmt.Sprintf("%.2f", row.FiftyTwoWeekLowChangePercent) != fmt.Sprintf("%.2f", stock.FiftyTwoWeekLowChangePercent) ||
		fmt.Sprintf("%.2f", row.FiftyTwoWeekHighChange) != fmt.Sprintf("%.2f", stock.FiftyTwoWeekHighChange) ||
		fmt.Sprintf("%.2f", row.FiftyTwoWeekHighChangePercent) != fmt.Sprintf("%.2f", stock.FiftyTwoWeekHighChangePercent) ||
		fmt.Sprintf("%.2f", row.FiftyTwoWeekLow) != fmt.Sprintf("%.2f", stock.FiftyTwoWeekLow) ||
		fmt.Sprintf("%.2f", row.FiftyTwoWeekHigh) != fmt.Sprintf("%.2f", stock.FiftyTwoWeekHigh) {
		t.Error("Error: stock data selected by symbol and date is not matched.")
	}
}
