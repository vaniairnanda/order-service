package repository

import (
    "database/sql"
	"fmt"
	"time"
    model "order-service/model"
)

type OrderRepository struct {
    db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
    return &OrderRepository{db: db}
}

// GetOrders retrieves orders based on search, filtering, and pagination
func (repo *OrderRepository) GetOrders(searchTerm, startDate, endDate, sortDirection string, page int) ([]model.JoinedOrder, int, error) {
	// Set a default value for page if it's empty
	if page == 0 {
		page = 1
	}

	// Calculate offset for pagination
	offset := (page - 1) * 5

	var args []interface{}

	// Start building the SQL query
	sqlQuery := `
		SELECT 
			o.id,
			o.order_name,
			cc.company_name AS customer_company_name,
			c.name AS customer_name,
			o.created_at AS order_date,
			COALESCE(SUM(d.delivered_quantity * oi.price_per_unit), 0) AS delivered_amount,
			COALESCE(SUM(oi.quantity * oi.price_per_unit), 0) AS total_amount
		FROM 
			orders o
		LEFT JOIN 
			customers c ON o.customer_id = c.user_id
		LEFT JOIN 
			customer_companies cc ON c.company_id = cc.company_id
		LEFT JOIN 
			order_items oi ON o.id = oi.order_id
		LEFT JOIN
			deliveries d ON oi.id = d.order_item_id`

	var location *time.Location

	// Specify the location for Melbourne, Australia
	location, err := time.LoadLocation("Australia/Melbourne")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return nil, 0, err
	}

	// Add date range condition if startDate and endDate are provided
	if startDate != "" {
		parsedStartDate, err := time.ParseInLocation("2006-01-02", startDate, location)
		if err != nil {
			fmt.Println("Error parsing startDate:", err)
			return nil, 0, err
		}
		// Set the time to the start of the day (00:00:00)
		parsedStartDate = parsedStartDate.Truncate(24 * time.Hour)

		sqlQuery += " WHERE o.created_at >= CAST($1 AS TIMESTAMP)"
		args = append(args, parsedStartDate)
	}

	if endDate != "" {
		// Parse endDate as end of the day in Melbourne time
		parsedEndDate, err := time.ParseInLocation("2006-01-02", endDate, location)
		if err != nil {
			fmt.Println("Error parsing endDate:", err)
			return nil, 0, err
		}
		// Set the time to the end of the day (23:59:59)
		parsedEndDate = parsedEndDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

		if startDate == "" {
			sqlQuery += " WHERE"
		} else {
			sqlQuery += " AND"
		}

		sqlQuery += " o.created_at <= CAST($2 AS TIMESTAMP)"
		args = append(args, parsedEndDate)
	}

	// Add search condition if searchTerm is provided
	if searchTerm != "" {
		if startDate == "" && endDate == "" {
			sqlQuery += " WHERE"
			sqlQuery += " (o.order_name ILIKE $1 OR oi.product ILIKE $1)"
		} else {
			sqlQuery += " AND"
			sqlQuery += " (o.order_name ILIKE $3 OR oi.product ILIKE $3)"
		}
		args = append(args, "%"+searchTerm+"%")
	}

	totalRecords, err := repo.countTotalRecords(sqlQuery, args)
	if err != nil {
		return nil, 0, err
	}

	// Calculate the total number of pages
	totalPages := (totalRecords + 5 - 1) / 5

	// Add pagination
	sqlQuery += fmt.Sprintf(" GROUP BY o.id, cc.company_name, c.name ORDER BY o.created_at %s LIMIT 5 OFFSET %d", sortDirection, offset)
	fmt.Println(sqlQuery)

	// Execute the query
	rows, err := repo.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Parse results
	var orders []model.JoinedOrder
	for rows.Next() {
		var order model.JoinedOrder
		err := rows.Scan(
			&order.ID,
			&order.OrderName,
			&order.CustomerCompanyName,
			&order.CustomerName,
			&order.OrderDate,
			&order.DeliveredAmount,
			&order.TotalAmount,
		)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	return orders, totalPages, nil
}

func (repo *OrderRepository) countTotalRecords(sqlQuery string, args []interface{}) (int, error) {
	sqlQuery += fmt.Sprintf(" GROUP BY o.id, cc.company_name, c.name")

	countQuery := "SELECT COUNT(*) FROM (" + sqlQuery + ") AS count_query"

	var totalRecords int
	err := repo.db.QueryRow(countQuery, args...).Scan(&totalRecords)
	if err != nil {
		return 0, err
	}

	return totalRecords, nil
}