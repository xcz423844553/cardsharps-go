package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//DaoTag is a struct to manipulate db_tag
type DaoTag struct {
	TblName string
}

//RowTag is a struct representing the row of db_tag
type RowTag struct {
	Tag string `json:"tag"`
}

//GetTableName
//Return: table name
func (dao *DaoTag) getTableName() string {
	if TestMode {
		return TblNameTag + "_test"
	}
	return TblNameTag
}

//SelectTagAll selects all the tags from the database
//Return: []]RowTag - list of tag structs
func (dao *DaoTag) SelectTagAll() ([]RowTag, error) {
	querySQL := "SELECT * FROM " + dao.getTableName()
	return dao.selectTag(querySQL)
}

//selectTag selects a list of tags from the database based on the given querySql
//Param: querySql - the sql to select the tags
//Return: []RowTag - list of tag structs
func (dao *DaoTag) selectTag(querySQL string) ([]RowTag, error) {
	var tags []RowTag
	db, dbConnErr := new(Dao).ConnectToDb()
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
			return tags, scanErr
		}
		tags = append(tags, RowTag{
			Tag: tag,
		})
	}
	return tags, nil
}

//DropTableIfExist drops table if the table exists
func (dao *DaoTag) DropTableIfExist() error {
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
func (dao *DaoTag) CreateTableIfNotExist() error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	sqlStr := "CREATE TABLE IF NOT EXISTS " +
		dao.getTableName() +
		" (" +
		"Tag VARCHAR(255) NOT NULL," +
		"PRIMARY KEY (Tag)" +
		")"
	_, tblCreateErr := db.Exec(sqlStr)
	if tblCreateErr != nil {
		return tblCreateErr
	}
	return nil
}

//InsertOrUpdateRow inserts a row if it does not exist, updates a row if it exists
func (dao *DaoTag) InsertOrUpdateRow(obj RowTag) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			dao.getTableName()+
			" WHERE Tag=?)",
		obj.Tag)
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
func (dao *DaoTag) InsertRow(obj RowTag) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		dao.getTableName() +
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

//UpdateRow update a row
func (dao *DaoTag) UpdateRow(obj RowTag) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		dao.getTableName() +
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

//DeleteRow delete a row
func (dao *DaoTag) DeleteRow(obj RowTag) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("DELETE FROM " +
		dao.getTableName() +
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
