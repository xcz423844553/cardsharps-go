package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// func producer(nums ...int) <-chan int {
// 	out := make(chan int, 13)
// 	go func() {
// 		for _, n := range nums {
// 			fmt.Println("Store ", n)
// 			out <- n
// 		}
// 		defer close(out)
// 		fmt.Println("Closed out ")
// 	}()
// 	return out
// }

// func consumer(in <-chan int) chan int {
// 	out := make(chan int)
// 	go func() {
// 		for n := range in {
// 			time.Sleep(3 * time.Second)
// 			fmt.Println("Ready ", n)
// 			out <- n * n
// 			fmt.Println("Output ", n)
// 		}
// 		close(out)
// 	}()
// 	return out
// }
// func main() {
// 	in := producer(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
// 	consChan1 := consumer(in)
// 	//consChan2 := consumer(in)
// 	// for n := range consChan1 {
// 	// 	fmt.Println("Received ", n)
// 	// }
// 	time.Sleep(20 * time.Second)
// 	close(consChan1)
// 	return
// }
func main() {
	fmt.Println("started-service @ " + time.Now().String())

	// initDb()
	// // return
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

	ticker := time.NewTicker(TICKER_VOLUME_CHECKER * 10)
	checker := new(Checker).InitChecker()
	for ; true; <-ticker.C {
		// if !isMarketOpen() && !BYPASS_MARKET_STATUS {
		// 	fmt.Println("Market not open yet")
		// 	continue
		// }
		checker.runChecker4()
	}

	// fmt.Println("started-service @ " + time.Now().Format("20060102"))
	// new(Orbit).runOptionReportForAllSymbol(20190510)
	// log.Fatal(http.ListenAndServe(":9999", nil))
	// return

	// go new(Showdown).runShowdown()
	// new(Showdown).runShowdown2() //use this
	fmt.Println("ended-service @ " + time.Now().String())
	r := mux.NewRouter()
	r.Handle("/post", http.HandlerFunc(handlerPost))
	r.Handle("/search", http.HandlerFunc(handlerSearch))
	r.Handle("/cluster", http.HandlerFunc(handlerCluster))
	r.Handle("/signup", http.HandlerFunc(handlerSignup))
	r.Handle("/login", http.HandlerFunc(handlerLogin))
	//Backend Endpoints
	http.Handle("/", r)
	//FrontEnd endpoints
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func handlerPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one post request")
	w.Header().Set("Content-Type", "application/json")
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
