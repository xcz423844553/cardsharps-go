package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//IStockData is the struct of input to db_stock_data
type IStockData interface {
	GetSymbol() string
	GetRegularMarketChange() float32
	GetRegularMarketOpen() float32
	GetRegularMarketDayHigh() float32
	GetRegularMarketDayLow() float32
	GetRegularMarketVolume() int64
	GetRegularMarketChangePercent() float32
	GetRegularMarketPreviousClose() float32
	GetRegularMarketPrice() float32
	GetRegularMarketTime() int64
	GetEarningsTimestamp() int64
	GetFiftyDayAverage() float32
	GetFiftyDayAverageChange() float32
	GetFiftyDayAverageChangePercent() float32
	GetTwoHundredDayAverage() float32
	GetTwoHundredDayAverageChange() float32
	GetTwoHundredDayAverageChangePercent() float32
	IsTradeable() bool
	GetMarketState() string
	GetPostMarketChangePercent() float32
	GetPostMarketTime() int64
	GetPostMarketPrice() float32
	GetPostMarketChange() float32
	GetBid() float32
	GetAsk() float32
	GetBidSize() int64
	GetAskSize() int64
	GetAverageDailyVolume3Month() int64
	GetAverageDailyVolume10Day() int64
	GetFiftyTwoWeekLowChange() float32
	GetFiftyTwoWeekLowChangePercent() float32
	GetFiftyTwoWeekHighChange() float32
	GetFiftyTwoWeekHighChangePercent() float32
	GetFiftyTwoWeekLow() float32
	GetFiftyTwoWeekHigh() float32
}

//DaoStockData is a struct to manipulate db_stock_data
type DaoStockData struct {
}

//RowStockData is a struct representing the row of db_stock_data
type RowStockData struct {
	Symbol                            string  `json:"symbol"`
	RegularMarketChange               float32 `json:"regularMarketChange"`
	RegularMarketOpen                 float32 `json:"regularMarketOpen"`
	RegularMarketDayHigh              float32 `json:"regularMarketDayHigh"`
	RegularMarketDayLow               float32 `json:"regularMarketDayLow"`
	RegularMarketVolume               int64   `json:"regularMarketVolume"`
	RegularMarketChangePercent        float32 `json:"regularMarketChangePercent"`
	RegularMarketPreviousClose        float32 `json:"regularMarketPreviousClose"`
	RegularMarketPrice                float32 `json:"regularMarketPrice"`
	RegularMarketTime                 int64   `json:"regularMarketTime"`
	EarningsTimestamp                 int64   `json:"earningsTimestamp"`
	FiftyDayAverage                   float32 `json:"fiftyDayAverage"`
	FiftyDayAverageChange             float32 `json:"fiftyDayAverageChange"`
	FiftyDayAverageChangePercent      float32 `json:"fiftyDayAverageChangePercent"`
	TwoHundredDayAverage              float32 `json:"twoHundredDayAverage"`
	TwoHundredDayAverageChange        float32 `json:"twoHundredDayAverageChange"`
	TwoHundredDayAverageChangePercent float32 `json:"twoHundredDayAverageChangePercent"`
	Tradeable                         bool    `json:"tradeable"`
	MarketState                       string  `json:"marketState"`
	PostMarketChangePercent           float32 `json:"postMarketChangePercent"`
	PostMarketTime                    int64   `json:"postMarketTime"`
	PostMarketPrice                   float32 `json:"postMarketPrice"`
	PostMarketChange                  float32 `json:"postMarketChange"`
	Bid                               float32 `json:"bid"`
	Ask                               float32 `json:"ask"`
	BidSize                           int64   `json:"bidSize"`
	AskSize                           int64   `json:"askSize"`
	AverageDailyVolume3Month          int64   `json:"averageDailyVolume3Month"`
	AverageDailyVolume10Day           int64   `json:"averageDailyVolume10Day"`
	FiftyTwoWeekLowChange             float32 `json:"fiftyTwoWeekLowChange"`
	FiftyTwoWeekLowChangePercent      float32 `json:"fiftyTwoWeekLowChangePercent"`
	FiftyTwoWeekHighChange            float32 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekHighChangePercent     float32 `json:"fiftyTwoWeekHighChangePercent"`
	FiftyTwoWeekLow                   float32 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh                  float32 `json:"fiftyTwoWeekHigh"`
}

//GetTableName returns table name
func (dao *DaoStockData) getTableName() string {
	if TestMode {
		return TblNameStockData + "_test"
	}
	return TblNameStockData
}

//DropTableIfExist drops table if the table exists
func (dao *DaoStockData) DropTableIfExist() error {
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
func (dao *DaoStockData) CreateTableIfNotExist() error {
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
	_, tblCreateErr := db.Exec(sqlStr)
	if tblCreateErr != nil {
		return tblCreateErr
	}
	return nil
}

//InsertOrUpdateStockData inserts or updates the stock data
func (dao *DaoStockData) InsertOrUpdateStockData(stock IStockData) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			dao.getTableName()+
			" WHERE Symbol=? AND Date=?)",
		stock.GetSymbol(), GetTimeInYYYYMMDD64())
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		return dbQueryErr
	}
	if exist == 0 {
		if err := dao.InsertStockData(stock); err != nil {
			return err
		}
	} else {
		if err := dao.UpdateStockData(stock); err != nil {
			return err
		}
	}
	return nil
}

