package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TblSymbol struct {
}

type TblSymbolRow struct {
	Symbol  string `json:"symbol"`
	Sp500   bool   `json:"sp500"`
	Nasdaq  bool   `json:"nasdaq"`
	Dow     bool   `json:"dow"`
	Russell bool   `json:"russell"`
}

func (tbl TblSymbol) SelectSymbolByFilter() ([]string, error) {
	var symbols []string
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return symbols, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Symbol FROM " + TBL_SYMBOL + " ")
	if dbPrepErr != nil {
		return symbols, dbPrepErr
	}
	defer stmt.Close()
	optionRows, dbQueryErr := stmt.Query()
	if dbQueryErr != nil {
		return symbols, dbQueryErr
	}
	defer optionRows.Close()
	var symbol string
	for optionRows.Next() {
		if scanErr := optionRows.Scan(&symbol); scanErr != nil {
			fmt.Println(scanErr)
			return symbols, scanErr
		}
		symbols = append(symbols, symbol)
	}
	return symbols, nil
}

func (tbl TblSymbol) SelectSymbolByTrader(trader string) ([]string, error) {
	sqlPrefix := "SELECT Symbol FROM " + TBL_SYMBOL + " "
	sqlPostfix := "WHERE "
	if trader == TRADER_SP500 {
		sqlPostfix += "Sp500=1"
	} else if trader == TRADER_NASDAQ {
		sqlPostfix += "Nasdaq=1"
	} else if trader == TRADER_DOW {
		sqlPostfix += "Dow=1"
	} else if trader == TRADER_RUSSELL {
		sqlPostfix += "Russell=1"
	}
	var symbols []string
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return symbols, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare(sqlPrefix + sqlPostfix)
	if dbPrepErr != nil {
		return symbols, dbPrepErr
	}
	defer stmt.Close()
	optionRows, dbQueryErr := stmt.Query()
	if dbQueryErr != nil {
		return symbols, dbQueryErr
	}
	defer optionRows.Close()
	var symbol string
	for optionRows.Next() {
		if scanErr := optionRows.Scan(&symbol); scanErr != nil {
			fmt.Println(scanErr)
			return symbols, scanErr
		}
		symbols = append(symbols, symbol)
	}
	return symbols, nil
}

func (tbl TblSymbol) DropTableIfExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + TBL_SYMBOL)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl TblSymbol) CreateTableIfNotExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string
	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		TBL_SYMBOL +
		" (" +
		"Symbol VARCHAR(10) NOT NULL," +
		"Sp500 BOOLEAN," +
		"Nasdaq BOOLEAN," +
		"Dow BOOLEAN, " +
		"Russell BOOLEAN," +
		"PRIMARY KEY (Symbol)" +
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

func (tbl TblSymbol) InsertOrUpdateSymbol(obj TblSymbolRow) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			TBL_SYMBOL+
			" WHERE symbol=?)",
		obj.Symbol)
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		panic(dbQueryErr)
		return dbQueryErr
	}
	if exist == 0 {
		if err := tbl.InsertSymbol(obj); err != nil {
			panic(err)
			return err
		}
	} else {
		if err := tbl.UpdateSymbol(obj); err != nil {
			panic(err)
			return err
		}
	}
	return nil
}

func (tbl TblSymbol) InsertSymbol(obj TblSymbolRow) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		TBL_SYMBOL +
		" (" +
		"Symbol, " +
		"Sp500, " +
		"Nasdaq, " +
		"Dow, " +
		"Russell) VALUES (?,?,?,?,?)")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		obj.Symbol,
		obj.Sp500,
		obj.Nasdaq,
		obj.Dow,
		obj.Russell)
	if dbExecErr != nil {
		panic(dbExecErr)
		return dbExecErr
	}
	return nil
}

func (tbl TblSymbol) UpdateSymbol(obj TblSymbolRow) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		TBL_SYMBOL +
		" SET " +
		"Sp500=?, " +
		"Nasdaq=?, " +
		"Dow=? " +
		"Russell=? " +
		"WHERE Symbol=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		obj.Sp500,
		obj.Nasdaq,
		obj.Dow,
		obj.Russell,
		obj.Symbol)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}
