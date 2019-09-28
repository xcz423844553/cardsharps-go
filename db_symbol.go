package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//DaoSymbol is a struct to manipulate db_symbol
type DaoSymbol struct {
	TblName string
}

//RowSymbol is a struct representing the row of db_symbol
type RowSymbol struct {
	Symbol string `json:"symbol"`
}

//GetTableName
//Return: table name
func (dao *DaoSymbol) getTableName() string {
	if TestMode {
		return TblNameSymbol + "_test"
	}
	return TblNameSymbol
}

//SelectSymbolAll selects all the symbols from the database
//Return: []]RowSymbol - list of symbol structs
func (dao *DaoSymbol) SelectSymbolAll() ([]RowSymbol, error) {
	querySQL := "SELECT * FROM " + dao.getTableName()
	return dao.selectSymbol(querySQL)
}

//selectSymbol selects a list of symbols from the database based on the given querySql
//Param: querySql - the sql to select the symbols
//Return: []RowSymbol - list of symbol structs
func (dao *DaoSymbol) selectSymbol(querySQL string) ([]RowSymbol, error) {
	var symbols []RowSymbol
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return symbols, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare(querySQL)
	if dbPrepErr != nil {
		return symbols, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query()
	if dbQueryErr != nil {
		return symbols, dbQueryErr
	}
	defer rows.Close()
	var symbol string
	for rows.Next() {
		if scanErr := rows.Scan(&symbol); scanErr != nil {
			return symbols, scanErr
		}
		symbols = append(symbols, RowSymbol{
			Symbol: symbol,
		})
	}
	return symbols, nil
}

//DropTableIfExist drops table if the table exists
func (dao *DaoSymbol) DropTableIfExist() error {
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
func (dao *DaoSymbol) CreateTableIfNotExist() error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	sqlStr := "CREATE TABLE IF NOT EXISTS " +
		dao.getTableName() +
		" (" +
		"Symbol VARCHAR(10) NOT NULL," +
		"PRIMARY KEY (Symbol)" +
		")"
	_, tblCreateErr := db.Exec(sqlStr)
	if tblCreateErr != nil {
		return tblCreateErr
	}
	return nil
}

//InsertOrUpdateRow inserts a row if it does not exist, updates a row if it exists
func (dao *DaoSymbol) InsertOrUpdateRow(obj RowSymbol) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			dao.getTableName()+
			" WHERE Symbol=?)",
		obj.Symbol)
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		return dbQueryErr
	}
	if exist == 0 {
		if err := dao.InsertRow(obj); err != nil {
			return err
		}
	} else {
		if err := dao.UpdateRow(obj); err != nil {
			return err
		}
	}
	return nil
}

//InsertRow insert a row
func (dao *DaoSymbol) InsertRow(obj RowSymbol) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		dao.getTableName() +
		" (" +
		"Symbol) VALUES (?)")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		obj.Symbol)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

//UpdateRow update a row
func (dao *DaoSymbol) UpdateRow(obj RowSymbol) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		dao.getTableName() +
		" SET " +
		"Symbol=? " +
		"WHERE Symbol=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		obj.Symbol,
		obj.Symbol)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

//DeleteRow delete a row
func (dao *DaoSymbol) DeleteRow(obj RowSymbol) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("DELETE FROM " +
		dao.getTableName() +
		" WHERE Symbol=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		obj.Symbol)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}
