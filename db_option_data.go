package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//IOptionData is the struct of input to db_option_data
type IOptionData interface {
	GetContractSymbol() string
	GetSymbol() string
	GetOptionType() string
	GetStrike() float32
	GetLastPrice() float32
	GetVolume() int64
	GetOpenInterest() int64
	GetImpliedVolatility() float32
	GetPercentChange() float32
	GetBid() float32
	GetAsk() float32
	GetExpiration() int64
	GetLastTradeDate() int64
}

//DaoOptionData is a struct to manipulate db_option_data
type DaoOptionData struct {
}

//RowOptionData is a struct representing the row of db_option_data
type RowOptionData struct {
	ContractSymbol    string  `json:"contractSymbol"` //Contract symbol is Symbol(6) + yymmdd(6) + P/C(1) + Price(8) = VARCHAR(21)
	Date              int64   `json:"date"`           //YYYYMMDD
	Symbol            string  `json:"symbol"`
	OptionType        string  `json:"optionType"`
	Strike            float32 `json:"strike"`
	LastPrice         float32 `json:"lastPrice"`
	Volume            int64   `json:"volume"`
	OpenInterest      int64   `json:"openInterest"`
	ImpliedVolatility float32 `json:"impliedVolatility"`
	PercentChange     float32 `json:"percentChange"`
	Bid               float32 `json:"bid"`
	Ask               float32 `json:"ask"`
	Expiration        int64   `json:"expiration"`    //YYYYMMDD
	LastTradeDate     int64   `json:"lastTradeDate"` //YYYYMMDD
}

//GetTableName returns table name
func (dao *DaoOptionData) getTableName() string {
	if TestMode {
		return TblNameOptionData + "_test"
	}
	return TblNameOptionData
}

//DropTableIfExist drops table if the table exists
func (dao *DaoOptionData) DropTableIfExist() error {
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
func (dao *DaoOptionData) CreateTableIfNotExist() error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()

	sqlStr := "CREATE TABLE IF NOT EXISTS " +
		dao.getTableName() +
		" (" +
		"ContractSymbol VARCHAR(21) NOT NULL," +
		"Date INT NOT NULL," +
		"Symbol VARCHAR(10) NOT NULL," +
		"OptionType VARCHAR(1) NOT NULL," +
		"Strike FLOAT(10,2)," +
		"LastPrice FLOAT(10,2)," +
		"Volume INT," +
		"OpenInterest INT," +
		"ImpliedVolatility FLOAT(10,2)," +
		"PercentChange FLOAT(10,2)," +
		"Bid FLOAT(10,2)," +
		"Ask FLOAT(10,2)," +
		"Expiration INT," +
		"LastTradeDate INT," +
		"UpdatedTime TIMESTAMP," +
		"PRIMARY KEY (ContractSymbol, Date)" +
		")"
	_, tblCreateErr := db.Exec(sqlStr)
	if tblCreateErr != nil {
		return tblCreateErr
	}
	return nil
}

//InsertOrUpdateOptionData inserts or updates the option data
func (dao *DaoOptionData) InsertOrUpdateOptionData(option IOptionData) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	resQuery := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM "+
			dao.getTableName()+
			" WHERE ContractSymbol=? AND Date=?)",
		option.GetContractSymbol(), GetTimeInYYYYMMDD64())
	var exist int
	dbQueryErr := resQuery.Scan(&exist)
	if dbQueryErr != nil && dbQueryErr != sql.ErrNoRows {
		return dbQueryErr
	}
	if exist == 0 {
		if err := dao.InsertOptionData(option); err != nil {
			return err
		}
	} else {
		if err := dao.UpdateOptionData(option); err != nil {
			return err
		}
	}
	return nil
}

//InsertOptionData inserts the option data
func (dao *DaoOptionData) InsertOptionData(option IOptionData) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		dao.getTableName() +
		" (" +
		"ContractSymbol, " +
		"Date, " +
		"Symbol, " +
		"OptionType, " +
		"Strike, " +
		"LastPrice, " +
		"Volume, " +
		"OpenInterest, " +
		"ImpliedVolatility, " +
		"PercentChange, " +
		"Bid, " +
		"Ask, " +
		"Expiration, " +
		"LastTradeDate, " +
		"UpdatedTime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		option.GetContractSymbol(),
		GetTimeInYYYYMMDD64(),
		option.GetSymbol(),
		option.GetOptionType(),
		option.GetStrike(),
		option.GetLastPrice(),
		option.GetVolume(),
		option.GetOpenInterest(),
		option.GetImpliedVolatility(),
		option.GetPercentChange(),
		option.GetBid(),
		option.GetAsk(),
		option.GetExpiration(),
		option.GetLastTradeDate(),
		GetTime())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

//UpdateOptionData updates the option data
func (dao *DaoOptionData) UpdateOptionData(option IOptionData) error {
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("UPDATE " +
		dao.getTableName() +
		" SET " +
		"LastPrice=?, " +
		"Volume=?, " +
		"OpenInterest=?, " +
		"ImpliedVolatility=?, " +
		"PercentChange=?, " +
		"Bid=?, " +
		"Ask=?, " +
		"LastTradeDate=?, " +
		"UpdatedTime=? " +
		"WHERE ContractSymbol=? AND Date=?")
	if dbPrepErr != nil {
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		option.GetLastPrice(),
		option.GetVolume(),
		option.GetOpenInterest(),
		option.GetImpliedVolatility(),
		option.GetPercentChange(),
		option.GetBid(),
		option.GetAsk(),
		option.GetLastTradeDate(),
		GetTime(),
		option.GetContractSymbol(),
		GetTimeInYYYYMMDD64())
	if dbExecErr != nil {
		return dbExecErr
	}
	return nil
}

