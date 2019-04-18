package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func initDb() {
	var err error
	if err = createDatabaseIfNotExist(DB_NAME); err != nil {
		panic(err)
	}
	if err = dropTableIfExist(TBL_OPTION_DATA_NAME); err != nil {
		panic(err)
	}
	if err = dropTableIfExist(TBL_OPTION_REPORT_NAME); err != nil {
		panic(err)
	}
	if err = dropTableIfExist(TBL_STOCK_DATA_NAME); err != nil {
		panic(err)
	}
	if err = dropTableIfExist(TBL_STOCK_REPORT_NAME); err != nil {
		panic(err)
	}
	if err = createTableIfNotExist(TBL_OPTION_DATA_NAME); err != nil {
		panic(err)
	}
	if err = createTableIfNotExist(TBL_OPTION_REPORT_NAME); err != nil {
		panic(err)
	}
	if err = createTableIfNotExist(TBL_STOCK_DATA_NAME); err != nil {
		panic(err)
	}
	if err = createTableIfNotExist(TBL_STOCK_REPORT_NAME); err != nil {
		panic(err)
	}
}

func createDatabaseIfNotExist(dbName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, dbCreateErr := db.Exec("CREATE SCHEMA IF NOT EXISTS " + dbName + " DEFAULT CHARACTER SET utf8")
	if dbCreateErr != nil {
		return dbCreateErr
	}
	return nil
}

func dropTableIfExist(tblName string) error {
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

func createTableIfNotExist(tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string
	switch tblName {
	case TBL_OPTION_DATA_NAME:
		sqlStr = "CREATE TABLE IF NOT EXISTS " +
			TBL_OPTION_DATA_NAME +
			" (" +
			//Contract symbol is Symbol(6) + yymmdd(6) + P/C(1) + Price(8) = VARCHAR(21)
			"ContractSymbol VARCHAR(21) NOT NULL," +
			"Date INT NOT NULL," +
			"Symbol VARCHAR(10) NOT NULL," +
			"OptionType VARCHAR(1) NOT NULL," +
			"Strike FLOAT(10,2)," +
			"LastPrice FLOAT(10,2)," +
			"PriceChange FLOAT(10,2)," +
			"PercentChange FLOAT(10,2)," +
			"Volume INT," +
			"OpenInterest INT," +
			"Bid FLOAT(10,2)," +
			"Ask FLOAT(10,2)," +
			"Expiration BIGINT," +
			"ImpliedVolatility FLOAT(10,2)," +
			"InTheMoney BOOLEAN," +
			"UpdatedTime TIMESTAMP," +
			"PRIMARY KEY (ContractSymbol, Date)" +
			")"
	case TBL_OPTION_REPORT_NAME:
		//TODO
		sqlStr = "CREATE TABLE IF NOT EXISTS " +
			TBL_OPTION_REPORT_NAME +
			" (" + "ContractSymbol VARCHAR(21) NOT NULL," +
			"UpdatedTime TIMESTAMP," +
			"PRIMARY KEY (ContractSymbol)" +
			")"
	case TBL_STOCK_DATA_NAME:
		sqlStr = "CREATE TABLE IF NOT EXISTS " +
			TBL_STOCK_DATA_NAME +
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
	case TBL_STOCK_REPORT_NAME:
		//TODO
		sqlStr = "CREATE TABLE IF NOT EXISTS " +
			TBL_STOCK_REPORT_NAME +
			" (" + "ContractSymbol VARCHAR(21) NOT NULL," +
			"UpdatedTime TIMESTAMP," +
			"PRIMARY KEY (ContractSymbol)" +
			")"
	}
	if sqlStr != "" {
		_, tblCreateErr := db.Exec(sqlStr)
		if tblCreateErr != nil {
			return tblCreateErr
		}
		return nil
	}
	return errors.New("failed to find preset table name")
}

func InsertOrUpdateOptionData(option YahooOption) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			TBL_OPTION_DATA_NAME+
			" WHERE ContractSymbol=? AND Date=?)",
		option.ContractSymbol, GetTimeInYYYYMMDD())
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		panic(dbQueryErr)
		return dbQueryErr
	}
	fmt.Printf("%+v\n", option)
	fmt.Println(strconv.Itoa(exist) + " " + option.ContractSymbol)
	if exist == 0 {
		if err := InsertOptionData(option); err != nil {
			panic(err)
			return err
		}
	} else {
		if err := UpdateOptionData(option); err != nil {
			panic(err)
			return err
		}
	}
	return nil
}

func InsertOptionData(option YahooOption) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()

	fmt.Printf("Here%+v\n", option)
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		TBL_OPTION_DATA_NAME +
		" (" +
		"ContractSymbol, " +
		"Date, " +
		"Symbol, " +
		"OptionType, " +
		"Strike, " +
		"LastPrice, " +
		"PriceChange, " +
		"PercentChange, " +
		"Volume, " +
		"OpenInterest, " +
		"Bid, " +
		"Ask, " +
		"Expiration, " +
		"ImpliedVolatility, " +
		"InTheMoney, " +
		"UpdatedTime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if dbPrepErr != nil {
		panic(dbPrepErr)
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		option.ContractSymbol,
		GetTimeInYYYYMMDD(),
		option.GetSymbol(),
		option.GetOptionType(),
		option.Strike,
		option.LastPrice,
		option.PriceChange,
		option.PercentChange,
		option.Volume,
		option.OpenInterest,
		option.Bid,
		option.Ask,
		option.Expiration,
		option.ImpliedVolatility,
		option.InTheMoney,
		GetTime())
	if dbExecErr != nil {
		panic(dbExecErr)
		return dbExecErr
	}
	return nil
}

func UpdateOptionData(option YahooOption) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		TBL_OPTION_DATA_NAME +
		" SET " +
		"LastPrice=?, " +
		"PriceChange=?, " +
		"PercentChange=?, " +
		"Volume=?, " +
		"OpenInterest=?, " +
		"Bid=?, " +
		"Ask=?, " +
		"ImpliedVolatility=?, " +
		"InTheMoney=?, " +
		"UpdatedTime=? " +
		"WHERE ContractSymbol=? AND Date=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		option.LastPrice,
		option.PriceChange,
		option.PercentChange,
		option.Volume,
		option.OpenInterest,
		option.Bid,
		option.Ask,
		option.ImpliedVolatility,
		option.InTheMoney,
		GetTime())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func InsertOrUpdateStockData(stock YahooQuote) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			TBL_STOCK_DATA_NAME+
			" WHERE Symbol=? AND Date=?)",
		stock.Symbol, GetTimeInYYYYMMDD())
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		panic(dbQueryErr)
		return dbQueryErr
	}
	if exist == 0 {
		if err := InsertStockData(stock); err != nil {
			return err
		}
	} else {
		if err := UpdateStockData(stock); err != nil {
			return err
		}
	}
	return nil
}

func InsertStockData(stock YahooQuote) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		TBL_STOCK_DATA_NAME +
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
		panic(dbPrepErr)
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
		panic(dbExecErr)
		return dbExecErr
	}
	return nil
}

func UpdateStockData(stock YahooQuote) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		TBL_STOCK_DATA_NAME +
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
		panic(dbPrepErr)
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
		GetTime())
	if dbExecErr != nil {
		panic(dbExecErr)
		return dbExecErr
	}
	return nil
}
