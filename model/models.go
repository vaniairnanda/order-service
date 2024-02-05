package models

import "time"

type CustomerCompany struct {
	CompanyID   int    `json:"company_id"`
	CompanyName string `json:"company_name"`
}

type Customer struct {
	UserID       int    `json:"user_id"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	CompanyID    int    `json:"company_id"`
	CreditCards  string `json:"credit_cards"`
}

type Delivery struct {
	ID                int `json:"id"`
	OrderItemID       int `json:"order_item_id"`
	DeliveredQuantity int `json:"delivered_quantity"`
}

type OrderItem struct {
	ID           int     `json:"id"`
	OrderID      int     `json:"order_id"`
	PricePerUnit float64 `json:"price_per_unit"`
	Quantity     int     `json:"quantity"`
	Product      string  `json:"product"`
}

type Order struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	OrderName string    `json:"order_name"`
	CustomerID int      `json:"customer_id"`
}

type JoinedOrder struct {
    ID                  int       `json:"id"`
    OrderName           string    `json:"order_name"`
    CustomerCompanyName string    `json:"customer_company_name"`
    CustomerName        string    `json:"customer_name"`
    OrderDate           time.Time `json:"order_date"`
    DeliveredAmount     float64   `json:"delivered_amount"`
    TotalAmount         float64   `json:"total_amount"`
}