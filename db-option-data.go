package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TblOptionData struct {
}

func (tbl TblOptionData) DropTableIfExist() error {
	return tbl.CreateTableIfNotExistByTblName(TBL_OPTION_DATA_NAME)
}

func (tbl TblOptionData) DropTableIfExistByTblName(tblName string) error {
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

func (tbl TblOptionData) CreateTableIfNotExist() error {
	return tbl.CreateTableIfNotExistByTblName(TBL_OPTION_DATA_NAME)
}

func (tbl TblOptionData) CreateTableIfNotExistByTblName(tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string
	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		tblName +
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
		"PrevVolume INT," +
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
	return tbl.InsertOrUpdateOptionDataToTbl(option, TBL_OPTION_DATA_NAME)
}

func (tbl TblOptionData) InsertOrUpdateOptionDataToEtf(option YahooOption) error {
	return tbl.InsertOrUpdateOptionDataToTbl(option, TBL_OPTION_DATA_ETF_NAME)
}

func (tbl TblOptionData) InsertOrUpdateOptionDataToTbl(option YahooOption, tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			tblName+
			" WHERE ContractSymbol=? AND Date=?)",
		option.ContractSymbol, GetTimeInYYYYMMDD())
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		return dbQueryErr
	}
	if exist == 0 {
		if err := tbl.InsertOptionData(option, tblName); err != nil {
			return err
		}
	} else {
		if err := tbl.UpdateOptionData(option, tblName); err != nil {
			return err
		}
	}
	return nil
}

func (tbl TblOptionData) InsertOptionData(option YahooOption, tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		tblName +
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
		"PrevVolume, " +
		"OpenInterest, " +
		"Bid, " +
		"Ask, " +
		"Expiration, " +
		"ImpliedVolatility, " +
		"InTheMoney, " +
		"UpdatedTime) VALUES (?,?,?,?,?,?,?,?,?,0,?,?,?,?,?,?,?)")
	if dbPrepErr != nil {
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

func (tbl TblOptionData) UpdateOptionData(option YahooOption, tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	prevVolume, selectErr := tbl.SelectOptionDataVolumeByContractSymbolAndDate(option.ContractSymbol, GetTimeInYYYYMMDD())
	if selectErr != nil {
		return selectErr
	}
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		tblName +
		" SET " +
		"LastPrice=?, " +
		"PriceChange=?, " +
		"PercentChange=?, " +
		"Volume=?, " +
		"PrevVolume=?, " +
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
		prevVolume,
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

func (tbl TblOptionData) SelectOptionDataVolumeByContractSymbolAndDate(contractSymbol string, date int) (int, error) {
	return tbl.SelectOptionDataVolumeByContractSymbolAndDateFromTbl(contractSymbol, date, TBL_OPTION_DATA_NAME)
}

func (tbl TblOptionData) SelectOptionDataVolumeByContractSymbolAndDateToEtf(contractSymbol string, date int) (int, error) {
	return tbl.SelectOptionDataVolumeByContractSymbolAndDateFromTbl(contractSymbol, date, TBL_OPTION_DATA_ETF_NAME)
}

func (tbl TblOptionData) SelectOptionDataVolumeByContractSymbolAndDateFromTbl(contractSymbol string, date int, tblName string) (int, error) {
	var volume int
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return volume, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Volume FROM " + tblName + " WHERE " + "ContractSymbol=? AND Date=?")
	if dbPrepErr != nil {
		return volume, dbPrepErr
	}
	defer stmt.Close()
	optionRows, dbQueryErr := stmt.Query(contractSymbol, date)
	if dbQueryErr != nil {
		return volume, dbQueryErr
	}
	defer optionRows.Close()
	for optionRows.Next() {
		if scanErr := optionRows.Scan(&volume); scanErr != nil {
			fmt.Println(scanErr)
			return volume, scanErr
		}
	}
	return volume, nil
}
