package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "karimkhalifa"
	dbname = "packform"
)

func connPostgres() *sql.DB {
	psqlInfo := fmt.Sprintf("port=%d host=%s user=%s dbname=%s sslmode=disable", port, host, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func getMongoRecords(customerID string) (string, string) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	packformDatabase := client.Database("packform")
	companiesCollection := packformDatabase.Collection("customer_companies")
	customersCollection := packformDatabase.Collection("customers")

	filterCursor1, err := customersCollection.Find(ctx, bson.M{"user_id": customerID})
	if err != nil {
		log.Fatal(err)
	}
	var customersFiltered []bson.M
	if err = filterCursor1.All(ctx, &customersFiltered); err != nil {
		log.Fatal(err)
	}

	filterCursor2, err := companiesCollection.Find(ctx, bson.M{"company_id": customersFiltered[0]["company_id"]})
	if err != nil {
		log.Fatal(err)
	}
	var companiesFiltered []bson.M
	if err = filterCursor2.All(ctx, &companiesFiltered); err != nil {
		log.Fatal(err)
	}

	return customersFiltered[0]["name"].(string), companiesFiltered[0]["company_name"].(string)
}

func getRecords(db *sql.DB, param1 string, param2 int) []map[string]string {

	sqlStatement := `SELECT order_name , created_at as order_date , customer_id , SUM(delivered_quantity*price_per_unit) as delivered_amount , SUM(price_per_unit*quantity) as total_amount
					 FROM orders , order_items , deliveries 
					 WHERE orders.id = order_id AND order_id = order_item_id ` + param1 + `
					 GROUP BY customer_id , order_name , created_at
					 ORDER BY order_name ASC 
					 LIMIT 5
					 OFFSET ` + strconv.Itoa(param2) + `;`

	defer db.Close()

	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println("Failed to run query", err)
		return nil
	}

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("Failed to get columns", err)
		return nil
	}

	rawResult := make([][]byte, len(cols))
	results := []map[string]string{}

	dest := make([]interface{}, len(cols))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			fmt.Println("Failed to scan row", err)
			return nil
		}
		result := map[string]string{}
		for i, raw := range rawResult {
			if raw == nil {
				result[cols[i]] = "0"
			} else {
				result[cols[i]] = string(raw)
			}
			println(cols[i], result[cols[i]])
		}
		t1, err := time.Parse(
			time.RFC3339,
			result["order_date"])
		if err != nil {
			fmt.Println("Failed to add amount", err)
			return nil
		}
		fmt.Println(t1)

		result["customer_name"], result["customer_company_name"] = getMongoRecords(result["customer_id"])
		delete(result, "customer_id")
		results = append(results, result)
	}
	fmt.Println(err)
	return results
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	page, _ := strconv.Atoi(mux.Vars(r)["page"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getRecords(connPostgres(), "", page*5))
}

func getOrdersLike(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	page, _ := strconv.Atoi(mux.Vars(r)["page"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getRecords(connPostgres(), "AND ( product LIKE '%"+mux.Vars(r)["search"]+"%' OR order_name LIKE '%"+mux.Vars(r)["search"]+"%')", page*5))
}

func filterOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	page, _ := strconv.Atoi(mux.Vars(r)["page"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getRecords(connPostgres(), "AND created_at BETWEEN '"+mux.Vars(r)["d1"]+"' AND '"+mux.Vars(r)["d2"]+"'", page*5))
}

func main() {
	fmt.Println("Starting...")
	r := mux.NewRouter()

	r.HandleFunc("/api/orders/{page}", getOrders).Methods("GET")
	r.HandleFunc("/api/orders/{search}/{page}", getOrdersLike).Methods("GET")
	r.HandleFunc("/api/orders/between/{d1}/{d2}/{page}", filterOrder).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
