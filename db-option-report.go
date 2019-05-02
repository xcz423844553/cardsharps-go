package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TblOptionReport struct {
}

func (tbl TblOptionReport) GenerateReport(symbol string, expDate int64) error {
	var sumPriceBuy float32
	var sumPriceSell float32
	var expectedPriceVol float32
	var expectedPriceOI float32
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
	//1. Download all option data for the symbol, expDate, and today's date
	var dbPrepErr error
	var dbQueryErr error
	var scanErr error
	var stmt *sql.Stmt
	var optionRows *sql.Rows
	var stockRow *sql.Row
	stmt, dbPrepErr = db.Prepare("SELECT OptionType, Strike, Volume, PrevVolume, OpenInterest FROM " +
		TBL_OPTION_DATA_NAME +
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
		GetTimeInYYYYMMDD(),
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
			fmt.Println(scanErr)
			return scanErr
		}
		sumPriceBuy += float32(volume) * strike
		sumPriceSell += float32(openInterest) * strike
		if optionType == "C" {
			oiCall += openInterest
			volCall += volume
			deltaVolCall += volume - prevVolume
			deltaSumPriceVolCall += float32(volume-prevVolume) * strike
		} else if optionType == "P" {
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
		TBL_STOCK_DATA_NAME +
		" WHERE " +
		"Symbol = ? AND " +
		"Date = ?;")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	stockRow = stmt.QueryRow(
		symbol,
		GetTimeInYYYYMMDD())
	if scanErr = stockRow.Scan(&stockPrice); scanErr != nil {
		return scanErr
	}
	//3. Save the result into database
	fmt.Printf("%f %f\n", expectedPriceVol, expectedPriceOI)
	new(TblOptionReport).InsertOptionReport(symbol, expectedPriceVol,
		expectedPriceOI,
		deltaVolCall,
		deltaVolPut,
		deltaVolTol,
		deltaExpectedPriceVolCall,
		deltaExpectedPriceVolPut,
		deltaExpectedPriceVolTol, oiCall, oiPut, oiTotal, volCall, volPut,
		volTotal, stockPrice)
	return nil
}

func (tbl TblOptionReport) InsertOptionReport(
	symbol string,
	expectedPriceVol float32,
	expectedPriceOI float32,
	deltaVolCall int,
	deltaVolPut int,
	deltaVolTol int,
	deltaExpectedPriceVolCall float32,
	deltaExpectedPriceVolPut float32,
	deltaExpectedPriceVolTol float32,
	oiCall int,
	oiPut int,
	oiTotal int,
	volCall int,
	volPut int,
	volTotal int,
	stockPrice float32) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		TBL_OPTION_REPORT_NAME +
		" (" +
		"Symbol, " +
		"Date, " +
		"ExpectedPriceVol, " +
		"ExpectedPriceOI, " +
		"DeltaVolCall," +
		"DeltaVolPut," +
		"DeltaVolTol," +
		"DeltaExpectedPriceVolCall," +
		"DeltaExpectedPriceVolPut," +
		"DeltaExpectedPriceVolTol," +
		"OiCall, " +
		"OiPut, " +
		"OiTotal, " +
		"VolCall, " +
		"VolPut, " +
		"VolTotal, " +
		"StockPrice, " +
		"UpdatedTime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if dbPrepErr != nil {
		panic(dbPrepErr)
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		symbol,
		GetTimeInYYYYMMDD(),
		expectedPriceVol,
		expectedPriceOI,
		deltaVolCall,
		deltaVolPut,
		deltaVolTol,
		deltaExpectedPriceVolCall,
		deltaExpectedPriceVolPut,
		deltaExpectedPriceVolTol,
		oiCall,
		oiPut,
		oiTotal,
		volCall,
		volPut,
		volTotal,
		stockPrice,
		GetTime())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func (tbl TblOptionReport) DropTableIfExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + TBL_OPTION_REPORT_NAME)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl TblOptionReport) CreateTableIfNotExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string
	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		TBL_OPTION_REPORT_NAME +
		" (" +
		"Symbol VARCHAR(10) NOT NULL," +
		"Date INT NOT NULL," +
		"ExpectedPriceVol FLOAT(10,2)," +
		"ExpectedPriceOI FLOAT(10,2)," +
		"DeltaVolCall INT," +
		"DeltaVolPut INT," +
		"DeltaVolTol INT," +
		"DeltaExpectedPriceVolCall FLOAT(10,2)," +
		"DeltaExpectedPriceVolPut FLOAT(10,2)," +
		"DeltaExpectedPriceVolTol FLOAT(10,2)," +
		"OiCall INT," +
		"OiPut INT," +
		"OiTotal INT," +
		"VolCall INT," +
		"VolPut INT," +
		"VolTotal INT," +
		"StockPrice FLOAT(10,2)," +
		"UpdatedTime TIMESTAMP," +
		"PRIMARY KEY (Symbol, Date, UpdatedTime)" +
		")"
	if sqlStr != "" {
		_, tblCreateErr := db.Exec(sqlStr)
		if tblCreateErr != nil {
			return tblCreateErr
		}
		return nil
	}
	return errors.New("failed to find preset table name")
}