//InsertStockData inserts the stock data
func (dao *DaoStockData) InsertStockData(stock IStockData) error {
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
		stock.GetSymbol(),
		GetTimeInYYYYMMDD64(),
		stock.GetRegularMarketChange(),
		stock.GetRegularMarketOpen(),
		stock.GetRegularMarketDayHigh(),
		stock.GetRegularMarketDayLow(),
		stock.GetRegularMarketVolume(),
		stock.GetRegularMarketChangePercent(),
		stock.GetRegularMarketPreviousClose(),
		stock.GetRegularMarketPrice(),
		stock.GetRegularMarketTime(),
		stock.GetEarningsTimestamp(),
		stock.GetFiftyDayAverage(),
		stock.GetFiftyDayAverageChange(),
		stock.GetFiftyDayAverageChangePercent(),
		stock.GetTwoHundredDayAverage(),
		stock.GetTwoHundredDayAverageChange(),
		stock.GetTwoHundredDayAverageChangePercent(),
		stock.GetPostMarketChangePercent(),
		stock.GetPostMarketTime(),
		stock.GetPostMarketPrice(),
		stock.GetPostMarketChange(),
		stock.GetBid(),
		stock.GetAsk(),
		stock.GetBidSize(),
		stock.GetAskSize(),
		stock.GetAverageDailyVolume3Month(),
		stock.GetAverageDailyVolume10Day(),
		stock.GetFiftyTwoWeekLowChange(),
		stock.GetFiftyTwoWeekLowChangePercent(),
		stock.GetFiftyTwoWeekHighChange(),
		stock.GetFiftyTwoWeekHighChangePercent(),
		stock.GetFiftyTwoWeekLow(),
		stock.GetFiftyTwoWeekHigh(),
		GetTime())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

//UpdateStockData updates the stock data
func (dao *DaoStockData) UpdateStockData(stock IStockData) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		dao.getTableName() +
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
		stock.GetRegularMarketChange(),
		stock.GetRegularMarketOpen(),
		stock.GetRegularMarketDayHigh(),
		stock.GetRegularMarketDayLow(),
		stock.GetRegularMarketVolume(),
		stock.GetRegularMarketChangePercent(),
		stock.GetRegularMarketPreviousClose(),
		stock.GetRegularMarketPrice(),
		stock.GetRegularMarketTime(),
		stock.GetEarningsTimestamp(),
		stock.GetFiftyDayAverage(),
		stock.GetFiftyDayAverageChange(),
		stock.GetFiftyDayAverageChangePercent(),
		stock.GetTwoHundredDayAverage(),
		stock.GetTwoHundredDayAverageChange(),
		stock.GetTwoHundredDayAverageChangePercent(),
		stock.GetPostMarketChangePercent(),
		stock.GetPostMarketTime(),
		stock.GetPostMarketPrice(),
		stock.GetPostMarketChange(),
		stock.GetBid(),
		stock.GetAsk(),
		stock.GetBidSize(),
		stock.GetAskSize(),
		stock.GetAverageDailyVolume3Month(),
		stock.GetAverageDailyVolume10Day(),
		stock.GetFiftyTwoWeekLowChange(),
		stock.GetFiftyTwoWeekLowChangePercent(),
		stock.GetFiftyTwoWeekHighChange(),
		stock.GetFiftyTwoWeekHighChangePercent(),
		stock.GetFiftyTwoWeekLow(),
		stock.GetFiftyTwoWeekHigh(),
		GetTime(),
		stock.GetSymbol(),
		GetTimeInYYYYMMDD64())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

