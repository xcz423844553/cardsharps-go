package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func initDb() {
	fmt.Println("initDb() starts")
	var err error
	if err = createDatabaseIfNotExist(DB_NAME); err != nil {
		panic(err)
	}
	// if err = new(TblSymbol).DropTableIfExist(); err != nil {
	// 	panic(err)
	// }
	if err = new(TblOptionData).DropTableIfExist(); err != nil {
		panic(err)
	}
	if err = new(TblOptionData).DropTableIfExistByTblName(TBL_OPTION_DATA_ETF_NAME); err != nil {
		panic(err)
	}
	if err = new(TblOptionReport).DropTableIfExist(); err != nil {
		panic(err)
	}
	if err = new(TblStockData).DropTableIfExist(); err != nil {
		panic(err)
	}
	if err = new(TblStockReport).DropTableIfExist(); err != nil {
		panic(err)
	}
	// if err = new(TblSymbol).CreateTableIfNotExist(); err != nil {
	// 	panic(err)
	// }
	if err = new(TblOptionData).CreateTableIfNotExist(); err != nil {
		panic(err)
	}
	if err = new(TblOptionData).CreateTableIfNotExistByTblName(TBL_OPTION_DATA_ETF_NAME); err != nil {
		panic(err)
	}
	if err = new(TblOptionReport).CreateTableIfNotExist(); err != nil {
		panic(err)
	}
	if err = new(TblStockData).CreateTableIfNotExist(); err != nil {
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
