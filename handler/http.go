package handler

import (
    "github.com/gin-gonic/gin"
	"strconv"
	"net/http"
    "order-service/repository"

)

type HTTPHandler struct {
    orderRepo *repository.OrderRepository
}

func NewHTTPHandler(orderRepo *repository.OrderRepository) *HTTPHandler {
    return &HTTPHandler{orderRepo: orderRepo}
}


func (h *HTTPHandler) GetOrders(c *gin.Context) {
		pageStr := c.Query("page")
		searchQuery := c.Query("search")
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}
	
		orders, err := h.orderRepo.GetOrders(searchQuery, startDateStr, endDateStr, pageStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	
		c.JSON(http.StatusOK, orders)
}
