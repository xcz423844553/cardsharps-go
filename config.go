package main

//"fmt"

const (
	MYSQL_DBNAME   = "mysql"
	MYSQL_USERNAME = "root"
	MYSQL_PASSWORD = "root"
	MYSQL_IP       = "127.0.0.1"
	MYSQL_PORT     = "3306"
	MYSQL_DBADDR   = MYSQL_USERNAME + ":" + MYSQL_PASSWORD + "@tcp(" + MYSQL_IP + ":" + MYSQL_PORT + ")/"

	DB_NAME = "card_sharps_test"
	// DB_NAME                  = "card_sharps"
	TBL_SYMBOL               = "symbol"
	TBL_OPTION_DATA_NAME     = "option_data"
	TBL_OPTION_DATA_ETF_NAME = "option_data_etf"
	TBL_STOCK_DATA_NAME      = "stock_data"
	TBL_OPTION_REPORT_NAME   = "option_report"
	TBL_STOCK_REPORT_NAME    = "stock_report"

	URL_OPTION = "https://query1.finance.yahoo.com/v7/finance/options/"
	URL_STOCK  = "https://query1.finance.yahoo.com/v7/finance/quote?symbols="

	//Option Filter Parameter
	DEFAULT_MAX_OPTION_PERCENT  = 0.3
	DEFAULT_MIN_OPTION_PERCENT  = -0.3
	DEFAULT_MAX_OPEN_INTEREST   = 0
	DEFAULT_MIN_OPEN_INTEREST   = 0
	DEFAULT_MAX_VOLUME          = 0
	DEFAULT_MIN_VOLUME          = 0
	DEFAULT_MAX_EXPIRATION_DATE = 365 * 24 * 3600
	DEFAULT_MIN_EXPIRATION_DATE = 0
)
