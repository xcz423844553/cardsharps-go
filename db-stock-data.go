package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type TblStockData struct {
}

func (tbl TblStockData) DropTableIfExist() error {
	if err := tbl.DropTableIfExistForTblName(TBL_STOCK_DATA_NAME); err != nil {
		return err
	}
	return tbl.DropTableIfExistForTblName(TBL_STOCK_DATA_ETF_NAME)
}

func (tbl TblStockData) DropTableIfExistForTblName(tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + tblName)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl TblStockData) CreateTableIfNotExist() error {
	if err := tbl.CreateTableIfNotExistForTblName(TBL_STOCK_DATA_NAME); err != nil {
		return err
	}
	return tbl.CreateTableIfNotExistForTblName(TBL_STOCK_DATA_ETF_NAME)
}

func (tbl TblStockData) CreateTableIfNotExistForTblName(tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string

	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		tblName +
		" (" +
		"Symbol VARCHAR(10) NOT NULL," +
		"Date INT NOT NULL," +
		"RegularMarketChange FLOAT(10,2)," +
		"RegularMarketOpen FLOAT(10,2)," +
		"RegularMarketDayHigh FLOAT(10,2)," +
		"RegularMarketDayLow FLOAT(10,2)," +
		"RegularMarketVolume BIGINT," +
		"RegularMarketChangePercent FLOAT(10,2)," +
		"RegularMarketPreviousClose FLOAT(10,2)," +
		"RegularMarketPrice FLOAT(10,2)," +
		"RegularMarketTime BIGINT," +
		"EarningsTimestamp BIGINT," +
		"FiftyDayAverage FLOAT(10,2)," +
		"FiftyDayAverageChange FLOAT(10,2)," +
		"FiftyDayAverageChangePercent FLOAT(10,2)," +
		"TwoHundredDayAverage FLOAT(10,2)," +
		"TwoHundredDayAverageChange FLOAT(10,2)," +
		"TwoHundredDayAverageChangePercent FLOAT(10,2)," +
		"PostMarketChangePercent FLOAT(10,2)," +
		"PostMarketTime BIGINT," +
		"PostMarketPrice FLOAT(10,2)," +
		"PostMarketChange FLOAT(10,2)," +
		"Bid FLOAT(10,2)," +
		"Ask FLOAT(10,2)," +
		"BidSize BIGINT," +
		"AskSize BIGINT," +
		"AverageDailyVolume3Month BIGINT," +
		"AverageDailyVolume10Day BIGINT," +
		"FiftyTwoWeekLowChange FLOAT(10,2)," +
		"FiftyTwoWeekLowChangePercent FLOAT(10,2)," +
		"FiftyTwoWeekHighChange FLOAT(10,2)," +
		"FiftyTwoWeekHighChangePercent FLOAT(10,2)," +
		"FiftyTwoWeekLow FLOAT(10,2)," +
		"FiftyTwoWeekHigh FLOAT(10,2)," +
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

func (tbl TblStockData) InsertOrUpdateStockData(stock YahooQuote, isEtf bool) error {
	if isEtf {
		return tbl.InsertOrUpdateStockDataToTbl(stock, TBL_STOCK_DATA_ETF_NAME)
	} else {
		return tbl.InsertOrUpdateStockDataToTbl(stock, TBL_STOCK_DATA_NAME)
	}
}

func (tbl TblStockData) InsertOrUpdateStockDataToTbl(stock YahooQuote, tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			tblName+
			" WHERE Symbol=? AND Date=?)",
		stock.Symbol, GetTimeInYYYYMMDD())
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		return dbQueryErr
	}
	if exist == 0 {
		if err := tbl.InsertStockData(stock, tblName); err != nil {
			return err
		}
	} else {
		if err := tbl.UpdateStockData(stock, tblName); err != nil {
			return err
		}
	}
	return nil
}

