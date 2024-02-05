package main

import (
    "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
    "github.com/gin-contrib/cors"
    "go.uber.org/zap"
    "order-service/handler"
    "order-service/repository"
)

func main() {
    loggerConfig := zap.NewDevelopmentConfig()
	zapLogger, _ := loggerConfig.Build()
	defer zapLogger.Sync()
	zap.ReplaceGlobals(zapLogger)

    err := godotenv.Load()
    if err != nil {
        zapLogger.Error(err.Error())
    }

    db, err := repository.NewDB()
    if err != nil {
        zapLogger.Error(err.Error())
    }
    defer db.Close()

    orderRepo := repository.NewOrderRepository(db)
    httpHandler := handler.NewHTTPHandler(orderRepo)

    router := gin.Default()

    // Use CORS middleware to enable CORS
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:5173"} 
    router.Use(cors.New(config))

    router.GET("/orders", httpHandler.GetOrders)

    router.Run(":8080")
}
