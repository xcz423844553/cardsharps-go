package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TblLogError struct {
}

type LogError struct {
	LogId      int64  `json:"logid"`
	LogType    string `json:"logtype"`
	LogContent string `json:"logcontent"`
	Date       int    `json:"date"`
}

func (tbl TblLogError) SelectLogErrorByDateRange(date1 int, date2 int) ([]LogError, error) {
	if date1 > date2 {
		tmp := date1
		date1 = date2
		date2 = tmp
	}
	var logErrors []LogError
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return logErrors, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT LogId, LogType, LogContent, Date FROM " + TBL_LOG_ERROR + " WHERE Date>=? AND Date<=?;")
	if dbPrepErr != nil {
		return logErrors, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query(date1, date2)
	if dbQueryErr != nil {
		return logErrors, dbQueryErr
	}
	defer rows.Close()
	logError := new(LogError)
	var varLogId int64
	var varLogType string
	var varLogContent string
	var varDate int
	for rows.Next() {
		if scanErr := rows.Scan(&varLogId, &varLogType, &varLogContent, &varDate); scanErr != nil {
			fmt.Println(scanErr)
			return logErrors, scanErr
		}
		logError.LogId = varLogId
		logError.LogType = varLogType
		logError.LogContent = varLogContent
		logError.Date = varDate
		logErrors = append(logErrors, *logError)
	}
	return logErrors, nil
}

func (tbl TblLogError) DropTableIfExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + TBL_LOG_ERROR)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl TblLogError) CreateTableIfNotExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string
	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		TBL_LOG_ERROR +
		" (" +
		"LogId INT NOT NULL AUTO_INCREMENT," +
		"LogType VARCHAR(50)," +
		"LogContent TEXT," +
		"Date INT NOT NULL," +
		"UpdatedTime TIMESTAMP, " +
		"PRIMARY KEY (LogId)" +
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

func (tbl TblLogError) InsertLogError(logType string, logContent string) {
	PrintMsgInConsole(MSGERROR, logType, logContent)
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		TBL_LOG_ERROR +
		" (" +
		"LogType, " +
		"LogContent, " +
		"Date, " +
		"UpdatedTime) VALUES (?,?,?,?)")
	if dbPrepErr != nil {
		return
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		logType,
		logContent,
		GetTimeInYYYYMMDD(),
		GetTime())
	if dbExecErr != nil {
		return
	}
	return
}
