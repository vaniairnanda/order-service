# order-service

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

6. Access the client `cd ..` `cd client/order-service-ui`

7. Run the client `npm run dev`
