package repository

import (
    "database/sql"
	"fmt"
	"strconv"
    model "order-service/model"
)

type OrderRepository struct {
    db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
    return &OrderRepository{db: db}
}

// GetOrders retrieves orders based on search, filtering, and pagination
func (repo *OrderRepository) GetOrders(searchTerm, startDate, endDate, sortDirection, page string) ([]model.JoinedOrder, error) {
	// Set a default value for page if it's empty
	if page == "" {
		page = "1"
	}

	fmt.Println("got here")

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		fmt.Println("Error converting page string to integer:", err)
		return nil, err
	}
	
	// Calculate offset for pagination
    offset := (pageInt - 1) * 5

    // Build the SQL query
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

		var args []interface{}

	// Add search condition if searchTerm is provided
	if searchTerm != "" {
		sqlQuery += " WHERE (o.order_name ILIKE $1 OR oi.product ILIKE $1)"
		args = append(args, "%"+searchTerm+"%")
	}

	// Add date range condition if startDate and endDate are provided
	if startDate != "" {
		if searchTerm == "" {
			sqlQuery += " WHERE"
		} else {
			sqlQuery += " AND"
		}
		sqlQuery += " o.created_at >= $2"
		args = append(args, startDate)
	}
	if endDate != "" {
		if searchTerm == "" && startDate == "" {
			sqlQuery += " WHERE"
		} else {
			sqlQuery += " AND"
		}
		sqlQuery += " o.created_at <= $3"
		args = append(args, endDate)
	}

	// Add pagination
	sqlQuery += fmt.Sprintf("GROUP BY o.id, cc.company_name, c.name ORDER BY o.created_at %s LIMIT 5 OFFSET %d ", sortDirection, offset)
	fmt.Println(sqlQuery)
    // Execute the query
    rows, err := repo.db.Query(sqlQuery, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

	fmt.Println(rows)

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
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
