package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TblLogSystem struct {
}

type LogSystem struct {
	LogId      int64  `json:"logid"`
	LogType    string `json:"logtype"`
	LogContent string `json:"logcontent"`
	Date       int    `json:"date"`
}

func (tbl TblLogSystem) SelectLogSystemByDate(date1 int, date2 int) ([]LogSystem, error) {
	if date1 > date2 {
		tmp := date1
		date1 = date2
		date2 = tmp
	}
	var logSystems []LogSystem
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return logSystems, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT LogId, LogType, LogContent, Date FROM " + TBL_LOG_SYSTEM + " WHERE Date>=? AND Date<=?;")
	if dbPrepErr != nil {
		return logSystems, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query(date1, date2)
	if dbQueryErr != nil {
		return logSystems, dbQueryErr
	}
	defer rows.Close()
	logSystem := new(LogSystem)
	var varLogId int64
	var varLogType string
	var varLogContent string
	var varDate int
	for rows.Next() {
		if scanErr := rows.Scan(&varLogId, &varLogType, &varLogContent, &varDate); scanErr != nil {
			fmt.Println(scanErr)
			return logSystems, scanErr
		}
		logSystem.LogId = varLogId
		logSystem.LogType = varLogType
		logSystem.LogContent = varLogContent
		logSystem.Date = varDate
		logSystems = append(logSystems, *logSystem)
	}
	return logSystems, nil
}

func (tbl TblLogSystem) DropTableIfExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + TBL_LOG_SYSTEM)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl TblLogSystem) CreateTableIfNotExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string
	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		TBL_LOG_SYSTEM +
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

func (tbl TblLogSystem) InsertLogSystem(logType string, logContent string) {
	PrintMsgInConsole(MSGSYSTEM, logType, logContent)
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		panic(dbConnErr)
		return
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		TBL_LOG_SYSTEM +
		" (" +
		"LogType, " +
		"LogContent, " +
		"Date, " +
		"UpdatedTime) VALUES (?,?,?,?)")
	if dbPrepErr != nil {
		panic(dbPrepErr)
		return
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		logType,
		logContent,
		GetTimeInYYYYMMDD(),
		GetTime())
	if dbExecErr != nil {
		panic(dbExecErr)
		return
	}
	return
}
