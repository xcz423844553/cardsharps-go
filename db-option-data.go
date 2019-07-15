package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TblOptionData struct {
}

type RowOptionData struct {
	ContractSymbol    string  `json:"contractSymbol"`
	Date              int     `json:"date"`
	Symbol            string  `json:"symbol"`
	OptionType        string  `json:"optionType"`
	Strike            float32 `json:"strike"`
	LastPrice         float32 `json:"lastPrice"`
	Volume            int     `json:"volume"`
	OpenInterest      int     `json:"openInterest"`
	ImpliedVolatility float32 `json:"impliedVolatility"`
}

func (tbl TblOptionData) DropTableIfExist() error {
	if err := tbl.DropTableIfExistForTblName(TBL_OPTION_DATA_NAME); err != nil {
		return err
	}
	return tbl.DropTableIfExistForTblName(TBL_OPTION_DATA_ETF_NAME)
}

func (tbl TblOptionData) DropTableIfExistForTblName(tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + tblName)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl TblOptionData) CreateTableIfNotExist() error {
	if err := tbl.CreateTableIfNotExistForTblName(TBL_OPTION_DATA_NAME); err != nil {
		return err
	}
	return tbl.CreateTableIfNotExistForTblName(TBL_OPTION_DATA_ETF_NAME)
}

func (tbl TblOptionData) CreateTableIfNotExistForTblName(tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string
	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		tblName +
		" (" +
		//Contract symbol is Symbol(6) + yymmdd(6) + P/C(1) + Price(8) = VARCHAR(21)
		"ContractSymbol VARCHAR(21) NOT NULL," +
		"Date INT NOT NULL," +
		"Symbol VARCHAR(10) NOT NULL," +
		"OptionType VARCHAR(1) NOT NULL," +
		"Strike FLOAT(10,2)," +
		"LastPrice FLOAT(10,2)," +
		"PriceChange FLOAT(10,2)," +
		"PercentChange FLOAT(10,2)," +
		"Volume INT," +
		"PrevVolume INT," +
		"OpenInterest INT," +
		"Bid FLOAT(10,2)," +
		"Ask FLOAT(10,2)," +
		"Expiration BIGINT," +
		"ImpliedVolatility FLOAT(10,2)," +
		"InTheMoney BOOLEAN," +
		"UpdatedTime TIMESTAMP," +
		"PRIMARY KEY (ContractSymbol, Date)" +
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

func (tbl TblOptionData) InsertOrUpdateOptionData(option YahooOption, isEtf bool) error {
	if isEtf {
		return tbl.InsertOrUpdateOptionDataToTbl(option, TBL_OPTION_DATA_ETF_NAME)
	} else {
		return tbl.InsertOrUpdateOptionDataToTbl(option, TBL_OPTION_DATA_NAME)
	}
}

func (tbl TblOptionData) InsertOrUpdateOptionDataToTbl(option YahooOption, tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			tblName+
			" WHERE ContractSymbol=? AND Date=?)",
		option.ContractSymbol, GetTimeInYYYYMMDD())
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		return dbQueryErr
	}
	if exist == 0 {
		if err := tbl.InsertOptionData(option, tblName); err != nil {
			return err
		}
	} else {
		if err := tbl.UpdateOptionData(option, tblName); err != nil {
			return err
		}
	}
	return nil
}

func (tbl TblOptionData) InsertOptionData(option YahooOption, tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		tblName +
		" (" +
		"ContractSymbol, " +
		"Date, " +
		"Symbol, " +
		"OptionType, " +
		"Strike, " +
		"LastPrice, " +
		"PriceChange, " +
		"PercentChange, " +
		"Volume, " +
		"PrevVolume, " +
		"OpenInterest, " +
		"Bid, " +
		"Ask, " +
		"Expiration, " +
		"ImpliedVolatility, " +
		"InTheMoney, " +
		"UpdatedTime) VALUES (?,?,?,?,?,?,?,?,?,0,?,?,?,?,?,?,?)")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		option.ContractSymbol,
		GetTimeInYYYYMMDD(),
		option.GetSymbol(),
		option.GetOptionType(),
		option.Strike,
		option.LastPrice,
		option.PriceChange,
		0,
		//option.PercentChange,
		option.Volume,
		option.OpenInterest,
		option.Bid,
		option.Ask,
		option.Expiration,
		option.ImpliedVolatility,
		option.InTheMoney,
		GetTime())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func (tbl TblOptionData) UpdateOptionData(option YahooOption, tblName string) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	prevVolume, selectErr := tbl.SelectOptionDataVolumeByContractSymbolAndDate(option.ContractSymbol, GetTimeInYYYYMMDD())
	if selectErr != nil {
		return selectErr
	}
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		tblName +
		" SET " +
		"LastPrice=?, " +
		"PriceChange=?, " +
		"PercentChange=?, " +
		"Volume=?, " +
		"PrevVolume=?, " +
		"OpenInterest=?, " +
		"Bid=?, " +
		"Ask=?, " +
		"ImpliedVolatility=?, " +
		"InTheMoney=?, " +
		"UpdatedTime=? " +
		"WHERE ContractSymbol=? AND Date=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		option.LastPrice,
		option.PriceChange,
		0,
		//option.PercentChange,
		option.Volume,
		prevVolume,
		option.OpenInterest,
		option.Bid,
		option.Ask,
		option.ImpliedVolatility,
		option.InTheMoney,
		GetTime(),
		option.ContractSymbol,
		GetTimeInYYYYMMDD())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

