package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func initDb() {
	fmt.Println("initDb() starts")
	if err2 := new(TblStockHist).DropTableIfExist(); err2 != nil {
		panic(err2)
	}
	if err3 := new(TblStockHist).CreateTableIfNotExist(); err3 != nil {
		panic(err3)
	}
	if err4 := new(TblStockReport).DropTableIfExist(); err4 != nil {
		panic(err4)
	}
	if err5 := new(TblStockReport).CreateTableIfNotExist(); err5 != nil {
		panic(err5)
	}
	fmt.Println("initDb() ends")
	return
	var err error
	if err = createDatabaseIfNotExist(DB_NAME); err != nil {
		panic(err)
	}
	// if err = new(TblSymbol).DropTableIfExist(); err != nil {
	// 	panic(err)
	// }
	if err = new(TblLogError).DropTableIfExist(); err != nil {
		panic(err)
	}
	if err = new(TblLogSystem).DropTableIfExist(); err != nil {
		panic(err)
	}
	if err = new(TblOptionData).DropTableIfExist(); err != nil {
		panic(err)
	}
	if err = new(TblOptionReport).DropTableIfExist(); err != nil {
		panic(err)
	}
	if err = new(TblStockData).DropTableIfExist(); err != nil {
		panic(err)
	}
	if err = new(TblStockHist).DropTableIfExist(); err != nil {
		panic(err)
	}
	if err = new(TblStockReport).DropTableIfExist(); err != nil {
		panic(err)
	}
	// if err = new(TblSymbol).CreateTableIfNotExist(); err != nil {
	// 	panic(err)
	// }
	if err = new(TblLogError).CreateTableIfNotExist(); err != nil {
		panic(err)
	}
	if err = new(TblLogSystem).CreateTableIfNotExist(); err != nil {
		panic(err)
	}
	if err = new(TblOptionData).CreateTableIfNotExist(); err != nil {
		panic(err)
	}
	if err = new(TblOptionReport).CreateTableIfNotExist(); err != nil {
		panic(err)
	}
	if err = new(TblStockData).CreateTableIfNotExist(); err != nil {
		panic(err)
	}
	if err = new(TblStockHist).CreateTableIfNotExist(); err != nil {
		panic(err)
	}
	if err = new(TblStockReport).CreateTableIfNotExist(); err != nil {
		panic(err)
	}
	fmt.Println("initDb() ends")
}

func createDatabaseIfNotExist(dbName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, dbCreateErr := db.Exec("CREATE SCHEMA IF NOT EXISTS " + dbName + " DEFAULT CHARACTER SET utf8")
	if dbCreateErr != nil {
		return dbCreateErr
	}
	return nil
}
