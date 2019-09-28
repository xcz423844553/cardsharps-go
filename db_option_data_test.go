package main

import (
	"fmt"
	"testing"
)

func TestDbOptionData(t *testing.T) {
	dao := new(DaoOptionData)
	if dbDropErr := dao.DropTableIfExist(); dbDropErr != nil {
		t.Error(dbDropErr.Error())
	}
	if dbCreateErr := dao.CreateTableIfNotExist(); dbCreateErr != nil {
		t.Error(dbCreateErr.Error())
	}
	contractSymbol := "SPY190918C00287000"
	symbol := "SPY"
	date := GetTimeInYYYYMMDD64()
	row := RowOptionData{
		ContractSymbol:    contractSymbol,
		Date:              date,
		Symbol:            symbol,
		OptionType:        "C",
		Strike:            287,
		LastPrice:         14.24,
		Volume:            10,
		OpenInterest:      0,
		ImpliedVolatility: 0.4062559375,
		PercentChange:     0,
		Bid:               13.29,
		Ask:               13.47,
		Expiration:        20190918,
		LastTradeDate:     20190913,
	}
	input := ApiYahooOption{
		ContractSymbol:    contractSymbol,
		Strike:            287,
		LastPrice:         14.24,
		PercentChange:     0,
		Volume:            10,
		OpenInterest:      0,
		Bid:               13.29,
		Ask:               13.47,
		Expiration:        1568764800,
		ImpliedVolatility: 0.4062559375,
		LastTradeDate:     1568398749,
	}
	if dbInsertErr := dao.InsertOrUpdateOptionData(&input); dbInsertErr != nil {
		t.Error(dbInsertErr.Error())
	}
	if dbUpdateErr := dao.InsertOrUpdateOptionData(&input); dbUpdateErr != nil {
		t.Error(dbUpdateErr.Error())
	}

	expArr, dbExpErr := dao.SelectExpirationBySymbolAndDate(symbol, date)
	if dbExpErr != nil {
		t.Error(dbExpErr.Error())
	}
	if len(expArr) != 1 {
		t.Error("Error: # of expiration date rows not matched.")
	}
	if expArr[0] != row.Expiration {
		t.Error("Error: Expiration date is not matched.")
	}
	vol, dbVolErr := dao.SelectOptionDataVolumeByContractSymbolAndDate(contractSymbol, date)
	if dbVolErr != nil {
		t.Error(dbVolErr.Error())
	}
	if vol != row.Volume {
		t.Error("Error: Volume is not matched.")
	}
	cs, dbCsErr := dao.SelectContractSymbolBySymbolAndDate(symbol, date)
	if dbCsErr != nil {
		t.Error(dbCsErr)
	}
	if len(cs) != 1 {
		t.Error("Error: # of contract symbol is not matched.")
	}
	if cs[0] != contractSymbol {
		t.Error("Error: contract symbol is not matched.")
	}
	odRow, dbOdRowErr := dao.SelectOptionDataByContractSymbol(contractSymbol)
	if dbOdRowErr != nil {
		t.Error(dbOdRowErr.Error())
	}
	if row.ContractSymbol != odRow[0].ContractSymbol ||
		row.Date != odRow[0].Date ||
		row.Symbol != odRow[0].Symbol ||
		row.OptionType != odRow[0].OptionType ||
		row.Strike != odRow[0].Strike ||
		row.LastPrice != odRow[0].LastPrice ||
		row.Volume != odRow[0].Volume ||
		row.OpenInterest != odRow[0].OpenInterest ||
		fmt.Sprintf("%.2f", row.ImpliedVolatility) != fmt.Sprintf("%.2f", odRow[0].ImpliedVolatility) ||
		row.PercentChange != odRow[0].PercentChange ||
		row.Bid != odRow[0].Bid ||
		row.Ask != odRow[0].Ask ||
		row.Expiration != odRow[0].Expiration ||
		row.LastTradeDate != odRow[0].LastTradeDate {
		t.Error("Error: option data selected by contract symbol is not matched.")
	}
	expDate, callVol, putVol, callOi, putOi, dbAllErr := dao.SelectExpirationVolumeOpenInterestBySymbolAndDate(symbol, date)
	if dbAllErr != nil {
		t.Error(dbAllErr.Error())
	}
	if len(expDate) != 1 {
		t.Error("Error: # of expiration date rows not matched.")
	}
	if expDate[0] != row.Expiration {
		t.Error("Error: Expiration date is not matched.")
	}
	if len(callVol) != 1 {
		t.Error("Error: # of elements in call volume map is not matched.")
	}
	if callVol[expDate[0]] != row.Volume {
		t.Error("Error: Volume in call volume map is not matched.")
	}
	if len(callOi) != 1 {
		t.Error("Error: # of elements in call open interest map is not matched.")
	}
	if callOi[expDate[0]] != row.OpenInterest {
		t.Error("Error: Open interest in call open interest map is not matched.")
	}
	if len(putVol) != 0 {
		t.Error("Error: # of elements in put volume map is not matched.")
	}
	if len(putOi) != 0 {
		t.Error("Error: # of elements in put open interest map is not matched.")
	}
}
