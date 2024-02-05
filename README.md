# order-service
# Order Service Web Application

This project includes a script for initializing the database and populating it with data from a CSV file and a server in Golang. Additionally, it provides a simple web application in Vue.js that displays customer order information.

## Sample API Call

Perform a GET request to retrieve a paginated list of orders with optional search, sorting, and date range parameters.

```bash
curl --location --request GET 'http://localhost:8080/orders?page=1&search=PO&sortDirection=ASC&startDate=2020-01-03&endDate=2020-01-05'
```

API Parameters:
- page: Page number for pagination.
- search: Search keyword for filtering orders.
- sortDirection: Sorting direction (ASC or DESC).
- startDate: Start date for filtering orders by date range.
- endDate: End date for filtering orders by date range.


## Sample API result
```
{
    "currentPage": 1,
    "orders": [
        {
            "id": 101,
            "order_name": "PO #101-I",
            "customer_company_name": "Sample Company A",
            "customer_name": "John Doe",
            "order_date": "2022-05-15T08:45:00Z",
            "delivered_amount": 15.75,
            "total_amount": 1520.89
        },
        {
            "id": 102,
            "order_name": "PO #102-P",
            "customer_company_name": "Sample Company B",
            "customer_name": "Jane Doe",
            "order_date": "2022-05-16T12:30:00Z",
            "delivered_amount": 298.25,
            "total_amount": 765.45
        },
        {
            "id": 103,
            "order_name": "PO #103-I",
            "customer_company_name": "Sample Company C",
            "customer_name": "Sam Smith",
            "order_date": "2022-05-18T16:20:00Z",
            "delivered_amount": 0,
            "total_amount": 1489.36
        }
    ],
    "totalPages": 2
}
```

## Getting Started

1. Access the server 
  ```cd server```

2. Copy `.env.example` to `.env` and configure the env
   
   ```
   cp .env.example .env
   ```

3. Run database via docker (or PostgreSQL in your local)

    ```bash
    docker-compose up -d
    ```

4. Seed the database using the following command
   ```
   go run script/csvparser.go
   ```

5. Run the server app `go run .`

6. Access the client `cd ../client/order-service-ui`

7. Run the client `npm run dev`
