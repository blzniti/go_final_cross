package controller

import (
	"go-final-cross/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthController(router *gin.Engine, db *gorm.DB) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", login(db))
		// สามารถเพิ่ม routes อื่นๆ เช่น register, logout ได้ที่นี่
	}
}

func login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ค้นหาลูกค้าด้วย email
		var customer model.Customer
		if err := db.Where("email = ?", loginRequest.Email).First(&customer).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// ตรวจสอบรหัสผ่าน
		if !customer.CheckPassword(loginRequest.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// สร้าง response โดยไม่รวมรหัสผ่าน
		response := model.CustomerResponse{
			CustomerID:  customer.CustomerID,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Email:       customer.Email,
			PhoneNumber: customer.PhoneNumber,
			Address:     customer.Address,
			CreatedAt:   customer.CreatedAt,
			UpdatedAt:   customer.UpdatedAt,
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Login successful",
			"customer": response,
			// สามารถเพิ่ม token JWT ได้ที่นี่หากต้องการ
		})
	}
}
