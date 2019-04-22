package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type TblStockReport struct {
}

func (tbl TblStockReport) DropTableIfExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + TBL_STOCK_REPORT_NAME)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl TblStockReport) CreateTableIfNotExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string

	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		TBL_STOCK_REPORT_NAME +
		" (" + "ContractSymbol VARCHAR(21) NOT NULL," +
		"UpdatedTime TIMESTAMP," +
		"PRIMARY KEY (ContractSymbol)" +
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
