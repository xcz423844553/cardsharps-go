package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//Dao is a struct for database manipulation
type Dao struct {
}

//ConnectToDb connects to the default database
func (dao *Dao) ConnectToDb() (*sql.DB, error) {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DbName)
	if dbConnErr != nil {
		return db, dbConnErr
	}
	return db, dbConnErr
}

//CreateDb creates table if the table is not existing yet
func (dao *Dao) CreateDb() error {
	db, dbConnErr := dao.ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, dbCreateErr := db.Exec("CREATE SCHEMA IF NOT EXISTS " + DbName + " DEFAULT CHARACTER SET utf8")
	if dbCreateErr != nil {
		return dbCreateErr
	}
	return nil
}

//InitDb creates the database and tables for the entire program
func (dao *Dao) InitDb() {
	fmt.Println("The program is initiating the database and tables. Type 'y' to continue. Type any other chars to skip.")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	if strings.Compare("y", text) == 0 {
		fmt.Println("Creating the database and tables.")
		if err := dao.CreateDb(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoSymbol).DropTableIfExist(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoSymbol).CreateTableIfNotExist(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoTag).DropTableIfExist(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoTag).CreateTableIfNotExist(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoSymbolTag).DropTableIfExist(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoSymbolTag).CreateTableIfNotExist(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoStockData).DropTableIfExist(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoStockData).CreateTableIfNotExist(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoStockHist).DropTableIfExist(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoStockHist).CreateTableIfNotExist(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoOptionData).DropTableIfExist(); err != nil {
			fmt.Println(err.Error())
		}
		if err := new(DaoOptionData).CreateTableIfNotExist(); err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Finished creating the database and tables.")
	} else {
		fmt.Println("Database initiation skipped.")
	}
}
