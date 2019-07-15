package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type TblStockHist struct {
}

type RowStockHist struct {
	Symbol      string
	Date        int
	MarketOpen  float32
	MarketHigh  float32
	MarketLow   float32
	MarketClose float32
	Volume      int
}

func (tbl *TblStockHist) SelectLastStockHist(symbol string) (RowStockHist, error) {
	hist := RowStockHist{}
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return hist, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Symbol, Date, MarketOpen, MarketHigh, MarketLow, MarketClose, Volume FROM " + TBL_STOCK_HIST_NAME + " WHERE Symbol='" + symbol + "' ORDER BY Date DESC LIMIT 1")
	if dbPrepErr != nil {
		return hist, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query()
	if dbQueryErr != nil {
		return hist, dbQueryErr
	}
	defer rows.Close()
	var sb string
	var date int
	var open float32
	var high float32
	var low float32
	var close float32
	var volume int
	for rows.Next() {
		if scanErr := rows.Scan(&sb, &date, &open, &high, &low, &close, &volume); scanErr != nil {
			fmt.Println(scanErr)
			return hist, scanErr
		}
	}
	hist.Symbol = sb
	hist.Date = date
	hist.MarketOpen = open
	hist.MarketHigh = high
	hist.MarketLow = low
	hist.MarketClose = close
	hist.Volume = volume
	return hist, nil
}

//Select # of stock hist (reversed array) before the date exclusively the function is called
func (tbl *TblStockHist) SelectLastStockHistByCountAndBeforeDate(symbol string, count int, beforeDate int) ([]RowStockHist, error) {
	histList := []RowStockHist{}
	hist := RowStockHist{}
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return histList, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Symbol, Date, MarketOpen, MarketHigh, MarketLow, MarketClose, Volume FROM " + TBL_STOCK_HIST_NAME + " WHERE Symbol='" + symbol + "' AND Date<" + strconv.Itoa(beforeDate) + " ORDER BY Date DESC LIMIT " + strconv.Itoa(count))
	if dbPrepErr != nil {
		return histList, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query()
	if dbQueryErr != nil {
		return histList, dbQueryErr
	}
	defer rows.Close()
	var sb string
	var date int
	var open float32
	var high float32
	var low float32
	var close float32
	var volume int
	for rows.Next() {
		if scanErr := rows.Scan(&sb, &date, &open, &high, &low, &close, &volume); scanErr != nil {
			fmt.Println(scanErr)
			return histList, scanErr
		}
		hist.Symbol = sb
		hist.Date = date
		hist.MarketOpen = open
		hist.MarketHigh = high
		hist.MarketLow = low
		hist.MarketClose = close
		hist.Volume = volume
		histList = append(histList, hist)
	}
	return histList, nil
}

func (tbl *TblStockHist) DropTableIfExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + TBL_STOCK_HIST_NAME)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl *TblStockHist) CreateTableIfNotExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string

	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		TBL_STOCK_HIST_NAME +
		" (" +
		"Symbol VARCHAR(10) NOT NULL," +
		"Date INT NOT NULL," +
		"MarketOpen FLOAT(10,2)," +
		"MarketHigh FLOAT(10,2)," +
		"MarketLow FLOAT(10,2)," +
		"MarketClose FLOAT(10,2)," +
		"Volume BIGINT," +
		"UpdatedTime TIMESTAMP," +
		"PRIMARY KEY (Symbol, Date)" +
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

func (tbl *TblStockHist) InsertOrUpdateStockData(stock YahooQuote, date int) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			TBL_STOCK_HIST_NAME+
			" WHERE Symbol=? AND Date=?)",
		stock.Symbol, date)
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		return dbQueryErr
	}
	if exist == 0 {
		if err := tbl.InsertStockData(stock, date); err != nil {
			return err
		}
	} else {
		if err := tbl.UpdateStockData(stock, date); err != nil {
			return err
		}
	}
	return nil
}

func (tbl *TblStockHist) InsertStockData(stock YahooQuote, date int) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		TBL_STOCK_HIST_NAME +
		" (" +
		"Symbol, " +
		"Date, " +
		"MarketOpen," +
		"MarketHigh," +
		"MarketLow," +
		"MarketClose," +
		"Volume," +
		"UpdatedTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		stock.Symbol,
		date,
		stock.RegularMarketOpen,
		stock.RegularMarketDayHigh,
		stock.RegularMarketDayLow,
		stock.RegularMarketPrice,
		stock.RegularMarketVolume,
		GetTime())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func (tbl *TblStockHist) UpdateStockData(stock YahooQuote, date int) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		TBL_STOCK_HIST_NAME +
		" SET " +
		"MarketOpen=?, " +
		"MarketHigh=?, " +
		"MarketLow=?, " +
		"MarketClose=?, " +
		"Volume=?, " +
		"UpdatedTime=? " +
		"WHERE Symbol=? AND Date=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		stock.RegularMarketOpen,
		stock.RegularMarketDayHigh,
		stock.RegularMarketDayLow,
		stock.RegularMarketPrice,
		stock.RegularMarketVolume,
		GetTime(),
		stock.Symbol,
		date)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}
