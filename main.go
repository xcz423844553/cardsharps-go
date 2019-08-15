package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	// new(Monitor).MonitorPCR("BA", int64(20190708))
	// new(Monitor).MonitorPCR("BA", int64(20190709))
	// new(Monitor).MonitorPCR("BA", int64(20190710))
	// new(Monitor).MonitorPCR("BA", int64(20190711))
	// new(Monitor).MonitorPCR("BA", int64(20190712))
	// new(Monitor).MonitorPCR("BA", int64(20190715))
	// new(Monitor).MonitorPCR("BA", int64(20190716))
	// new(Monitor).MonitorPCR("BA", int64(20190717))
	// new(Monitor).MonitorPCR("BA", int64(20190718))
	// new(Monitor).MonitorPCR("BA", int64(20190719))
	// return

	// initDb()
	// return

	// new(Shuffler).RecoverHistoricalStockDataFromYahoo(20190708, 20190807)
	// return

	fmt.Println("started-service @ " + time.Now().String())

	// new(Simulator).RunSimulator()
	// return
	//symbols := []string{"MSFT", "AMZN", "AAPL", "GOOGL", "GOOG", "FB", "JPM", "JNJ", "XOM", "V", "WMT", "BAC", "PG", "MA", "SPY"}

	//Download csv from IEX
	// symbols, symbolSelectErr := new(TblSymbol).SelectSymbolByFilter()
	// if symbolSelectErr != nil {
	// 	new(TblLogError).InsertLogError(LOGTYPE_DB_SYMBOL, symbolSelectErr.Error())
	// }
	// for _, symbol := range symbols {
	// 	charts, _ := new(IexApiManager).GetStockDayChartBySymbolAndRange(symbol, "2y")
	// 	new(IexApiManager).WriteIexDayChartToCsvFile(symbol, charts)
	// }
	// return

	// new(Checker).runChecker()
	// new(Checker).runChecker2()
	// go new(Checker).runChecker3()
	// go new(Checker).runSpyChecker("SPY")

	// ticker := time.NewTicker(TICKER_VOLUME_CHECKER * 10)
	// for ; true; <-ticker.C {
	// 	if !isMarketOpen() && !BYPASS_MARKET_STATUS {
	// 		fmt.Println("Market not open yet")
	// 		continue
	// 	}
	// new(Monitor).MonitorAllStock()
	// }
	// return

	/*****
	ticker := time.NewTicker(TICKER_VOLUME_CHECKER * 10)
	for true {
		for ; !isMarketOpen(); <-ticker.C {
			fmt.Println("Waiting for market to open")
		}
		for ; isMarketOpen(); <-ticker.C {
			new(Monitor).MonitorAllStock()
		}
		new(Showdown).runShowdown2()
	}
	/******/

	// fmt.Println("started-service @ " + time.Now().Format("20060102"))
	// new(Orbit).runOptionReportForAllSymbol(20190510)
	// log.Fatal(http.ListenAndServe(":9999", nil))
	// return

	// go new(Showdown).runShowdown()
	new(Monitor).MonitorAllStock()
	new(Showdown).runShowdown2() //use this
	fmt.Println("ended-service @ " + time.Now().String())
	r := mux.NewRouter()
	r.Handle("/post", http.HandlerFunc(handlerPost))
	r.Handle("/search", http.HandlerFunc(handlerSearch))
	r.Handle("/cluster", http.HandlerFunc(handlerCluster))
	r.Handle("/signup", http.HandlerFunc(handlerSignup))
	r.Handle("/login", http.HandlerFunc(handlerLogin))

	r.Handle("/optionList/{symbol}/{start_date}", http.HandlerFunc(handlerOptionList)).Methods("GET")
	r.Handle("/optionData/{contractSymbol}", http.HandlerFunc(handlerOptionData)).Methods("GET")
	r.Handle("/api/symbol_list/{action}/{symbols}", http.HandlerFunc(handlerSymbolList)).Methods("GET")
	r.Handle("/api/tag_list/{action}/{tag}", http.HandlerFunc(handlerTagList)).Methods("GET")
	r.Handle("/api/tag_list/{action}", http.HandlerFunc(handlerTagList)).Methods("GET")
	r.Handle("/api/symbol_tag_list/{action}/{symbol}/{tag}", http.HandlerFunc(handlerSymbolTagList)).Methods("GET")
	r.Handle("/api/read_all_symbol", http.HandlerFunc(handlerReadAllSymbol)).Methods("GET")
	r.Handle("/api/read_symbol_by_tag/{tags}", http.HandlerFunc(handlerReadSymbolByTags)).Methods("GET")
	r.Handle("/api/read_tag_by_symbol/{symbol}", http.HandlerFunc(handlerReadTagBySymbol)).Methods("GET")
	//Backend Endpoints
	http.Handle("/", r)
	//FrontEnd endpoints
	log.Fatal(http.ListenAndServe(":9999",
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}

func handlerPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one post request")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Write([]byte("/post reached"))
}

func handlerSearch(w http.ResponseWriter, r *http.Request) {
}

func handlerCluster(w http.ResponseWriter, r *http.Request) {
}

