package handler

import (
    "github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"fmt"
	"net/url"
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
		startDateStr := c.Query("startDate")
		endDateStr := c.Query("endDate")
		sortDirection := c.DefaultQuery("sortDirection", "ASC")

		decodedSearchQuery, err := url.QueryUnescape(searchQuery)
		if err != nil {
			fmt.Println("Error decoding search term:", err)
			return
		}
	
		fmt.Println("Decoded search term:", decodedSearchQuery)
	
		page, err := strconv.Atoi(pageStr)

		if err != nil || page <= 0 {
			page = 1
		}
	
		orders, totalPages, err := h.orderRepo.GetOrders(searchQuery, startDateStr, endDateStr, sortDirection, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"orders": orders, "currentPage": page, "totalPages": totalPages})
}