func (tbl TblOptionData) SelectExpirationBySymbolAndDate(symbol string, date int, isEtf bool) ([]int64, error) {
	var expDates []int64
	var tblName string
	if isEtf {
		tblName = TBL_OPTION_DATA_ETF_NAME
	} else {
		tblName = TBL_OPTION_DATA_NAME
	}
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return expDates, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Expiration FROM " + tblName + " WHERE " + "Symbol=? AND Date=? GROUP BY Expiration")
	if dbPrepErr != nil {
		return expDates, dbPrepErr
	}
	defer stmt.Close()
	rows, dbQueryErr := stmt.Query(symbol, date)
	if dbQueryErr != nil {
		return expDates, dbQueryErr
	}
	defer rows.Close()
	var expDate int64
	for rows.Next() {
		if scanErr := rows.Scan(&expDate); scanErr != nil {
			return expDates, scanErr
		}
		expDates = append(expDates, expDate)
	}
	return expDates, nil
}

func (tbl TblOptionData) SelectOptionDataVolumeByContractSymbolAndDate(contractSymbol string, date int) (int, error) {
	return tbl.SelectOptionDataVolumeByContractSymbolAndDateFromTbl(contractSymbol, date, TBL_OPTION_DATA_NAME)
}

func (tbl TblOptionData) SelectOptionDataVolumeByContractSymbolAndDateToEtf(contractSymbol string, date int) (int, error) {
	return tbl.SelectOptionDataVolumeByContractSymbolAndDateFromTbl(contractSymbol, date, TBL_OPTION_DATA_ETF_NAME)
}

func (tbl TblOptionData) SelectOptionDataVolumeByContractSymbolAndDateFromTbl(contractSymbol string, date int, tblName string) (int, error) {
	var volume int
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return volume, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Volume FROM " + tblName + " WHERE " + "ContractSymbol=? AND Date=?")
	if dbPrepErr != nil {
		return volume, dbPrepErr
	}
	defer stmt.Close()
	optionRows, dbQueryErr := stmt.Query(contractSymbol, date)
	if dbQueryErr != nil {
		return volume, dbQueryErr
	}
	defer optionRows.Close()
	for optionRows.Next() {
		if scanErr := optionRows.Scan(&volume); scanErr != nil {
			fmt.Println(scanErr)
			return volume, scanErr
		}
	}
	return volume, nil
}

func (tbl *TblOptionData) SelectContractSymbolListBySymbolAndDate(symbol string, date int) ([]string, error) {
	var contractSymbol []string
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return contractSymbol, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT ContractSymbol FROM " + TBL_OPTION_DATA_NAME + " WHERE " + "Symbol=? AND Date>? GROUP BY ContractSymbol")
	if dbPrepErr != nil {
		return contractSymbol, dbPrepErr
	}
	defer stmt.Close()
	optionRows, dbQueryErr := stmt.Query(symbol, date)
	if dbQueryErr != nil {
		return contractSymbol, dbQueryErr
	}
	defer optionRows.Close()
	var cs string
	for optionRows.Next() {
		if scanErr := optionRows.Scan(&cs); scanErr != nil {
			fmt.Println(scanErr)
			return contractSymbol, scanErr
		}
		contractSymbol = append(contractSymbol, cs)
	}
	return contractSymbol, nil
}

func (tbl *TblOptionData) SelectContractSymbolDataByContractSymbol(contractSymbol string) ([]RowOptionData, error) {
	var rows []RowOptionData
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return rows, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Date, Symbol, OptionType, Strike, LastPrice, Volume, OpenInterest, ImpliedVolatility FROM " + TBL_OPTION_DATA_NAME + " WHERE " + "ContractSymbol=? ORDER BY Date ASC")
	if dbPrepErr != nil {
		return rows, dbPrepErr
	}
	defer stmt.Close()
	optionRows, dbQueryErr := stmt.Query(contractSymbol)
	if dbQueryErr != nil {
		return rows, dbQueryErr
	}
	defer optionRows.Close()
	var date int
	var symbol string
	var optionType string
	var strike float32
	var lastPrice float32
	var volume int
	var openInterest int
	var impliedVolatility float32
	for optionRows.Next() {
		if scanErr := optionRows.Scan(&date, &symbol, &optionType, &strike, &lastPrice, &volume, &openInterest, &impliedVolatility); scanErr != nil {
			fmt.Println(scanErr)
			return rows, scanErr
		}
		rows = append(rows, RowOptionData{
			ContractSymbol:    contractSymbol,
			Date:              date,
			Symbol:            symbol,
			OptionType:        optionType,
			Strike:            strike,
			LastPrice:         lastPrice,
			Volume:            volume,
			OpenInterest:      openInterest,
			ImpliedVolatility: impliedVolatility,
		})
	}
	return rows, nil
}