//SelectExpirationBySymbolAndDate returns an array of expiration dates of the given symbol and date
func (dao *DaoOptionData) SelectExpirationBySymbolAndDate(symbol string, date int64) ([]int64, error) {
	var expDates []int64
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return expDates, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Expiration FROM " + dao.getTableName() + " WHERE " + "Symbol=? AND Date=? GROUP BY Expiration")
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

//SelectOptionDataVolumeByContractSymbolAndDate returns the option volume of the given contract symbol and date
func (dao *DaoOptionData) SelectOptionDataVolumeByContractSymbolAndDate(contractSymbol string, date int64) (int64, error) {
	var volume int64
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return volume, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT Volume FROM " + dao.getTableName() + " WHERE " + "ContractSymbol=? AND Date=?")
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

//SelectContractSymbolBySymbolAndDate returns an array of contract symbol of the given symbol and date (later or equal to date)
func (dao *DaoOptionData) SelectContractSymbolBySymbolAndDate(symbol string, date int64) ([]string, error) {
	var contractSymbol []string
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return contractSymbol, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT ContractSymbol FROM " + dao.getTableName() + " WHERE " + "Symbol=? AND Date>=? GROUP BY ContractSymbol")
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

//SelectOptionDataByContractSymbol returns an array of option data of the given contract symbol
func (dao *DaoOptionData) SelectOptionDataByContractSymbol(contractSymbol string) ([]RowOptionData, error) {
	var rows []RowOptionData
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return rows, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare(
		"SELECT `Date`, Symbol, OptionType, Strike, LastPrice, Volume, OpenInterest, ImpliedVolatility, PercentChange, Bid, Ask, Expiration, LastTradeDate FROM " +
			dao.getTableName() +
			" WHERE" +
			" ContractSymbol=? ORDER BY Date ASC")
	if dbPrepErr != nil {
		return rows, dbPrepErr
	}
	defer stmt.Close()
	optionRows, dbQueryErr := stmt.Query(contractSymbol)
	if dbQueryErr != nil {
		return rows, dbQueryErr
	}
	defer optionRows.Close()
	var date int64
	var symbol string
	var optionType string
	var strike float32
	var lastPrice float32
	var volume int64
	var openInterest int64
	var impliedVolatility float32
	var percentChange float32
	var bid float32
	var ask float32
	var expiration int64
	var lastTradeDate int64
	for optionRows.Next() {
		if scanErr := optionRows.Scan(&date, &symbol, &optionType, &strike, &lastPrice, &volume, &openInterest, &impliedVolatility, &percentChange, &bid, &ask, &expiration, &lastTradeDate); scanErr != nil {
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
			PercentChange:     percentChange,
			Bid:               bid,
			Ask:               ask,
			Expiration:        expiration,
			LastTradeDate:     lastTradeDate,
		})
	}
	return rows, nil
}

//SelectExpirationVolumeOpenInterestBySymbolAndDate returns an array of expiration date, a map of call volume, a map of put volume, a map of call open interest, and a map of put open interest
func (dao *DaoOptionData) SelectExpirationVolumeOpenInterestBySymbolAndDate(symbol string, date int64) ([]int64, map[int64]int64, map[int64]int64, map[int64]int64, map[int64]int64, error) {
	var expDate []int64
	exp := make(map[int64]struct{}) //serve as a hashset for expDate
	callVol := make(map[int64]int64)
	putVol := make(map[int64]int64)
	callOi := make(map[int64]int64)
	putOi := make(map[int64]int64)
	db, dbConnErr := new(Dao).ConnectToDb()
	if dbConnErr != nil {
		return expDate, callVol, putVol, callOi, putOi, dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("SELECT OptionType, Volume, OpenInterest, Expiration FROM " + dao.getTableName() + " WHERE " + "Symbol=? AND Date=?")
	if dbPrepErr != nil {
		return expDate, callVol, putVol, callOi, putOi, dbPrepErr
	}
	defer stmt.Close()
	optionRows, dbQueryErr := stmt.Query(symbol, date)
	if dbQueryErr != nil {
		return expDate, callVol, putVol, callOi, putOi, dbQueryErr
	}
	defer optionRows.Close()
	var optionType string
	var volume int64
	var openInterest int64
	var expiration int64
	for optionRows.Next() {
		if scanErr := optionRows.Scan(&optionType, &volume, &openInterest, &expiration); scanErr != nil {
			fmt.Println(scanErr)
			return expDate, callVol, putVol, callOi, putOi, scanErr
		}
		if _, exist := exp[expiration]; !exist {
			exp[expiration] = struct{}{}
		}
		if optionType == "C" {
			callVol[expiration] += volume
			callOi[expiration] += openInterest
		} else {
			putVol[expiration] += volume
			putOi[expiration] += openInterest
		}
	}
	for key := range exp {
		expDate = append(expDate, key)
	}
	return expDate, callVol, putVol, callOi, putOi, nil
}
