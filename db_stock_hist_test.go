package main

import (
	"fmt"
	"testing"
)

func TestDbStockHist(t *testing.T) {
	dao := new(DaoStockHist)
	if dbDropErr := dao.DropTableIfExist(); dbDropErr != nil {
		t.Error(dbDropErr.Error())
	}
	if dbCreateErr := dao.CreateTableIfNotExist(); dbCreateErr != nil {
		t.Error(dbCreateErr.Error())
	}
	symbol := "AAPL"
	row1 := RowStockHist2{
		Symbol:      symbol,
		MarketOpen:  219.96,
		MarketHigh:  220.82,
		MarketLow:   219.13,
		Volume:      16092987,
		MarketClose: 220.7,
	}
	row2 := RowStockHist2{
		Symbol:      symbol,
		MarketOpen:  2190.96,
		MarketHigh:  2200.82,
		MarketLow:   2190.13,
		Volume:      16087,
		MarketClose: 2200.7,
	}
	input1 := ApiYahooQuote{
		Symbol:               symbol,
		RegularMarketOpen:    219.96,
		RegularMarketDayHigh: 220.82,
		RegularMarketDayLow:  219.13,
		RegularMarketVolume:  16092987,
		RegularMarketPrice:   220.7,
	}
	input2 := ApiYahooQuote{
		Symbol:               symbol,
		RegularMarketOpen:    2190.96,
		RegularMarketDayHigh: 2200.82,
		RegularMarketDayLow:  2190.13,
		RegularMarketVolume:  16087,
		RegularMarketPrice:   2200.7,
	}
	if dbInsertErr := dao.InsertOrUpdateStockHist(&input1, 20190101); dbInsertErr != nil {
		t.Error(dbInsertErr.Error())
	}
	if dbInsertErr := dao.InsertOrUpdateStockHist(&input2, 20190102); dbInsertErr != nil {
		t.Error(dbInsertErr.Error())
	}
	if dbUpdateErr := dao.InsertOrUpdateStockHist(&input1, 20190101); dbUpdateErr != nil {
		t.Error(dbUpdateErr.Error())
	}

	stock, dbStockErr := dao.SelectLastStockHist(symbol)
	if dbStockErr != nil {
		t.Error(dbStockErr.Error())
	}
	if row2.Symbol != stock.Symbol ||
		fmt.Sprintf("%.2f", row2.MarketOpen) != fmt.Sprintf("%.2f", stock.MarketOpen) ||
		fmt.Sprintf("%.2f", row2.MarketHigh) != fmt.Sprintf("%.2f", stock.MarketHigh) ||
		fmt.Sprintf("%.2f", row2.MarketLow) != fmt.Sprintf("%.2f", stock.MarketLow) ||
		row2.Volume != stock.Volume ||
		fmt.Sprintf("%.2f", row2.MarketClose) != fmt.Sprintf("%.2f", stock.MarketClose) {
		fmt.Printf("%.2f, %.2f, %.2f, %.2f, %d, %s", row2.MarketOpen, row2.MarketClose, row2.MarketHigh, row2.MarketLow, row2.Volume, row2.Symbol)
		fmt.Printf("%.2f, %.2f, %.2f, %.2f, %d, %s", stock.MarketOpen, stock.MarketClose, stock.MarketHigh, stock.MarketLow, stock.Volume, stock.Symbol)
		t.Error("Error: last available stock hist selected by symbol is not matched.")
	}

	stockList, dbStockListErr := dao.SelectLastNumberStockHistBeforeDate(symbol, 2, 20190103)
	if dbStockListErr != nil {
		t.Error(dbStockListErr.Error())
	}
	if len(stockList) != 2 {
		t.Error("Error: last number of avaliable stock hists selected by symbol and before date are not matched in number.")
	}
	if row1.Symbol != stockList[1].Symbol ||
		fmt.Sprintf("%.2f", row1.MarketOpen) != fmt.Sprintf("%.2f", stockList[1].MarketOpen) ||
		fmt.Sprintf("%.2f", row1.MarketHigh) != fmt.Sprintf("%.2f", stockList[1].MarketHigh) ||
		fmt.Sprintf("%.2f", row1.MarketLow) != fmt.Sprintf("%.2f", stockList[1].MarketLow) ||
		row1.Volume != stockList[1].Volume ||
		fmt.Sprintf("%.2f", row1.MarketClose) != fmt.Sprintf("%.2f", stockList[1].MarketClose) ||
		row2.Symbol != stockList[0].Symbol ||
		fmt.Sprintf("%.2f", row2.MarketOpen) != fmt.Sprintf("%.2f", stockList[0].MarketOpen) ||
		fmt.Sprintf("%.2f", row2.MarketHigh) != fmt.Sprintf("%.2f", stockList[0].MarketHigh) ||
		fmt.Sprintf("%.2f", row2.MarketLow) != fmt.Sprintf("%.2f", stockList[0].MarketLow) ||
		row2.Volume != stockList[0].Volume ||
		fmt.Sprintf("%.2f", row2.MarketClose) != fmt.Sprintf("%.2f", stockList[0].MarketClose) {
		t.Error("Error: last number of available stock hists selected by symbol and before date are not matched.")
	}
}
