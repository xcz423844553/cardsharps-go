package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("started-service @ " + time.Now().String())
	fmt.Println("started-service @ " + time.Now().Format("20060102"))
	go runCore()
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