//SelectStockDataBySymbolAndDate returns an array of expiration date, a map of call volume, a map of put volume, a map of call open interest, and a map of put open interest
func (dao *DaoStockData) SelectStockDataBySymbolAndDate(sym string, date int64) (RowStockData, error) {
	var stock RowStockData
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return stock, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Symbol, RegularMarketChange, RegularMarketOpen, RegularMarketDayHigh," +
		"RegularMarketDayLow, RegularMarketVolume, RegularMarketChangePercent, RegularMarketPreviousClose, RegularMarketPrice, RegularMarketTime," +
		"EarningsTimestamp, FiftyDayAverage, FiftyDayAverageChange, FiftyDayAverageChangePercent, TwoHundredDayAverage, TwoHundredDayAverageChange," +
		"TwoHundredDayAverageChangePercent, PostMarketChangePercent, PostMarketTime, PostMarketPrice, PostMarketChange,	Bid, Ask, BidSize, AskSize," +
		"AverageDailyVolume3Month, AverageDailyVolume10Day, FiftyTwoWeekLowChange, FiftyTwoWeekLowChangePercent, FiftyTwoWeekHighChange, FiftyTwoWeekHighChangePercent," +
		"FiftyTwoWeekLow, FiftyTwoWeekHigh FROM " + dao.getTableName() + " WHERE " + "Symbol=? AND Date=?")
	if dbPrepErr != nil {
		return stock, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query(sym, date)
	if dbQueryErr != nil {
		return stock, dbQueryErr
	}
	defer rows.Close()
	var symbol string
	var regularMarketChange float32
	var regularMarketOpen float32
	var regularMarketDayHigh float32
	var regularMarketDayLow float32
	var regularMarketVolume int64
	var regularMarketChangePercent float32
	var regularMarketPreviousClose float32
	var regularMarketPrice float32
	var regularMarketTime int64
	var earningsTimestamp int64
	var fiftyDayAverage float32
	var fiftyDayAverageChange float32
	var fiftyDayAverageChangePercent float32
	var twoHundredDayAverage float32
	var twoHundredDayAverageChange float32
	var twoHundredDayAverageChangePercent float32
	var postMarketChangePercent float32
	var postMarketTime int64
	var postMarketPrice float32
	var postMarketChange float32
	var bid float32
	var ask float32
	var bidSize int64
	var askSize int64
	var averageDailyVolume3Month int64
	var averageDailyVolume10Day int64
	var fiftyTwoWeekLowChange float32
	var fiftyTwoWeekLowChangePercent float32
	var fiftyTwoWeekHighChange float32
	var fiftyTwoWeekHighChangePercent float32
	var fiftyTwoWeekLow float32
	var fiftyTwoWeekHigh float32
	for rows.Next() {
		if scanErr := rows.Scan(&symbol,
			&regularMarketChange,
			&regularMarketOpen,
			&regularMarketDayHigh,
			&regularMarketDayLow,
			&regularMarketVolume,
			&regularMarketChangePercent,
			&regularMarketPreviousClose,
			&regularMarketPrice,
			&regularMarketTime,
			&earningsTimestamp,
			&fiftyDayAverage,
			&fiftyDayAverageChange,
			&fiftyDayAverageChangePercent,
			&twoHundredDayAverage,
			&twoHundredDayAverageChange,
			&twoHundredDayAverageChangePercent,
			&postMarketChangePercent,
			&postMarketTime,
			&postMarketPrice,
			&postMarketChange,
			&bid,
			&ask,
			&bidSize,
			&askSize,
			&averageDailyVolume3Month,
			&averageDailyVolume10Day,
			&fiftyTwoWeekLowChange,
			&fiftyTwoWeekLowChangePercent,
			&fiftyTwoWeekHighChange,
			&fiftyTwoWeekHighChangePercent,
			&fiftyTwoWeekLow,
			&fiftyTwoWeekHigh); scanErr != nil {
			fmt.Println(scanErr)
			return stock, scanErr
		}
		stock.Symbol = symbol
		stock.RegularMarketChange = regularMarketChange
		stock.RegularMarketOpen = regularMarketOpen
		stock.RegularMarketDayHigh = regularMarketDayHigh
		stock.RegularMarketDayLow = regularMarketDayLow
		stock.RegularMarketVolume = regularMarketVolume
		stock.RegularMarketChangePercent = regularMarketChangePercent
		stock.RegularMarketPreviousClose = regularMarketPreviousClose
		stock.RegularMarketPrice = regularMarketPrice
		stock.RegularMarketTime = regularMarketTime
		stock.EarningsTimestamp = earningsTimestamp
		stock.FiftyDayAverage = fiftyDayAverage
		stock.FiftyDayAverageChange = fiftyDayAverageChange
		stock.FiftyDayAverageChangePercent = fiftyDayAverageChangePercent
		stock.TwoHundredDayAverage = twoHundredDayAverage
		stock.TwoHundredDayAverageChange = twoHundredDayAverageChange
		stock.TwoHundredDayAverageChangePercent = twoHundredDayAverageChangePercent
		stock.PostMarketChangePercent = postMarketChangePercent
		stock.PostMarketTime = postMarketTime
		stock.PostMarketPrice = postMarketPrice
		stock.PostMarketChange = postMarketChange
		stock.Bid = bid
		stock.Ask = ask
		stock.BidSize = bidSize
		stock.AskSize = askSize
		stock.AverageDailyVolume3Month = averageDailyVolume3Month
		stock.AverageDailyVolume10Day = averageDailyVolume10Day
		stock.FiftyTwoWeekLowChange = fiftyTwoWeekLowChange
		stock.FiftyTwoWeekLowChangePercent = fiftyTwoWeekLowChangePercent
		stock.FiftyTwoWeekHighChange = fiftyTwoWeekHighChange
		stock.FiftyTwoWeekHighChangePercent = fiftyTwoWeekHighChangePercent
		stock.FiftyTwoWeekLow = fiftyTwoWeekLow
		stock.FiftyTwoWeekHigh = fiftyTwoWeekHigh

	}
	return stock, nil
}
