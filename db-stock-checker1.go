package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type TblStockChecker1 struct {
}

func (tbl TblStockChecker1) DropTableIfExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	_, tblDropErr := db.Exec("DROP TABLE IF EXISTS " + TBL_STOCK_CHECKER1_NAME)
	if tblDropErr != nil {
		return tblDropErr
	}
	return nil
}

func (tbl TblStockChecker1) CreateTableIfNotExist() error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	var sqlStr string
	sqlStr = "CREATE TABLE IF NOT EXISTS " +
		TBL_STOCK_CHECKER1_NAME +
		" (" +
		"Date INT NOT NULL," +
		"Minute INT NOT NULL," +
		"MSFT_PRICE FLOAT(10,2)," +
		"MSFT_VOLUME INT," +
		"AMZN_PRICE FLOAT(10,2)," +
		"AMZN_VOLUME INT," +
		"AAPL_PRICE FLOAT(10,2)," +
		"AAPL_VOLUME INT," +
		"GOOGL_PRICE FLOAT(10,2)," +
		"GOOGL_VOLUME INT," +
		"GOOG_PRICE FLOAT(10,2)," +
		"GOOG_VOLUME INT," +
		"FB_PRICE FLOAT(10,2)," +
		"FB_VOLUME INT," +
		"JPM_PRICE FLOAT(10,2)," +
		"JPM_VOLUME INT," +
		"JNJ_PRICE FLOAT(10,2)," +
		"JNJ_VOLUME INT," +
		"XOM_PRICE FLOAT(10,2)," +
		"XOM_VOLUME INT," +
		"V_PRICE FLOAT(10,2)," +
		"V_VOLUME INT," +
		"WMT_PRICE FLOAT(10,2)," +
		"WMT_VOLUME INT," +
		"BAC_PRICE FLOAT(10,2)," +
		"BAC_VOLUME INT," +
		"PG_PRICE FLOAT(10,2)," +
		"PG_VOLUME INT," +
		"MA_PRICE FLOAT(10,2)," +
		"MA_VOLUME INT," +
		"CSCO_PRICE FLOAT(10,2)," +
		"CSCO_VOLUME INT," +
		"PFE_PRICE FLOAT(10,2)," +
		"PFE_VOLUME INT," +
		"DIS_PRICE FLOAT(10,2)," +
		"DIS_VOLUME INT," +
		"VZ_PRICE FLOAT(10,2)," +
		"VZ_VOLUME INT," +
		"T_PRICE FLOAT(10,2)," +
		"T_VOLUME INT," +
		"CVX_PRICE FLOAT(10,2)," +
		"CVX_VOLUME INT," +
		"UNH_PRICE FLOAT(10,2)," +
		"UNH_VOLUME INT," +
		"HD_PRICE FLOAT(10,2)," +
		"HD_VOLUME INT," +
		"KO_PRICE FLOAT(10,2)," +
		"KO_VOLUME INT," +
		"MRK_PRICE FLOAT(10,2)," +
		"MRK_VOLUME INT," +
		"INTC_PRICE FLOAT(10,2)," +
		"INTC_VOLUME INT," +
		"WFC_PRICE FLOAT(10,2)," +
		"WFC_VOLUME INT," +
		"ORCL_PRICE FLOAT(10,2)," +
		"ORCL_VOLUME INT," +
		"CMCSA_PRICE FLOAT(10,2)," +
		"CMCSA_VOLUME INT," +
		"PEP_PRICE FLOAT(10,2)," +
		"PEP_VOLUME INT," +
		"NFLX_PRICE FLOAT(10,2)," +
		"NFLX_VOLUME INT," +
		"MCD_PRICE FLOAT(10,2)," +
		"MCD_VOLUME INT," +
		"C_PRICE FLOAT(10,2)," +
		"C_VOLUME INT," +
		"UpdatedTime TIMESTAMP," +
		"PRIMARY KEY (Date, Minute)" +
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

