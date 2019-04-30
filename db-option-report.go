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
	var expectedPriceBuy float32
	var expectedPriceSell float32
	var oiCall int
	var oiPut int
	var oiTotal int
	var volCall int
	var volPut int
	var volTotal int
	var stockPrice float32
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
	stmt, dbPrepErr = db.Prepare("SELECT OptionType, Strike, Volume, OpenInterest FROM " +
		TBL_OPTION_DATA_NAME +
		" WHERE " +
		"Symbol = ? AND " +
		"Date = ? AND " +
		"Expiration = ?;")
	if dbPrepErr != nil {
		panic(dbPrepErr)
		return dbPrepErr
	}
	defer stmt.Close()
	optionRows, dbQueryErr = stmt.Query(
		symbol,
		GetTimeInYYYYMMDD(),
		expDate)
	defer optionRows.Close()
	if dbQueryErr != nil {
		panic(dbQueryErr)
		return dbQueryErr
	}
	for optionRows.Next() {
		var optionType string
		var strike float32
		var volume int
		var openInterest int
		if scanErr = optionRows.Scan(&optionType, &strike, &volume, &openInterest); scanErr != nil {
			fmt.Println(scanErr)
			return scanErr
		}
		sumPriceBuy += float32(volume) * strike
		sumPriceSell += float32(openInterest) * strike
		if optionType == "C" {
			oiCall += openInterest
			volCall += volume
		} else if optionType == "P" {
			oiPut += openInterest
			volPut += volume
		}
		oiTotal += openInterest
		volTotal += volume
	}
	if volTotal == 0 {
		expectedPriceBuy = 0
	} else {
		expectedPriceBuy = sumPriceBuy / float32(volTotal)
	}
	if oiTotal == 0 {
		expectedPriceSell = 0
	} else {
		expectedPriceSell = sumPriceSell / float32(oiTotal)
	}
	//Download stock data for the symbol and today's date
	stmt, dbPrepErr = db.Prepare("SELECT RegularMarketPrice FROM " +
		TBL_STOCK_DATA_NAME +
		" WHERE " +
		"Symbol = ? AND " +
		"Date = ?;")
	if dbPrepErr != nil {
		panic(dbPrepErr)
		return dbPrepErr
	}
	defer stmt.Close()
	stockRow = stmt.QueryRow(
		symbol,
		GetTimeInYYYYMMDD())
	if scanErr = stockRow.Scan(&stockPrice); scanErr != nil {
		panic(scanErr)
		return scanErr
	}
	//Save the result into database
	fmt.Printf("%f %f\n", expectedPriceBuy, expectedPriceSell)
	new(TblOptionReport).InsertOptionReport(symbol, expectedPriceBuy,
		expectedPriceSell, oiCall, oiPut, oiTotal, volCall, volPut,
		volTotal, stockPrice)
	return nil
}

func (tbl TblOptionReport) InsertOptionReport(
	symbol string,
	expectedPriceBuy float32,
	expectedPriceSell float32,
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
		"ExpectedPriceBuy, " +
		"ExpectedPriceSell, " +
		"OiCall, " +
		"OiPut, " +
		"OiTotal, " +
		"VolCall, " +
		"VolPut, " +
		"VolTotal, " +
		"StockPrice, " +
		"UpdatedTime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")
	if dbPrepErr != nil {
		panic(dbPrepErr)
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		symbol,
		GetTimeInYYYYMMDD(),
		expectedPriceBuy,
		expectedPriceSell,
		oiCall,
		oiPut,
		oiTotal,
		volCall,
		volPut,
		volTotal,
		stockPrice,
		GetTime())
	if dbExecErr != nil {
		panic(dbExecErr)
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
		"ExpectedPriceBuy FLOAT(10,2)," +
		"ExpectedPriceSell FLOAT(10,2)," +
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