func handlerSignup(w http.ResponseWriter, r *http.Request) {
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
}

func handlerSymbolTagList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one symbolTagList request.")
	//Handle parameters
	params := mux.Vars(r)
	action, _ := url.PathUnescape(params["action"])
	symbol, _ := url.PathUnescape(params["symbol"])
	tag, _ := url.PathUnescape(params["tag"])

	//Sanity Check
	if action != ACTION_CREATE && action != ACTION_DELETE {
		http.Error(w, "Unknown action.", http.StatusBadRequest)
		return
	}
	if symbol == "" || tag == "" {
		http.Error(w, "Invalid symbol or tag", http.StatusBadRequest)
		return
	}

	//Initiate Variables
	var result string
	var symbolTagErr error

	//Process request
	row := TblSymbolTagRow{
		Symbol: symbol,
		Tag:    tag,
	}
	if action == ACTION_CREATE {
		symbolTagErr = new(TblSymbolTag).InsertOneRow(row)
		if symbolTagErr != nil {
			http.Error(w, "Unexpected Error", http.StatusInternalServerError)
			return
		}
		result = "Tagged symbol " + symbol + " with " + tag
	} else if action == ACTION_DELETE {
		symbolTagErr = new(TblSymbolTag).DeleteOneRow(row)
		if symbolTagErr != nil {
			http.Error(w, "Unexpected Error", http.StatusInternalServerError)
			return
		}
		result = "Remove tag " + tag + " on symbol " + symbol
	}
	//Prepare returns
	jsonResponse, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SERVER, "Json Marshal Error: "+jsonErr.Error())
		http.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResponse)
}

func handlerReadSymbolByTags(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one ReadSymbolByTag request.")
	//Handle parameters
	params := mux.Vars(r)
	tags, _ := url.PathUnescape(params["tags"])
	tagSlice := strings.Split(tags, ";")

	//Sanity Check
	if tags == "" {
		http.Error(w, "Invalid tags.", http.StatusBadRequest)
		return
	}

	//Initiate Variables
	var result []SymbolListEntity

	//Process request
	for index, tag := range tagSlice {
		row := TblSymbolTagRow{
			Tag: tag,
		}
		resultRows, symbolTagErr := new(TblSymbolTag).SelectSymbolRowByTag(row)
		if symbolTagErr != nil {
			http.Error(w, "Unexpected Error", http.StatusInternalServerError)
			return
		}
		if index == 0 {
			for _, resultRow := range resultRows {
				result = append(result, SymbolListEntity{
					Quote:   resultRow.Symbol,
					Company: resultRow.Symbol,
				})
			}
		} else {
			var tempResult []SymbolListEntity
			i := 0
			j := 0
			for i < len(result) && j < len(resultRows) {
				if strings.Compare(result[i].Quote, resultRows[j].Symbol) == 0 {
					tempResult = append(tempResult, result[i])
					i++
					j++
				} else if strings.Compare(result[i].Quote, resultRows[j].Symbol) > 0 {
					j++
				} else {
					i++
				}
			}
			result = tempResult
		}
	}

	//Prepare returns
	jsonResponse, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SERVER, "Json Marshal Error: "+jsonErr.Error())
		http.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	} else if len(result) == 0 {
		http.Error(w, "No Symbol Available", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResponse)
}

func handlerReadTagBySymbol(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one ReadTagBySymbol request.")
	//Handle parameters
	params := mux.Vars(r)
	symbol, _ := url.PathUnescape(params["symbol"])

	//Sanity Check
	if symbol == "" {
		http.Error(w, "Invalid symbol.", http.StatusBadRequest)
		return
	}

	//Initiate Variables
	var result []string

	//Process request
	row := TblSymbolTagRow{
		Symbol: symbol,
	}
	resultRows, symbolTagErr := new(TblSymbolTag).SelectTagRowBySymbol(row)
	if symbolTagErr != nil {
		http.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}
	for _, resultRow := range resultRows {
		result = append(result, resultRow.Tag)
	}

	//Prepare returns
	jsonResponse, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SERVER, "Json Marshal Error: "+jsonErr.Error())
		http.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResponse)
}

//SymbolListEntity is the entry of symbol list returned to front-end
type SymbolListEntity struct {
	Quote   string `json:"quote"`
	Company string `json:"company"`
}

//handlerReadAllSymbol returns a list of symbol with the given tag name
//Return: json formatted list of SymbolListEntity struct
func handlerReadAllSymbol(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one symbolList request.")
	var results []SymbolListEntity

	//Process request
	rows, symbolSelectErr := new(TblSymbol).SelectAllSymbolRow()
	if symbolSelectErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SERVER, "Http Handler Error: "+symbolSelectErr.Error())
		http.Error(w, "Failed to retrieve the symbols.", http.StatusBadRequest)
		return
	}
	for _, row := range rows {
		results = append(results, SymbolListEntity{
			Quote:   row.Symbol,
			Company: row.Symbol,
		})
	}

	//Prepare returns
	jsonResponse, jsonErr := json.Marshal(results)
	if jsonErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SERVER, "Json Marshal Error: "+jsonErr.Error())
		http.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResponse)
}

