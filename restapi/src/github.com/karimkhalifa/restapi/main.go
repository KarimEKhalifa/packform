package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var psqlInfo string
var mongoURI string
var apiPort string

func allowCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	page, _ := strconv.Atoi(mux.Vars(r)["page"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getRecords(connPostgres(), "", page*5))
}

func getOrdersLike(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	page, _ := strconv.Atoi(mux.Vars(r)["page"])
	searchString := mux.Vars(r)["search"]
	fmt.Println(searchString)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getRecords(connPostgres(), "LOWER(product) LIKE LOWER('%"+searchString+"%') OR LOWER(order_name) LIKE LOWER('%"+searchString+"%')", page*5))
}

func filterOrder(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	page, _ := strconv.Atoi(mux.Vars(r)["page"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getRecords(connPostgres(), "order_date >= '"+mux.Vars(r)["d1"]+"' AND  order_date <= '"+mux.Vars(r)["d2"]+"'", page*5))
}

func filterOrderBetween(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	page, _ := strconv.Atoi(mux.Vars(r)["page"])
	searchString := mux.Vars(r)["search"]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getRecords(connPostgres(), "(LOWER(product) LIKE LOWER('%"+searchString+"%') OR LOWER(order_name) LIKE LOWER('%"+searchString+"%')) AND (order_date >= '"+mux.Vars(r)["d1"]+"' AND  order_date <= '"+mux.Vars(r)["d2"]+"')", page*5))
}

func gettingSetupParams() (string, string, string) {
	fmt.Println("Before servers start make sure that you have:")
	fmt.Println("1- MongoDB running")
	fmt.Println("2- Postgres running")
	fmt.Println("3- Ran the python script before starting the rest api, to initialize and populate the databases")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter MongoDB connection URI:Port ex. localhost:27017: ")
	mongoURI, _ := reader.ReadString('\n')
	fmt.Println("Your MongoDB URI is: " + mongoURI)

	fmt.Print("Enter Postgres connection ex. port='5432' host='localhost' user='username' dbname='mydb' sslmode=disable: ")
	postgresURI, _ := reader.ReadString('\n')
	fmt.Println("Your Postgres parameters are: " + postgresURI)

	fmt.Print("Enter the port you want to use for the api: ")
	port, _ := reader.ReadString('\n')
	fmt.Println("You chose port: " + port)

	return strings.ReplaceAll(postgresURI, "\n", ""), strings.ReplaceAll(mongoURI, "\n", ""), strings.ReplaceAll(port, "\n", "")
}

func main() {
	psqlInfo, mongoURI, apiPort = gettingSetupParams()
	fmt.Println("Starting...")

	r := mux.NewRouter()

	r.HandleFunc("/api/orders/{page}", getOrders).Methods("GET")
	r.HandleFunc("/api/orders/{search}/{page}", getOrdersLike).Methods("GET")
	r.HandleFunc("/api/orders/between/{d1}/{d2}/{page}", filterOrder).Methods("GET")
	r.HandleFunc("/api/orders/search_between/{search}/{d1}/{d2}/{page}", filterOrderBetween).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+apiPort, r))
}
