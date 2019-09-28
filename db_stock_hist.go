package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

//IStockHist is the struct of input to db_stock_hist
type IStockHist interface {
	GetSymbol() string
	GetMarketOpen() float32
	GetMarketHigh() float32
	GetMarketLow() float32
	GetMarketClose() float32
	GetVolume() int64
}

//DaoStockHist is a struct to manipulate db_stock_hist
type DaoStockHist struct {
}

//RowStockHist is a struct representing the row of db_stock_hist
type RowStockHist2 struct {
	Symbol      string
	Date        int64
	MarketOpen  float32
	MarketHigh  float32
	MarketLow   float32
	MarketClose float32
	Volume      int64
}

//GetTableName returns table name
func (dao *DaoStockHist) getTableName() string {
	if TestMode {
		return TblNameStockHist + "_test"
	}
	return TblNameStockHist
}

//DropTableIfExist drops table if the table exists
func (dao *DaoStockHist) DropTableIfExist() error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + dao.getTableName())
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

//CreateTableIfNotExist creates table if the table does not exist
func (dao *DaoStockHist) CreateTableIfNotExist() error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	sqlStr := "CREATE TABLE IF NOT EXISTS " +
		dao.getTableName() +
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
	_, tblCreateErr := db.Exec(sqlStr)
	if tblCreateErr != nil {
		return tblCreateErr
	}
	return nil
}

//InsertOrUpdateStockHist inserts or updates the stock hist
func (dao *DaoStockHist) InsertOrUpdateStockHist(stock IStockHist, date int64) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			dao.getTableName()+
			" WHERE Symbol=? AND Date=?)",
		stock.GetSymbol(), date)
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		return dbQueryErr
	}
	if exist == 0 {
		if err := dao.InsertStockHist(stock, date); err != nil {
			return err
		}
	} else {
		if err := dao.UpdateStockHist(stock, date); err != nil {
			return err
		}
	}
	return nil
}

//InsertStockHist inserts the stock hist
func (dao *DaoStockHist) InsertStockHist(stock IStockHist, date int64) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		dao.getTableName() +
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
		stock.GetSymbol(),
		date,
		stock.GetMarketOpen(),
		stock.GetMarketHigh(),
		stock.GetMarketLow(),
		stock.GetMarketClose(),
		stock.GetVolume(),
		GetTime())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

//UpdateStockHist updates the stock hist
func (dao *DaoStockHist) UpdateStockHist(stock IStockHist, date int64) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		dao.getTableName() +
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
		stock.GetMarketOpen(),
		stock.GetMarketHigh(),
		stock.GetMarketLow(),
		stock.GetMarketClose(),
		stock.GetVolume(),
		GetTime(),
		stock.GetSymbol(),
		date)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

//SelectLastStockHist selects the stock hist of last available date
func (dao *DaoStockHist) SelectLastStockHist(symbol string) (RowStockHist2, error) {
	hist := RowStockHist2{}
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return hist, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Symbol, Date, MarketOpen, MarketHigh, MarketLow, MarketClose, Volume FROM " +
		dao.getTableName() +
		" WHERE Symbol='" +
		symbol +
		"' ORDER BY Date DESC LIMIT 1")
	if dbPrepErr != nil {
		return hist, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query()
	if dbQueryErr != nil {
		return hist, dbQueryErr
	}
	defer rows.Close()
	var sym string
	var date int64
	var open float32
	var high float32
	var low float32
	var close float32
	var volume int64
	for rows.Next() {
		if scanErr := rows.Scan(&sym, &date, &open, &high, &low, &close, &volume); scanErr != nil {
			fmt.Println(scanErr)
			return hist, scanErr
		}
	}
	hist.Symbol = sym
	hist.Date = date
	hist.MarketOpen = open
	hist.MarketHigh = high
	hist.MarketLow = low
	hist.MarketClose = close
	hist.Volume = volume
	return hist, nil
}

//SelectLastNumberStockHistBeforeDate selects # rows of stock hist before a certain date exclusively; First element of the return is the latest date
func (dao *DaoStockHist) SelectLastNumberStockHistBeforeDate(symbol string, count int64, beforeDate int64) ([]RowStockHist2, error) {
	histList := []RowStockHist2{}
	hist := RowStockHist2{}
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return nil, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Symbol, Date, MarketOpen, MarketHigh, MarketLow, MarketClose, Volume FROM " +
		dao.getTableName() +
		" WHERE Symbol='" +
		symbol +
		"' AND Date<" +
		strconv.FormatInt(beforeDate, 10) +
		" ORDER BY Date DESC LIMIT " +
		strconv.FormatInt(count, 10))
	if dbPrepErr != nil {
		return histList, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query()
	if dbQueryErr != nil {
		return histList, dbQueryErr
	}
	defer rows.Close()
	var sym string
	var date int64
	var open float32
	var high float32
	var low float32
	var close float32
	var volume int64
	for rows.Next() {
		if scanErr := rows.Scan(&sym, &date, &open, &high, &low, &close, &volume); scanErr != nil {
			fmt.Println(scanErr)
			return histList, scanErr
		}
		hist.Symbol = sym
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
