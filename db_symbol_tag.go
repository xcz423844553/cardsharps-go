package main

import (
	_ "github.com/go-sql-driver/mysql"
)

//DaoSymbolTag is a struct to manipulate db_symbol_tag
type DaoSymbolTag struct {
}

//RowSymbolTag is a struct representing the row of db_symbol_tag
type RowSymbolTag struct {
	Symbol string `json:"symbol"`
	Tag    string `json:"tag"`
}

//GetTableName
//Return: table name
func (dao *DaoSymbolTag) getTableName() string {
	if TestMode {
		return TblNameSymbolTag + "_test"
	}
	return TblNameSymbolTag
}

//SelectSymbolByTag selects a list of symbols based on the given tag
//Param: RowSymbolTag - the tag of the symbols
//Return: []RowSymbolTag - list of symbol with a given tag
func (dao *DaoSymbolTag) SelectSymbolByTag(obj RowSymbolTag) ([]RowSymbolTag, error) {
	var symbols []RowSymbolTag
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return symbols, dbConnErr
	}
	defer db.Close()
	querySQL := "SELECT * FROM " + dao.getTableName() + " WHERE Tag=? ORDER BY Symbol ASC"
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
			return symbols, scanErr
		}
		symbols = append(symbols, RowSymbolTag{
			Symbol: symbol,
			Tag:    tag,
		})
	}
	return symbols, nil
}

//SelectTagBySymbol selects a list of tags based on the given symbol
//Param: RowSymbolTag - the symbol of the tags
//Return: []RowSymbolTag - list of tags of a given symbol
func (dao *DaoSymbolTag) SelectTagBySymbol(obj RowSymbolTag) ([]RowSymbolTag, error) {
	var tags []RowSymbolTag
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return tags, dbConnErr
	}
	defer db.Close()
	querySQL := "SELECT * FROM " + dao.getTableName() + " WHERE Symbol=?"
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
			return tags, scanErr
		}
		tags = append(tags, RowSymbolTag{
			Symbol: symbol,
			Tag:    tag,
		})
	}
	return tags, nil
}

//DropTableIfExist drops table if the table exists
func (dao *DaoSymbolTag) DropTableIfExist() error {
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
func (dao *DaoSymbolTag) CreateTableIfNotExist() error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	sqlStr := "CREATE TABLE IF NOT EXISTS " +
		dao.getTableName() +
		" (" +
		"Symbol VARCHAR(10) NOT NULL," +
		"Tag VARCHAR(255) NOT NULL," +
		"PRIMARY KEY (Symbol, Tag)" +
		")"
	_, tblCreateErr := db.Exec(sqlStr)
	if tblCreateErr != nil {
		return tblCreateErr
	}
	return nil
}

//InsertRow insert a row
func (dao *DaoSymbolTag) InsertRow(obj RowSymbolTag) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		dao.getTableName() +
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

//DeleteRow delete a row
func (dao *DaoSymbolTag) DeleteRow(obj RowSymbolTag) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("DELETE FROM " +
		dao.getTableName() +
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
