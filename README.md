# order-service

## Getting Started

1. Copy `.env.example` to `.env` and configure the env
   
   ```
   cp .env.example .env
   ```

2. Run database via docker (or PostgreSQL in your local)

    ```bash
    docker-compose up -d
    ```

3. Seed the database using the following command
   ```
   go run script/csvparser.go
   ```

3. Run the app `go run .`
