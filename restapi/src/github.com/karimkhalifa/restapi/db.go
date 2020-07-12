package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connPostgres() *sql.DB {
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
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + mongoURI + "/"))
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

func getExtraInfo(db *sql.DB, param1 string) rune {
	sqlStatement := `SELECT customer_id, order_items.product, order_name, created_at as order_date, sum(order_items.quantity*order_items.price_per_unit) as total_amount, sum(deliveries.delivered_quantity*order_items.price_per_unit) as delivered_amount
						FROM orders, order_items, deliveries
						WHERE orders.id = order_items.order_id AND deliveries.order_item_id = order_items.id
						GROUP BY order_name, order_date, customer_id, order_items.product`
	if param1 == "" {
		sqlStatement = `SELECT count(*), sum(total_amount)
		FROM (` + sqlStatement + `) AS subquery;`
	} else {
		sqlStatement = `SELECT count(*), sum(total_amount)
		FROM (` + sqlStatement + `) AS subquery
		WHERE ` + param1 + `;`
	}
	defer db.Close()

	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println("Failed to run query", err)
	}
	count := ' '
	rows.Next()
	rows.Scan(&count)
	return count
}

func getRecords(db *sql.DB, param1 string, param2 int) []map[string]string {

	sqlStatement := `SELECT customer_id, order_items.product, order_name, created_at as order_date, sum(order_items.quantity*order_items.price_per_unit) as total_amount, sum(deliveries.delivered_quantity*order_items.price_per_unit) as delivered_amount
				   FROM orders, order_items, deliveries
				   WHERE orders.id = order_items.order_id AND deliveries.order_item_id = order_items.id
				   GROUP BY order_name, order_date, customer_id, order_items.product
				   ORDER BY order_name ASC`
	if param1 == "" {
		sqlStatement = `SELECT customer_id, order_name, order_date, sum(total_amount) as total_amount, sum(delivered_amount) as delivered_amount
		FROM (` + sqlStatement + `) AS subquery
		GROUP BY order_name, order_date, customer_id
		ORDER BY order_date ASC
		LIMIT 5
		OFFSET ` + strconv.Itoa(param2) + `;`
	} else {
		sqlStatement = `SELECT customer_id, order_name, order_date, sum(total_amount) as total_amount, sum(delivered_amount) as delivered_amount
		FROM (` + sqlStatement + `) AS subquery
		WHERE ` + param1 + `
		GROUP BY order_name, order_date, customer_id
		ORDER BY order_date ASC
		LIMIT 5
		OFFSET ` + strconv.Itoa(param2) + `;`
	}

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
		}
		if err != nil {
			fmt.Println("Failed to add amount", err)
			return nil
		}

		result["customer_name"], result["customer_company_name"] = getMongoRecords(result["customer_id"])
		delete(result, "customer_id")
		results = append(results, result)
	}
	return results
}
