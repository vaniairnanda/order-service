package main

import (
    "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
    "order-service/handler"
    "order-service/repository"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        // handle error
        panic(err)
    }

    db, err := repository.NewDB()
    if err != nil {
        // handle error
        panic(err)
    }
    defer db.Close()

    orderRepo := repository.NewOrderRepository(db)
    httpHandler := handler.NewHTTPHandler(orderRepo)

    router := gin.Default()

    router.GET("/orders", httpHandler.GetOrders)

    router.Run(":8080")
}
