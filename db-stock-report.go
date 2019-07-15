package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type TblStockReport struct {
}

type RowStockReport struct {
	Symbol string
	Date   int
	MA60   float32
	MA120  float32
}

func (tbl *TblStockReport) SelectLastStockReport(symbol string) (RowStockReport, error) {
	report := RowStockReport{}
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return report, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Symbol, Date, MA60, MA120 FROM " + TBL_STOCK_REPORT_NAME + " WHERE Symbol='" + symbol + "' ORDER BY Date DESC LIMIT 1")
	if dbPrepErr != nil {
		return report, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query()
	if dbQueryErr != nil {
		return report, dbQueryErr
	}
	defer rows.Close()
	var sb string
	var date int
	var ma60 float32
	var ma120 float32
	for rows.Next() {
		if scanErr := rows.Scan(&sb, &date, &ma60, &ma120); scanErr != nil {
			fmt.Println(scanErr)
			return report, scanErr
		}
	}
	report.Symbol = sb
	report.Date = date
	report.MA60 = ma60
	report.MA120 = ma120
	return report, nil
}

//Select # of stock report (reversed array) before the date exclusively the function is called
func (tbl *TblStockReport) SelectLastStockReportByCountAndBeforeDate(symbol string, count int, beforeDate int) ([]RowStockReport, error) {
	reportList := []RowStockReport{}
	report := RowStockReport{}
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return reportList, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Symbol, Date, MA60, MA120 FROM " + TBL_STOCK_REPORT_NAME + " WHERE Symbol='" + symbol + "' AND Date<" + strconv.Itoa(beforeDate) + " ORDER BY Date DESC LIMIT " + strconv.Itoa(count))
	if dbPrepErr != nil {
		return reportList, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query()
	if dbQueryErr != nil {
		return reportList, dbQueryErr
	}
	defer rows.Close()
	var sb string
	var date int
	var ma60 float32
	var ma120 float32
	for rows.Next() {
		if scanErr := rows.Scan(&sb, &date, &ma60, &ma120); scanErr != nil {
			fmt.Println(scanErr)
			return reportList, scanErr
		}
		report.Symbol = sb
		report.Date = date
		report.MA60 = ma60
		report.MA120 = ma120
		reportList = append(reportList, report)
	}
	return reportList, nil
}

func (tbl *TblStockReport) DropTableIfExist() error {
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

func (tbl *TblStockReport) CreateTableIfNotExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string

	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		TBL_STOCK_REPORT_NAME +
		" (" +
		"Symbol VARCHAR(10) NOT NULL," +
		"Date INT NOT NULL," +
		"MA60 FLOAT(10,2)," +
		"MA120 FLOAT(10,2)," +
		"UpdatedTime TIMESTAMP," +
		"PRIMARY KEY (Symbol, Date)" +
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

func (tbl *TblStockReport) InsertOrUpdateStockData(row RowStockReport) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			TBL_STOCK_REPORT_NAME+
			" WHERE Symbol=? AND Date=?)",
		row.Symbol, row.Date)
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		return dbQueryErr
	}
	if exist == 0 {
		if err := tbl.InsertStockData(row); err != nil {
			return err
		}
	} else {
		if err := tbl.UpdateStockData(row); err != nil {
			return err
		}
	}
	return nil
}

func (tbl *TblStockReport) InsertStockData(row RowStockReport) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		TBL_STOCK_REPORT_NAME +
		" (" +
		"Symbol, " +
		"Date, " +
		"MA60," +
		"MA120," +
		"UpdatedTime) VALUES (?, ?, ?, ?, ?)")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		row.Symbol,
		row.Date,
		row.MA60,
		row.MA120,
		GetTime())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func (tbl *TblStockReport) UpdateStockData(row RowStockReport) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		TBL_STOCK_REPORT_NAME +
		" SET " +
		"MA60=?, " +
		"MA120=?, " +
		"UpdatedTime=? " +
		"WHERE Symbol=? AND Date=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		row.MA60,
		row.MA120,
		GetTime(),
		row.Symbol,
		row.Date)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}
