package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	stocks "stocks-chat/bot/services"
	"strings"
)

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Stocks bot for Jobsity!")
}

func getStock(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToLower(r.URL.Query().Get("ticker"))
	if !stocks.GetStockfromAPI(ticker) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func Start() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/health", health).Methods("GET")
	router.HandleFunc("/stock", getStock).Methods("GET")
	log.Println("stocks bot listening on", "localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
