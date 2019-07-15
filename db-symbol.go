package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type TblSymbol struct {
}

type TblSymbolRow struct {
	Symbol          string `json:"symbol"`
	Sp500           bool   `json:"sp500"`
	Nasdaq          bool   `json:"nasdaq"`
	Dow             bool   `json:"dow"`
	Russell         bool   `json:"russell"`
	ETF             bool   `json:"etf"`
	StockMonitored  bool   `json:"stockMonitored"`
	OptionMonitored bool   `json:"optionMonitored"`
}

//SelectAllSymbol selects all the symbols from the database
//Return: symbols - list of symbol strings
func (tbl *TblSymbol) SelectAllSymbol() ([]string, error) {
	querySQL := "SELECT * FROM " + TBL_SYMBOL
	return tbl.SelectSymbol(querySQL)
}

//SelectAllSymbolRow selects all the symbols from the database
//Return: Array of TblSymbolRow - list of symbol structs
func (tbl *TblSymbol) SelectAllSymbolRow() ([]TblSymbolRow, error) {
	querySQL := "SELECT * FROM " + TBL_SYMBOL
	return tbl.SelectSymbolRow(querySQL)
}

//SelectSymbolByAndFilter selects a list of symbols from the database which satisfy all the given criteria at the same time
//Param: isSp500 - true if the symbol is marked as Sp500
//Param: isNasdaq - true if the symbol is marked as Nasdaq
//Param: isDow - true if the symbol is marked as Dow
//Param: isRussell - true if the symbol is marked as Russell
//Param: isETF - true if the symbol is ETF
//Param: isStockMonitored - true if the symbol's stock is monitored
//Param: isOptionMonitored - true if the symbol's option is monitored
//Return: symbols - list of symbol strings
func (tbl *TblSymbol) SelectSymbolByAndFilter(isSp500 bool, isNasdaq bool, isDow bool, isRussell bool, isETF bool, isStockMonitored bool, isOptionMonitored bool) ([]string, error) {
	var symbols []string
	list, err := tbl.SelectSymbolRowByAndFilter(isSp500, isNasdaq, isDow, isRussell, isETF, isStockMonitored, isOptionMonitored)
	if err != nil {
		return symbols, err
	}
	for _, row := range list {
		symbols = append(symbols, row.Symbol)
	}
	return symbols, nil
}

//SelectSymbolByOrFilter selects a list of symbols from the database which satisfy one of the given criteria
//Param: isSp500 - true if the symbol is marked as Sp500
//Param: isNasdaq - true if the symbol is marked as Nasdaq
//Param: isDow - true if the symbol is marked as Dow
//Param: isRussell - true if the symbol is marked as Russell
//Param: isETF - true if the symbol is ETF
//Param: isStockMonitored - true if the symbol's stock is monitored
//Param: isOptionMonitored - true if the symbol's option is monitored
//Return: symbols - list of symbol strings
func (tbl *TblSymbol) SelectSymbolByOrFilter(isSp500 bool, isNasdaq bool, isDow bool, isRussell bool, isETF bool, isStockMonitored bool, isOptionMonitored bool) ([]string, error) {
	var symbols []string
	list, err := tbl.SelectSymbolRowByOrFilter(isSp500, isNasdaq, isDow, isRussell, isETF, isStockMonitored, isOptionMonitored)
	if err != nil {
		return symbols, err
	}
	for _, row := range list {
		symbols = append(symbols, row.Symbol)
	}
	return symbols, nil
}

//SelectSymbol selects a list of symbols from the database based on the given querySql
//Param: querySql - the sql to select the symbols
//Return: symbols - list of symbol strings
func (tbl *TblSymbol) SelectSymbol(querySQL string) ([]string, error) {
	var symbols []string
	list, err := tbl.SelectSymbolRow(querySQL)
	if err != nil {
		return symbols, err
	}
	for _, row := range list {
		symbols = append(symbols, row.Symbol)
	}
	return symbols, nil
}

//SelectSymbolRowByOrFilter selects a list of symbols from the database which satisfy one of the given criteria
//Param: isSp500 - true if the symbol is marked as Sp500
//Param: isNasdaq - true if the symbol is marked as Nasdaq
//Param: isDow - true if the symbol is marked as Dow
//Param: isRussell - true if the symbol is marked as Russell
//Param: isETF - true if the symbol is ETF
//Param: isStockMonitored - true if the symbol's stock is monitored
//Param: isOptionMonitored - true if the symbol's option is monitored
//Return: Array of TblSymbolRow - list of symbol structs
func (tbl *TblSymbol) SelectSymbolRowByOrFilter(isSp500 bool, isNasdaq bool, isDow bool, isRussell bool, isETF bool, isStockMonitored bool, isOptionMonitored bool) ([]TblSymbolRow, error) {
	querySQL := "SELECT * FROM " + TBL_SYMBOL
	var crit []string
	if isSp500 {
		crit = append(crit, "Sp500=true")
	}
	if isNasdaq {
		crit = append(crit, "Nasdaq=true")
	}
	if isDow {
		crit = append(crit, "Dow=true")
	}
	if isRussell {
		crit = append(crit, "Russell=true")
	}
	if isETF {
		crit = append(crit, "ETF=true")
	}
	if isStockMonitored {
		crit = append(crit, "StockMonitored=true")
	}
	if isOptionMonitored {
		crit = append(crit, "OptionMonitored=true")
	}
	if len(crit) > 0 {
		querySQL += " WHERE " + strings.Join(crit, " OR ")
	}
	return tbl.SelectSymbolRow(querySQL)
}

