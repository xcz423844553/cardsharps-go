package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TblSymbolTag struct {
}

type TblSymbolTagRow struct {
	Symbol string `json:"symbol"`
	Tag    string `json:"tag"`
}

//SelectSymbolRowByTag selects a list of symbols based on the given tag
//Param: TblSymbolTagRow - the tag of the symbols
//Return: Array of TblSymbolTagRow - list of symbol with a given tag
func (tbl *TblSymbolTag) SelectSymbolRowByTag(obj TblSymbolTagRow) ([]TblSymbolTagRow, error) {
	var symbols []TblSymbolTagRow
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return symbols, dbConnErr
	}
	defer db.Close()
	var querySQL string = "SELECT * FROM " + TBL_SYMBOL_TAG + " WHERE Tag=? ORDER BY Symbol ASC"
	stmt, dbPrepErr := db.Prepare(querySQL)
	if dbPrepErr != nil {
		return symbols, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query(obj.Tag)
	if dbQueryErr != nil {
		return symbols, dbQueryErr
	}
	defer rows.Close()
	var symbol string
	var tag string
	for rows.Next() {
		if scanErr := rows.Scan(&symbol, &tag); scanErr != nil {
			fmt.Println(scanErr)
			return symbols, scanErr
		}
		symbols = append(symbols, TblSymbolTagRow{
			Symbol: symbol,
			Tag:    tag,
		})
	}
	return symbols, nil
}

//SelectTagRowBySymbol selects a list of tags based on the given symbol
//Param: TblSymbolTagRow - the symbol of the tags
//Return: Array of TblSymbolTagRow - list of tags of a given symbol
func (tbl *TblSymbolTag) SelectTagRowBySymbol(obj TblSymbolTagRow) ([]TblSymbolTagRow, error) {
	var tags []TblSymbolTagRow
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return tags, dbConnErr
	}
	defer db.Close()
	var querySQL string = "SELECT * FROM " + TBL_SYMBOL_TAG + " WHERE Symbol=?"
	stmt, dbPrepErr := db.Prepare(querySQL)
	if dbPrepErr != nil {
		return tags, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query(obj.Symbol)
	if dbQueryErr != nil {
		return tags, dbQueryErr
	}
	defer rows.Close()
	var symbol string
	var tag string
	for rows.Next() {
		if scanErr := rows.Scan(&symbol, &tag); scanErr != nil {
			fmt.Println(scanErr)
			return tags, scanErr
		}
		tags = append(tags, TblSymbolTagRow{
			Symbol: symbol,
			Tag:    tag,
		})
	}
	return tags, nil
}

func (tbl *TblSymbolTag) DropTableIfExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + TBL_SYMBOL_TAG)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl *TblSymbolTag) CreateTableIfNotExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string
	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		TBL_SYMBOL_TAG +
		" (" +
		"Symbol VARCHAR(10) NOT NULL," +
		"Tag VARCHAR(255) NOT NULL," +
		"PRIMARY KEY (Symbol, Tag)" +
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

func (tbl *TblSymbolTag) InsertOneRow(obj TblSymbolTagRow) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		TBL_SYMBOL_TAG +
		" (" +
		"Symbol, " +
		"Tag) VALUES (?,?)")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		obj.Symbol,
		obj.Tag)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func (tbl *TblSymbolTag) DeleteOneRow(obj TblSymbolTagRow) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("DELETE FROM " +
		TBL_SYMBOL_TAG +
		" WHERE Symbol=? AND Tag=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		obj.Symbol, obj.Tag)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}
