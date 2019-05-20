package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Orbit struct {
}

func (orbit Orbit) runOptionReportForAllSymbol(date int) {
	tblSymbol := new(TblSymbol)
	tblLogError := new(TblLogError)
	//run non-etf symbols
	symbols, symbolSelectErr := tblSymbol.SelectSymbolByFilter()
	if symbolSelectErr != nil {
		tblLogError.InsertLogError(LOGTYPE_DB_SYMBOL, symbolSelectErr.Error())
	}
	fmt.Println(symbols)
	for _, symbol := range symbols {
		orbit.runOptionReportForSymbol(symbol, date, false)
	}
	return
}

func (orbit Orbit) runOptionReportForSymbol(symbol string, date int, isEtf bool) {
	tblOptionData := new(TblOptionData)
	tblLogError := new(TblLogError)
	expDates, selectErr := tblOptionData.SelectExpirationBySymbolAndDate(symbol, date, isEtf)
	if selectErr != nil {
		tblLogError.InsertLogError(LOGTYPE_ORBIT, selectErr.Error())
		PrintMsgInConsole(MSGERROR, LOGTYPE_ORBIT, "Failed to select expiration dates for "+symbol+" on "+strconv.Itoa(date))
		return
	}
	for _, expDate := range expDates {
		createErr := orbit.createOptionReport(symbol, expDate, date, isEtf)
		if createErr != nil {
			tblLogError.InsertLogError(LOGTYPE_ORBIT, createErr.Error())
			PrintMsgInConsole(MSGERROR, LOGTYPE_ORBIT, "Failed to create option report for "+symbol+" @ "+strconv.FormatInt(expDate, 10)+" on "+strconv.Itoa(date))
			continue
		}
	}
}

func (orbit Orbit) createOptionReport(symbol string, expDate int64, date int, isEtf bool) error {
	var tblOptionName string
	var tblStockName string
	if isEtf {
		tblOptionName = TBL_OPTION_DATA_ETF_NAME
		tblStockName = TBL_STOCK_DATA_ETF_NAME
	} else {
		tblOptionName = TBL_OPTION_DATA_NAME
		tblStockName = TBL_STOCK_DATA_NAME
	}
	tblOptionReport := new(TblOptionReport)
	var sumPriceBuy float32
	var sumPriceSell float32
	var expectedPriceVol float32
	var expectedPriceOI float32
	var sumPriceVolCall float32
	var sumPriceVolPut float32
	var sumPriceOiCall float32
	var sumPriceOiPut float32
	var expectedPriceVolCall float32
	var expectedPriceVolPut float32
	var expectedPriceOiCall float32
	var expectedPriceOiPut float32
	var oiCall int
	var oiPut int
	var oiTotal int
	var volCall int
	var volPut int
	var volTotal int
	var stockPrice float32
	var deltaVolCall int
	var deltaVolPut int
	var deltaVolTol int
	var deltaSumPriceVolCall float32
	var deltaSumPriceVolPut float32
	var deltaSumPriceVolTol float32
	var deltaExpectedPriceVolCall float32
	var deltaExpectedPriceVolPut float32
	var deltaExpectedPriceVolTol float32
	//Connect to database
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	//1. Download all option data for the symbol, expDate, and date
	var dbPrepErr error
	var dbQueryErr error
	var scanErr error
	var stmt *sql.Stmt
	var optionRows *sql.Rows
	var stockRow *sql.Row
	stmt, dbPrepErr = db.Prepare("SELECT OptionType, Strike, Volume, PrevVolume, OpenInterest FROM " +
		tblOptionName +
		" WHERE " +
		"Symbol = ? AND " +
		"Date = ? AND " +
		"Expiration = ?;")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	optionRows, dbQueryErr = stmt.Query(
		symbol,
		date,
		expDate)
	defer optionRows.Close()
	if dbQueryErr != nil {
		return dbQueryErr
	}
	for optionRows.Next() {
		var optionType string
		var strike float32
		var volume int
		var prevVolume int
		var openInterest int
		if scanErr = optionRows.Scan(&optionType, &strike, &volume, &prevVolume, &openInterest); scanErr != nil {
			return scanErr
		}
		sumPriceBuy += float32(volume) * strike
		sumPriceSell += float32(openInterest) * strike
		if optionType == "C" {
			sumPriceVolCall += float32(volume) * strike
			sumPriceOiCall += float32(openInterest) * strike
			oiCall += openInterest
			volCall += volume
			deltaVolCall += volume - prevVolume
			deltaSumPriceVolCall += float32(volume-prevVolume) * strike
		} else if optionType == "P" {
			sumPriceVolPut += float32(volume) * strike
			sumPriceOiPut += float32(openInterest) * strike
			oiPut += openInterest
			volPut += volume
			deltaVolPut += volume - prevVolume
			deltaSumPriceVolPut += float32(volume-prevVolume) * strike
		}
		oiTotal += openInterest
		volTotal += volume
		deltaVolTol += volume - prevVolume
		deltaSumPriceVolTol += float32(volume-prevVolume) * strike
	}
	if volTotal == 0 {
		expectedPriceVol = 0
	} else {
		expectedPriceVol = sumPriceBuy / float32(volTotal)
	}
	if oiTotal == 0 {
		expectedPriceOI = 0
	} else {
		expectedPriceOI = sumPriceSell / float32(oiTotal)
	}
	if volCall == 0 {
		expectedPriceVolCall = 0
	} else {
		expectedPriceVolCall = sumPriceVolCall / float32(volCall)
	}
	if oiCall == 0 {
		expectedPriceOiCall = 0
	} else {
		expectedPriceOiCall = sumPriceOiCall / float32(oiCall)
	}
	if volPut == 0 {
		expectedPriceVolPut = 0
	} else {
		expectedPriceVolPut = sumPriceVolPut / float32(volPut)
	}
	if oiPut == 0 {
		expectedPriceOiPut = 0
	} else {
		expectedPriceOiPut = sumPriceOiPut / float32(oiPut)
	}
	if deltaVolCall == 0 {
		deltaExpectedPriceVolCall = 0
	} else {
		deltaExpectedPriceVolCall = deltaSumPriceVolCall / float32(deltaVolCall)
	}
	if deltaVolPut == 0 {
		deltaExpectedPriceVolPut = 0
	} else {
		deltaExpectedPriceVolPut = deltaSumPriceVolPut / float32(deltaVolPut)
	}
	if deltaVolTol == 0 {
		deltaExpectedPriceVolTol = 0
	} else {
		deltaExpectedPriceVolTol = deltaSumPriceVolTol / float32(deltaVolTol)
	}
	//2. Download stock data for the symbol and today's date
	stmt, dbPrepErr = db.Prepare("SELECT RegularMarketPrice FROM " +
		tblStockName +
		" WHERE " +
		"Symbol = ? AND " +
		"Date = ?;")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	stockRow = stmt.QueryRow(
		symbol,
		date)
	if scanErr = stockRow.Scan(&stockPrice); scanErr != nil {
		return scanErr
	}
	//3. Save the result into database
	tblOptionReport.InsertOptionReport(symbol, date, expDate, expectedPriceVol,
		expectedPriceOI,
		expectedPriceVolCall,
		expectedPriceVolPut,
		expectedPriceOiCall,
		expectedPriceOiPut,
		deltaVolCall,
		deltaVolPut,
		deltaVolTol,
		deltaExpectedPriceVolCall,
		deltaExpectedPriceVolPut,
		deltaExpectedPriceVolTol, oiCall, oiPut, oiTotal, volCall, volPut,
		volTotal, stockPrice)
	return nil
}
