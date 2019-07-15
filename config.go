package main

import "time"

//Create Properties.go to store private constants, including:
//  MYSQL_DBNAME   = "mysql"
//  MYSQL_USERNAME = ""
// 	MYSQL_PASSWORD = ""
// 	MYSQL_IP       = "127.0.0.1"
// 	MYSQL_PORT     = ""
// 	MYSQL_DBADDR   = MYSQL_USERNAME + ":" + MYSQL_PASSWORD + "@tcp(" + MYSQL_IP + ":" + MYSQL_PORT + ")/"

// 	EMAIL_SENDER   = "xxxxxxxxxxx@gmail.com"
// 	EMAIL_PASSWORD = "xxxxxxxxxxx"
// 	EMAIL_RECEIVER = "xxxxxxxxxxx@gmail.com"

const (
	// DB_NAME = "card_sharps_test"
	DB_NAME                  = "card_sharps"
	TBL_SYMBOL               = "symbol"
	TBL_TAG                  = "tag"
	TBL_SYMBOL_TAG           = "symbol_tag"
	TBL_OPTION_DATA_NAME     = "option_data"
	TBL_OPTION_DATA_ETF_NAME = "option_data_etf"
	TBL_STOCK_DATA_NAME      = "stock_data"
	TBL_STOCK_HIST_NAME      = "stock_hist"
	TBL_STOCK_DATA_ETF_NAME  = "stock_data_etf"
	TBL_OPTION_REPORT_NAME   = "option_report"
	TBL_STOCK_REPORT_NAME    = "stock_report"
	TBL_LOG_ERROR            = "log_error"
	TBL_LOG_SYSTEM           = "log_system"
	TBL_STOCK_CHECKER1_NAME  = "stock_checker1"
	TBL_SPY_CHECKER1_NAME    = "spy_checker1"

	URL_OPTION     = "https://query1.finance.yahoo.com/v7/finance/options/"
	URL_STOCK      = "https://query1.finance.yahoo.com/v7/finance/quote?symbols="
	URL_STOCK_HIST = "https://query1.finance.yahoo.com/v7/finance/download/%s?period1=%d&period2=%d&interval=1d&events=history&crumb=%s"

	URL_IEX_CHART_PART1 = "https://api.iextrading.com/1.0/stock/"
	URL_IEX_CHART_PART2 = "/chart/"

	URL_MACROTRENDS = "http://download.macrotrends.net/assets/php/stock_data_export.php?t="

	//Option Filter Parameter
	DEFAULT_MAX_OPTION_PERCENT  = 0.3
	DEFAULT_MIN_OPTION_PERCENT  = -0.3
	DEFAULT_MAX_OPEN_INTEREST   = 0
	DEFAULT_MIN_OPEN_INTEREST   = 0
	DEFAULT_MAX_VOLUME          = 0
	DEFAULT_MIN_VOLUME          = 0
	DEFAULT_MAX_EXPIRATION_DATE = 365 * 24 * 3600
	DEFAULT_MIN_EXPIRATION_DATE = 0

	MSGERROR  = "ERROR"
	MSGSYSTEM = "SYSTEM"

	LOGTYPE_SHOWDOWN    = "SHOWDOWN"
	LOGTYPE_DEALER      = "DEALER"
	LOGTYPE_ORBIT       = "ORBIT"
	LOGTYPE_CHECKER     = "CHECKER"
	LOGTYPE_SPY_CHECKER = "SPY_CHECKER"
	LOGTYPE_SERVER      = "SERVER"
	LOGTYPE_SHUFFLER    = "SHUFFLER"
	LOGTYPE_BOARD       = "BOARD"
	LOGTYPE_MONITOR     = "MONITOR"
	LOGTYPE_SHARPER     = "SHARPER"
	LOGTYPE_MAIN        = "MAIN"

	LOGTYPE_YAHOO_API_MANAGER  = "YAHOO_API_MANAGER"
	LOGTYPE_IEX_API_MANAGER    = "IEX_API_MANAGER"
	LOGTYPE_DB_SYMBOL          = "DB_SYMBOL"
	LOGTYPE_DB_OPTION_DATA     = "DB_OPTION_DATA"
	LOGTYPE_DB_OPTION_DATA_ETF = "DB_OPTION_DATA_ETF"
	LOGTYPE_DB_STOCK_DATA      = "DB_STOCK_DATA"
	LOGTYPE_DB_STOCK_DATA_ETF  = "DB_STOCK_DATA_ETF"
	LOGTYPE_DB_OPTION_REPORT   = "DB_OPTION_REPORT"
	LOGTYPE_DB_STOCK_REPORT    = "DB_STOCK_REPORT"

	MAX_NUM_SYMBOL_ON_EACH_GOROUTINE = 200

	MAX_LENGTH_OF_MINUTE_MONEY_CHECKER = 30
	MIN_LENGTH_OF_MINUTE_MONEY_CHECKER = 15
	MULTI_THRESHOLD_MONEY_CHECKER      = 1.5
	TICKER_MONEY_CHECKER               = 1 * time.Minute

	MAX_LENGTH_OF_MINUTE_VOLUME_CHECKER = 30
	MIN_LENGTH_OF_MINUTE_VOLUME_CHECKER = 15
	MULTI_THRESHOLD_VOLUME_CHECKER      = 2
	TICKER_VOLUME_CHECKER               = 1 * time.Minute

	UNIX_TWO_WEEK          = int64(1209600)
	SPY_CHECK_STRIKE_RANGE = 10

	BYPASS_MARKET_STATUS = false

	CHANNEL_PRODUCER_SYMBOL_BUFFER_CAPACITY = 30

	TRADER_ALL     = "All"
	TRADER_SP500   = "Sp500"
	TRADER_NASDAQ  = "Nasdaq"
	TRADER_DOW     = "Dow"
	TRADER_RUSSELL = "Russell"

	DRAGON_TAIL_LENGTH = 20
	DRAGON_TAIL_LIMIT  = 0.7
	DRAGON_TAIL_MA_GAP = 0.1
	SHIELD_HEIGHT      = 20 //#% to be considered shield

	ACTION_CREATE = "create"
	ACTION_READ   = "read"
	ACTION_UPDATE = "update"
	ACTION_DELETE = "delete"

	SYMBOLTAG_ALL        = "ALL"
	SYMBOLTAG_SP500      = "SP500"
	SYMBOLTAG_NASDAQ     = "NASDAQ"
	SYMBOLTAG_DOW        = "DOW"
	SYMBOLTAG_RUSSELL    = "RUSSELL"
	SYMBOLTAG_STOCKSTAR  = "Stock Star"
	SYMBOLTAG_OPTIONSTAR = "Option Star"
)
