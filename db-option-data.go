package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type TblOptionData struct {
}

func (tbl TblOptionData) DropTableIfExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + TBL_OPTION_DATA_NAME)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl TblOptionData) CreateTableIfNotExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string
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
	if sqlStr != "" {
		_, tblCreateErr := db.Exec(sqlStr)
		if tblCreateErr != nil {
			return tblCreateErr
		}
		return nil
	}
	return errors.New("failed to find preset table name")
}

func (tbl TblOptionData) InsertOrUpdateOptionData(option YahooOption) error {
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
	if exist == 0 {
		if err := tbl.InsertOptionData(option); err != nil {
			panic(err)
			return err
		}
	} else {
		if err := tbl.UpdateOptionData(option); err != nil {
			panic(err)
			return err
		}
	}
	return nil
}

func (tbl TblOptionData) InsertOptionData(option YahooOption) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
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

func (tbl TblOptionData) UpdateOptionData(option YahooOption) error {
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
		GetTime(),
		option.ContractSymbol,
		GetTimeInYYYYMMDD())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}
