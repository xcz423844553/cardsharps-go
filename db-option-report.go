package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TblOptionReport struct {
}

func (tbl TblOptionReport) InsertOptionReport(
	symbol string,
	date int,
	expDate int64,
	expectedPriceVol float32,
	expectedPriceOI float32,
	expectedPriceVolCall float32,
	expectedPriceVolPut float32,
	expectedPriceOiCall float32,
	expectedPriceOiPut float32,
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
	fmt.Println(symbol, expDate, expectedPriceVol, expectedPriceOI, stockPrice)
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
		"Expiration," +
		"ExpectedPriceVol, " +
		"ExpectedPriceOI, " +
		"ExpectedPriceVolCall, " +
		"ExpectedPriceVolPut, " +
		"ExpectedPriceOiCall, " +
		"ExpectedPriceOiPut, " +
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
		"UpdatedTime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		symbol,
		date,
		expDate,
		expectedPriceVol,
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
		"Expiration BIGINT," +
		"ExpectedPriceVol FLOAT(10,2)," +
		"ExpectedPriceOI FLOAT(10,2)," +
		"ExpectedPriceVolCall FLOAT(10,2)," +
		"ExpectedPriceVolPut FLOAT(10,2)," +
		"ExpectedPriceOiCall FLOAT(10,2)," +
		"ExpectedPriceOiPut FLOAT(10,2)," +
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
		"PRIMARY KEY (Symbol, Expiration, Date, UpdatedTime)" +
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