//SelectSymbolRowByAndFilter selects a list of symbols from the database which satisfy all the given criteria at the same time
//Param: isSp500 - true if the symbol is marked as Sp500
//Param: isNasdaq - true if the symbol is marked as Nasdaq
//Param: isDow - true if the symbol is marked as Dow
//Param: isRussell - true if the symbol is marked as Russell
//Param: isETF - true if the symbol is ETF
//Param: isStockMonitored - true if the symbol's stock is monitored
//Param: isOptionMonitored - true if the symbol's option is monitored
//Return: Array of TblSymbolRow - list of symbol structs
func (tbl *TblSymbol) SelectSymbolRowByAndFilter(isSp500 bool, isNasdaq bool, isDow bool, isRussell bool, isETF bool, isStockMonitored bool, isOptionMonitored bool) ([]TblSymbolRow, error) {
	querySQL := "SELECT * FROM " + TBL_SYMBOL
	var crit []string
	if isSp500 {
		crit = append(crit, "Sp500=true")
	}
	if isNasdaq {
		crit = append(crit, "Nasdaq=true")
	}
	if isDow {
		crit = append(crit, "Dow=true")
	}
	if isRussell {
		crit = append(crit, "Russell=true")
	}
	if isETF {
		crit = append(crit, "ETF=true")
	}
	if isStockMonitored {
		crit = append(crit, "StockMonitored=true")
	}
	if isOptionMonitored {
		crit = append(crit, "OptionMonitored=true")
	}
	if len(crit) > 0 {
		querySQL += " WHERE " + strings.Join(crit, " AND ")
	}
	return tbl.SelectSymbolRow(querySQL)
}

//SelectSymbolRow selects a list of symbols from the database based on the given querySql
//Param: querySql - the sql to select the symbols
//Return: Array of TblSymbolRow - list of symbol structs
func (tbl *TblSymbol) SelectSymbolRow(querySQL string) ([]TblSymbolRow, error) {
	var symbols []TblSymbolRow
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
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
	var sp500 bool
	var nasdaq bool
	var dow bool
	var russell bool
	var etf bool
	var stockMonitored bool
	var optionMonitored bool
	for rows.Next() {
		if scanErr := rows.Scan(&symbol, &sp500, &nasdaq, &dow, &russell, &etf, &stockMonitored, &optionMonitored); scanErr != nil {
			fmt.Println(scanErr)
			return symbols, scanErr
		}
		symbols = append(symbols, TblSymbolRow{
			Symbol:          symbol,
			Sp500:           sp500,
			Nasdaq:          nasdaq,
			Dow:             dow,
			Russell:         russell,
			ETF:             etf,
			StockMonitored:  stockMonitored,
			OptionMonitored: optionMonitored,
		})
	}
	return symbols, nil
}

func (tbl *TblSymbol) DropTableIfExist() error {
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

func (tbl *TblSymbol) CreateTableIfNotExist() error {
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
		"ETF BOOLEAN," +
		"StockMonitored BOOLEAN," +
		"OptionMonitored BOOLEAN," +
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

func (tbl *TblSymbol) InsertOrUpdateOneRow(obj TblSymbolRow) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			TBL_SYMBOL+
			" WHERE Symbol=?)",
		obj.Symbol)
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		return dbQueryErr
	}
	if exist == 0 {
		if err := tbl.InsertOneRow(obj); err != nil {
			return err
		}
	} else {
		if err := tbl.UpdateOneRow(obj); err != nil {
			return err
		}
	}
	return nil
}

func (tbl *TblSymbol) InsertOneRow(obj TblSymbolRow) error {
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
		"Russell, " +
		"ETF, " +
		"StockMonitored, " +
		"OptionMonitored) VALUES (?,?,?,?,?,?,?,?)")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		obj.Symbol,
		obj.Sp500,
		obj.Nasdaq,
		obj.Dow,
		obj.Russell,
		obj.ETF,
		obj.StockMonitored,
		obj.OptionMonitored)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func (tbl *TblSymbol) UpdateOneRow(obj TblSymbolRow) error {
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
		"ETF=? " +
		"StockMonitored=? " +
		"OptionMonitored=? " +
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
		obj.ETF,
		obj.StockMonitored,
		obj.OptionMonitored,
		obj.Symbol)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func (tbl *TblSymbol) DeleteOneRow(obj TblSymbolRow) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("DELETE FROM " +
		TBL_SYMBOL +
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