//handlerSymbolList returns a list of tags
//Return: json formatted list of symbol string
func handlerSymbolList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one symbolList request.")
	//Handle parameters
	params := mux.Vars(r)
	action, _ := url.PathUnescape(params["action"])
	symbols, _ := url.PathUnescape(params["symbols"])
	symbolSlice := strings.Split(symbols, ";")

	//Sanity Check
	if action != ACTION_CREATE && action != ACTION_DELETE {
		http.Error(w, "Unknown action.", http.StatusBadRequest)
		return
	}
	if symbols == "" {
		http.Error(w, "Invalid symbol.", http.StatusBadRequest)
		return
	}

	//Initiate Variables
	var result string
	var symbolErr error

	//Process request
	if action == ACTION_CREATE {
		for _, symbol := range symbolSlice {
			if symbol == "" {
				continue
			}
			row := TblSymbolRow{
				Symbol:          symbol,
				Sp500:           false,
				Nasdaq:          false,
				Dow:             false,
				Russell:         false,
				ETF:             false,
				StockMonitored:  false,
				OptionMonitored: false,
			}
			symbolErr = new(TblSymbol).InsertOneRow(row)
			if symbolErr != nil {
				result += "Symbol failed to create: " + symbol + " " + symbolErr.Error() + "\r\n"
			} else {
				result += "Symbol Created: " + symbol + "\r\n"
			}
		}
	} else if action == ACTION_DELETE {
		for _, symbol := range symbolSlice {
			if symbol == "" {
				continue
			}
			row := TblSymbolRow{
				Symbol: symbol,
			}
			symbolErr = new(TblSymbol).DeleteOneRow(row)
			if symbolErr != nil {
				result += "Symbol failed to delete: " + symbol + " " + symbolErr.Error() + "\r\n"
			} else {
				result += "Symbol Deleted: " + symbol
			}
		}
	}
	//Prepare returns
	jsonResponse, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SERVER, "Json Marshal Error: "+jsonErr.Error())
		http.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResponse)
}

//handlerTagList returns a list of tags
//Return: json formatted list of tag string
func handlerTagList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one tagList request.")
	//Handle parameters
	params := mux.Vars(r)
	action, _ := url.PathUnescape(params["action"])
	tag, _ := url.PathUnescape(params["tag"])

	//Sanity Check
	if action != ACTION_CREATE && action != ACTION_READ && action != ACTION_DELETE {
		http.Error(w, "Unknown action.", http.StatusBadRequest)
		return
	}
	if (action == ACTION_CREATE || action == ACTION_DELETE) && tag == "" {
		http.Error(w, "Invalid tag.", http.StatusBadRequest)
		return
	}

	//Initiate Variables
	var results []string
	var tagErr error

	//Process request
	if action == ACTION_CREATE {
		row := TblTagRow{
			Tag: tag,
		}
		tagErr = new(TblTag).InsertOrUpdateOneRow(row)
	} else if action == ACTION_READ {
		var rows []TblTagRow
		rows, tagErr = new(TblTag).SelectAllRow()
		for _, row := range rows {
			results = append(results, row.Tag)
		}
		if len(results) == 0 {
			http.Error(w, "No tag available.", http.StatusInternalServerError)
			return
		}
	} else if action == ACTION_DELETE {
		row := TblTagRow{
			Tag: tag,
		}
		tagErr = new(TblTag).DeleteOneRow(row)
	}
	if tagErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SERVER, "Http Handler Error: "+tagErr.Error())
		http.Error(w, "Failed to process the request.", http.StatusBadRequest)
		return
	}

	//Prepare returns
	jsonResponse, jsonErr := json.Marshal(results)
	if jsonErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SERVER, "Json Marshal Error: "+jsonErr.Error())
		http.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResponse)
}

func handlerOptionList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	symbol := params["symbol"]
	startDate, err := strconv.Atoi(params["start_date"])
	if err != nil {
		http.Error(w, "Failed to parse the start date", http.StatusBadRequest)
		return
	}
	fmt.Println("Received one optionList request for " + symbol + " " + strconv.Itoa(startDate))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	contractSymbol, err1 := new(TblOptionData).SelectContractSymbolListBySymbolAndDate(symbol, startDate)
	if err1 != nil {
		http.Error(w, "Failed to retrieve the contract symbol of "+symbol, http.StatusBadRequest)
		return
	}
	jsonResponse, jsonErr := json.Marshal(contractSymbol)
	if jsonErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SERVER, "Json Marshal Error")
		http.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func handlerOptionData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	contractSymbol := params["contractSymbol"]
	fmt.Println("Received one optionData request for " + contractSymbol)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	optionDataList, err1 := new(TblOptionData).SelectContractSymbolDataByContractSymbol(contractSymbol)
	if err1 != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SERVER, err1.Error())
		http.Error(w, "Failed to retrieve the option data of "+contractSymbol, http.StatusBadRequest)
		return
	}
	jsonResponse, jsonErr := json.Marshal(optionDataList)
	if jsonErr != nil {
		PrintMsgInConsole(MSGERROR, LOGTYPE_SERVER, "Json Marshal Error")
		http.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