func (tbl TblStockData) InsertStockData(stock YahooQuote, tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		tblName +
		" (" +
		"Symbol, " +
		"Date, " +
		"RegularMarketChange, " +
		"RegularMarketOpen, " +
		"RegularMarketDayHigh, " +
		"RegularMarketDayLow, " +
		"RegularMarketVolume, " +
		"RegularMarketChangePercent, " +
		"RegularMarketPreviousClose, " +
		"RegularMarketPrice, " +
		"RegularMarketTime, " +
		"EarningsTimestamp, " +
		"FiftyDayAverage, " +
		"FiftyDayAverageChange, " +
		"FiftyDayAverageChangePercent, " +
		"TwoHundredDayAverage, " +
		"TwoHundredDayAverageChange, " +
		"TwoHundredDayAverageChangePercent, " +
		"PostMarketChangePercent, " +
		"PostMarketTime, " +
		"PostMarketPrice, " +
		"PostMarketChange, " +
		"Bid, " +
		"Ask, " +
		"BidSize, " +
		"AskSize, " +
		"AverageDailyVolume3Month, " +
		"AverageDailyVolume10Day, " +
		"FiftyTwoWeekLowChange, " +
		"FiftyTwoWeekLowChangePercent, " +
		"FiftyTwoWeekHighChange, " +
		"FiftyTwoWeekHighChangePercent, " +
		"FiftyTwoWeekLow, " +
		"FiftyTwoWeekHigh, " +
		"UpdatedTime) VALUES (?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		stock.Symbol,
		GetTimeInYYYYMMDD(),
		stock.RegularMarketChange,
		stock.RegularMarketOpen,
		stock.RegularMarketDayHigh,
		stock.RegularMarketDayLow,
		stock.RegularMarketVolume,
		stock.RegularMarketChangePercent,
		stock.RegularMarketPreviousClose,
		stock.RegularMarketPrice,
		stock.RegularMarketTime,
		stock.EarningsTimestamp,
		stock.FiftyDayAverage,
		stock.FiftyDayAverageChange,
		stock.FiftyDayAverageChangePercent,
		stock.TwoHundredDayAverage,
		stock.TwoHundredDayAverageChange,
		stock.TwoHundredDayAverageChangePercent,
		stock.PostMarketChangePercent,
		stock.PostMarketTime,
		stock.PostMarketPrice,
		stock.PostMarketChange,
		stock.Bid,
		stock.Ask,
		stock.BidSize,
		stock.AskSize,
		stock.AverageDailyVolume3Month,
		stock.AverageDailyVolume10Day,
		stock.FiftyTwoWeekLowChange,
		stock.FiftyTwoWeekLowChangePercent,
		stock.FiftyTwoWeekHighChange,
		stock.FiftyTwoWeekHighChangePercent,
		stock.FiftyTwoWeekLow,
		stock.FiftyTwoWeekHigh,
		GetTime())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func (tbl TblStockData) UpdateStockData(stock YahooQuote, tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		tblName +
		" SET " +
		"RegularMarketChange=?, " +
		"RegularMarketOpen=?, " +
		"RegularMarketDayHigh=?, " +
		"RegularMarketDayLow=?, " +
		"RegularMarketVolume=?, " +
		"RegularMarketChangePercent=?, " +
		"RegularMarketPreviousClose=?, " +
		"RegularMarketPrice=?, " +
		"RegularMarketTime=?, " +
		"EarningsTimestamp=?, " +
		"FiftyDayAverage=?, " +
		"FiftyDayAverageChange=?, " +
		"FiftyDayAverageChangePercent=?, " +
		"TwoHundredDayAverage=?, " +
		"TwoHundredDayAverageChange=?, " +
		"TwoHundredDayAverageChangePercent=?, " +
		"PostMarketChangePercent=?, " +
		"PostMarketTime=?, " +
		"PostMarketPrice=?, " +
		"PostMarketChange=?, " +
		"Bid=?, " +
		"Ask=?, " +
		"BidSize=?, " +
		"AskSize=?, " +
		"AverageDailyVolume3Month=?, " +
		"AverageDailyVolume10Day=?, " +
		"FiftyTwoWeekLowChange=?, " +
		"FiftyTwoWeekLowChangePercent=?, " +
		"FiftyTwoWeekHighChange=?, " +
		"FiftyTwoWeekHighChangePercent=?, " +
		"FiftyTwoWeekLow=?, " +
		"FiftyTwoWeekHigh=?, " +
		"UpdatedTime=? " +
		"WHERE Symbol=? AND Date=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		stock.RegularMarketChange,
		stock.RegularMarketOpen,
		stock.RegularMarketDayHigh,
		stock.RegularMarketDayLow,
		stock.RegularMarketVolume,
		stock.RegularMarketChangePercent,
		stock.RegularMarketPreviousClose,
		stock.RegularMarketPrice,
		stock.RegularMarketTime,
		stock.EarningsTimestamp,
		stock.FiftyDayAverage,
		stock.FiftyDayAverageChange,
		stock.FiftyDayAverageChangePercent,
		stock.TwoHundredDayAverage,
		stock.TwoHundredDayAverageChange,
		stock.TwoHundredDayAverageChangePercent,
		stock.PostMarketChangePercent,
		stock.PostMarketTime,
		stock.PostMarketPrice,
		stock.PostMarketChange,
		stock.Bid,
		stock.Ask,
		stock.BidSize,
		stock.AskSize,
		stock.AverageDailyVolume3Month,
		stock.AverageDailyVolume10Day,
		stock.FiftyTwoWeekLowChange,
		stock.FiftyTwoWeekLowChangePercent,
		stock.FiftyTwoWeekHighChange,
		stock.FiftyTwoWeekHighChangePercent,
		stock.FiftyTwoWeekLow,
		stock.FiftyTwoWeekHigh,
		GetTime(),
		stock.Symbol,
		GetTimeInYYYYMMDD())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}
