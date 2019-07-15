package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TblTag struct {
}

type TblTagRow struct {
	Tag string `json:"tag"`
}

//SelectTagRow returns all the tags in the database
//Return: Array of TblTagRow - list of tag structs
func (tbl *TblTag) SelectAllRow() ([]TblTagRow, error) {
	var tags []TblTagRow
	var querySQL string = "SELECT * FROM " + TBL_TAG
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return tags, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare(querySQL)
	if dbPrepErr != nil {
		return tags, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query()
	if dbQueryErr != nil {
		return tags, dbQueryErr
	}
	defer rows.Close()
	var tag string
	for rows.Next() {
		if scanErr := rows.Scan(&tag); scanErr != nil {
			fmt.Println(scanErr)
			return tags, scanErr
		}
		tags = append(tags, TblTagRow{
			Tag: tag,
		})
	}
	return tags, nil
}

func (tbl *TblTag) DropTableIfExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + TBL_TAG)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl *TblTag) CreateTableIfNotExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string
	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		TBL_TAG +
		" (" +
		"Tag VARCHAR(255) NOT NULL," +
		"PRIMARY KEY (Tag)" +
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

func (tbl *TblTag) InsertOrUpdateOneRow(obj TblTagRow) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			TBL_TAG+
			" WHERE Tag=?)",
		obj.Tag)
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

func (tbl *TblTag) InsertOneRow(obj TblTagRow) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		TBL_TAG +
		" (" +
		"Tag) VALUES (?)")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		obj.Tag)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func (tbl *TblTag) UpdateOneRow(obj TblTagRow) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		TBL_TAG +
		" SET " +
		"Tag=? " +
		"WHERE Tag=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		obj.Tag,
		obj.Tag)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func (tbl *TblTag) DeleteOneRow(obj TblTagRow) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("DELETE FROM " +
		TBL_TAG +
		" WHERE Tag=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		obj.Tag)
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}
