package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	repository	"order-service/repository"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func parseCSV(filePath string) ([][]string, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	records, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}

	// Skip the first record (header)
	if len(records) > 0 {
		records = records[1:]
	}

	return records, nil
}

func createCustomerCompaniesTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS customer_companies (
			id SERIAL PRIMARY KEY,
			company_name VARCHAR(255)
		);
	`)
	return err
}

func createCustomersTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			user_id VARCHAR(255) PRIMARY KEY,
			login VARCHAR(255),
			password VARCHAR(255),
			name VARCHAR(255),
			company_id INT,
			credit_cards VARCHAR(255)
		);
	`)
	return err
}

func createDeliveriesTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS deliveries (
			id SERIAL PRIMARY KEY,
			order_item_id INT,
			delivered_quantity INT
		);
	`)
	return err
}

func createOrderItemsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INT,
			price_per_unit DECIMAL(10, 2),
			quantity INT,
			product VARCHAR(255)
		);
	`)
	return err
}

func createOrdersTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP,
			order_name VARCHAR(255),
			customer_id VARCHAR(255)
		);
	`)
	return err
}

func seedCustomerCompanies(db *sql.DB, filePath string) error {
	records, err := parseCSV(filePath)
	if err != nil {
		return err
	}

	for _, record := range records {
		id, _ := strconv.Atoi(string(record[0])) 
		companyName := string(record[1])

		err := createCustomerCompaniesTable(db)
		if err != nil {
			return err
		}

		err = insertCustomerCompany(db, id, companyName)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedCustomers(db *sql.DB, filePath string) error {
	records, err := parseCSV(filePath)
	if err != nil {
		return err
	}

	for _, record := range records {
		userID := string(record[0])
		login := string(record[1])
		password := string(record[2])
		name := string(record[3])
		companyID, _ := strconv.Atoi(string(record[4]))
		creditCards := string(record[5])

		err := createCustomersTable(db)
		if err != nil {
			return err
		}

		err = insertCustomer(db, userID, login, password, name, companyID, creditCards)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedDeliveries(db *sql.DB, filePath string) error {
	records, err := parseCSV(filePath)
	if err != nil {
		return err
	}

	for _, record := range records {
		id, _ := strconv.Atoi(string(record[0]))
		orderItemID, _ := strconv.Atoi(string(record[1]))
		deliveredQuantity, _ := strconv.Atoi(string(record[2])) 


		err := createDeliveriesTable(db)
		if err != nil {
			return err
		}
		err = insertDelivery(db, id, orderItemID, deliveredQuantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertCustomerCompany(db *sql.DB, id int, companyName string) error {
	_, err := db.Exec(`
		INSERT INTO customer_companies (id, company_name)
		VALUES ($1, $2)`,
		id, companyName,
	)
	return err
}

func insertCustomer(db *sql.DB, userID string, login, password, name string, companyID int, creditCards string) error {
	_, err := db.Exec(`
		INSERT INTO customers (user_id, login, password, name, company_id, credit_cards)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, login, password, name, companyID, creditCards,
	)
	return err
}

func insertDelivery(db *sql.DB, id, orderItemID, deliveredQuantity int) error {
	_, err := db.Exec(`
		INSERT INTO deliveries (id, order_item_id, delivered_quantity)
		VALUES ($1, $2, $3)`,
		id, orderItemID, deliveredQuantity,
	)
	return err
}

func seedOrderItems(db *sql.DB, filePath string) error {
	err := createOrderItemsTable(db)
	if err != nil {
		return err
	}

	err = insertOrderItemsFromCSV(db, filePath)
	if err != nil {
		return err
	}

	return nil
}

func seedOrders(db *sql.DB, filePath string) error {
	err := createOrdersTable(db)
	if err != nil {
		return err
	}

	err = insertOrdersFromCSV(db, filePath)
	if err != nil {
		return err
	}

	return nil
}

func insertOrderItemsFromCSV(db *sql.DB, filePath string) error {
	records, err := parseCSV(filePath)
	if err != nil {
		return err
	}

	for _, record := range records {
		orderID, _ := strconv.Atoi(string(record[1]))
		pricePerUnitStr := record[2]          
		pricePerUnit, _ := strconv.ParseFloat(pricePerUnitStr, 64) 
		quantity, _ := strconv.Atoi(string(record[3]))
		product := record[4]

		err := insertOrderItem(db, orderID, pricePerUnit, quantity, product)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertOrdersFromCSV(db *sql.DB, filePath string) error {
	records, err := parseCSV(filePath)
	if err != nil {
		return err
	}

	for _, record := range records {
		createdAtStr := record[1] 
		fmt.Println(createdAtStr)
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return err
		}
		orderName := record[2]
		customerID := record[3]

		err = insertOrder(db, createdAt, orderName, customerID)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertOrderItem(db *sql.DB, orderID int, pricePerUnit float64, quantity int, product string) error {
	_, err := db.Exec(`
		INSERT INTO order_items (order_id, price_per_unit, quantity, product)
		VALUES ($1, $2, $3, $4)`,
		orderID, pricePerUnit, quantity, product,
	)
	return err
}

func insertOrder(db *sql.DB, createdAt time.Time, orderName string, customerID string) error {
	_, err := db.Exec(`
		INSERT INTO orders (created_at, order_name, customer_id)
		VALUES ($1, $2, $3)`,
		createdAt.Format("2006-01-02 15:04:05"), orderName, customerID,
	)
	return err
}	


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := repository.NewDB()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	Seed data from CSV files
	err = seedCustomerCompanies(db, "./script/resources/customer_companies.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = seedCustomers(db, "./script/resources/customers.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = seedDeliveries(db, "./script/resources/deliveries.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = seedOrderItems(db, "./script/resources/order_items.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = seedOrders(db, "./script/resources/orders.csv")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data seeding completed successfully!")
}