func (tbl TblStockChecker1) InsertStockCheckerData(date int, minute int, volumes []int, prices []float32) error {
	db, dbConnErr := sql.Open(MYSQL_DBNAME, MYSQL_DBADDR+DB_NAME)
	if dbConnErr != nil {
		return dbConnErr
	}
	defer db.Close()
	stmt, dbPrepErr := db.Prepare("INSERT INTO " +
		TBL_STOCK_CHECKER1_NAME +
		" (" +
		"Date," +
		"Minute," +
		"MSFT_PRICE," +
		"MSFT_VOLUME," +
		"AMZN_PRICE," +
		"AMZN_VOLUME," +
		"AAPL_PRICE," +
		"AAPL_VOLUME," +
		"GOOGL_PRICE," +
		"GOOGL_VOLUME," +
		"GOOG_PRICE," +
		"GOOG_VOLUME," +
		"FB_PRICE," +
		"FB_VOLUME," +
		"JPM_PRICE," +
		"JPM_VOLUME," +
		"JNJ_PRICE," +
		"JNJ_VOLUME," +
		"XOM_PRICE," +
		"XOM_VOLUME," +
		"V_PRICE," +
		"V_VOLUME," +
		"WMT_PRICE," +
		"WMT_VOLUME," +
		"BAC_PRICE," +
		"BAC_VOLUME," +
		"PG_PRICE," +
		"PG_VOLUME," +
		"MA_PRICE," +
		"MA_VOLUME," +
		"CSCO_PRICE," +
		"CSCO_VOLUME," +
		"PFE_PRICE," +
		"PFE_VOLUME," +
		"DIS_PRICE," +
		"DIS_VOLUME," +
		"VZ_PRICE," +
		"VZ_VOLUME," +
		"T_PRICE," +
		"T_VOLUME," +
		"CVX_PRICE," +
		"CVX_VOLUME," +
		"UNH_PRICE," +
		"UNH_VOLUME," +
		"HD_PRICE," +
		"HD_VOLUME," +
		"KO_PRICE," +
		"KO_VOLUME," +
		"MRK_PRICE," +
		"MRK_VOLUME," +
		"INTC_PRICE," +
		"INTC_VOLUME," +
		"WFC_PRICE," +
		"WFC_VOLUME," +
		"ORCL_PRICE," +
		"ORCL_VOLUME," +
		"CMCSA_PRICE," +
		"CMCSA_VOLUME," +
		"PEP_PRICE," +
		"PEP_VOLUME," +
		"NFLX_PRICE," +
		"NFLX_VOLUME," +
		"MCD_PRICE," +
		"MCD_VOLUME," +
		"C_PRICE," +
		"C_VOLUME," +
		"UpdatedTime) VALUES (?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?)")
	if dbPrepErr != nil {
		//panic(dbPrepErr)
		return dbPrepErr
	}
	defer stmt.Close()
	_, dbExecErr := stmt.Exec(
		date,
		minute,
		prices[0],
		volumes[0],
		prices[1],
		volumes[1],
		prices[2],
		volumes[2],
		prices[3],
		volumes[3],
		prices[4],
		volumes[4],
		prices[5],
		volumes[5],
		prices[6],
		volumes[6],
		prices[7],
		volumes[7],
		prices[8],
		volumes[8],
		prices[9],
		volumes[9],
		prices[10],
		volumes[10],
		prices[11],
		volumes[11],
		prices[12],
		volumes[12],
		prices[13],
		volumes[13],
		prices[14],
		volumes[14],
		prices[15],
		volumes[15],
		prices[16],
		volumes[16],
		prices[17],
		volumes[17],
		prices[18],
		volumes[18],
		prices[19],
		volumes[19],
		prices[20],
		volumes[20],
		prices[21],
		volumes[21],
		prices[22],
		volumes[22],
		prices[23],
		volumes[23],
		prices[24],
		volumes[24],
		prices[25],
		volumes[25],
		prices[26],
		volumes[26],
		prices[27],
		volumes[27],
		prices[28],
		volumes[28],
		prices[29],
		volumes[29],
		prices[30],
		volumes[30],
		prices[31],
		volumes[31],
		GetTime())
	if dbExecErr != nil {
		//panic(dbExecErr)
		return dbExecErr
	}
	return nil
}
