package controller

import (
	"go-final-cross/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CustomerController(router *gin.Engine, db *gorm.DB) {
	customers := router.Group("/customers")
	{
		customers.POST("/", createCustomer(db))
		customers.GET("/:id", getCustomer(db))
		customers.PUT("/:id/change-password", changePassword(db))
		customers.PUT("/:id/change-address", changeAddress(db)) // เพิ่ม route นี้
		customers.GET("/:id/carts", getCustomerCarts(db))       // เพิ่ม route นี้
		// สามารถเพิ่ม routes อื่นๆ ได้ที่นี่
	}
}

func createCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer model.Customer
		if err := c.ShouldBindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// BeforeCreate hook จะทำการ hash password ให้อัตโนมัติ
		if err := db.Create(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

		c.JSON(http.StatusCreated, response)
	}
}

func getCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var customer model.Customer
		if err := db.First(&customer, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
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

		c.JSON(http.StatusOK, response)
	}
}

func changePassword(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึง customer ID จาก URL
		customerID := c.Param("id")

		// 2. รับข้อมูลจาก request
		var changePassRequest struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}

		if err := c.ShouldBindJSON(&changePassRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 3. ค้นหาลูกค้าในฐานข้อมูล
		var customer model.Customer
		if err := db.First(&customer, customerID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		// 4. ตรวจสอบรหัสผ่านเก่า
		if !customer.CheckPassword(changePassRequest.OldPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
			return
		}

		// 5. อัปเดตรหัสผ่านใหม่ (BeforeSave hook จะทำการ hash ให้อัตโนมัติ)
		customer.Password = changePassRequest.NewPassword

		// 6. บันทึกลงฐานข้อมูล
		if err := db.Save(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		// 7. ส่ง response กลับ
		c.JSON(http.StatusOK, gin.H{
			"message":     "Password changed successfully",
			"customer_id": customer.CustomerID,
		})
	}
}

func changeAddress(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึง customer ID จาก URL
		customerID := c.Param("id")

		// 2. รับข้อมูลที่อยู่ใหม่จาก request
		var changeAddressRequest struct {
			NewAddress string `json:"new_address" binding:"required"`
		}

		if err := c.ShouldBindJSON(&changeAddressRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 3. ค้นหาลูกค้าในฐานข้อมูล
		var customer model.Customer
		if err := db.First(&customer, customerID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		// 4. อัปเดตที่อยู่ใหม่
		customer.Address = changeAddressRequest.NewAddress

		// 5. บันทึกลงฐานข้อมูล
		if err := db.Save(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
			return
		}

		// 6. สร้าง response โดยไม่รวมรหัสผ่าน
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

		// 7. ส่ง response กลับ
		c.JSON(http.StatusOK, gin.H{
			"message":  "Address changed successfully",
			"customer": response,
		})
	}
}

func getCustomerCarts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึง customer ID จาก URL
		customerID := c.Param("id")

		// 2. ค้นหาลูกค้าในฐานข้อมูลเพื่อตรวจสอบว่ามีอยู่จริง
		var customer model.Customer
		if err := db.First(&customer, customerID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		// 3. ดึงข้อมูลรถเข็นทั้งหมดของลูกค้า พร้อมรายการสินค้า
		var carts []model.Cart
		if err := db.Preload("CartItems.Product").
			Where("customer_id = ?", customerID).
			Find(&carts).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carts"})
			return
		}

		// 4. สร้างโครงสร้าง response
		type CartItemResponse struct {
			CartItemID  uint    `json:"cart_item_id"`
			ProductID   uint    `json:"product_id"`
			ProductName string  `json:"product_name"`
			Quantity    int     `json:"quantity"`
			Price       float64 `json:"price"`
			Subtotal    float64 `json:"subtotal"`
		}

		type CartResponse struct {
			CartID    uint               `json:"cart_id"`
			CartName  string             `json:"cart_name"`
			CreatedAt time.Time          `json:"created_at"`
			UpdatedAt time.Time          `json:"updated_at"`
			Items     []CartItemResponse `json:"items"`
			Total     float64            `json:"total"`
		}

		// 5. แปลงข้อมูลให้อยู่ในรูปแบบ response
		var response []CartResponse
		for _, cart := range carts {
			cartResponse := CartResponse{
				CartID:    cart.CartID,
				CartName:  cart.CartName,
				CreatedAt: cart.CreatedAt,
				UpdatedAt: cart.UpdatedAt,
				Items:     []CartItemResponse{},
				Total:     0,
			}

			var total float64 = 0
			for _, item := range cart.CartItems {
				subtotal := item.Product.Price * float64(item.Quantity)
				cartItem := CartItemResponse{
					CartItemID:  item.CartItemID,
					ProductID:   item.ProductID,
					ProductName: item.Product.ProductName,
					Quantity:    item.Quantity,
					Price:       item.Product.Price,
					Subtotal:    subtotal,
				}
				cartResponse.Items = append(cartResponse.Items, cartItem)
				total += subtotal
			}
			cartResponse.Total = total
			response = append(response, cartResponse)
		}

		// 6. ส่ง response กลับ
		c.JSON(http.StatusOK, gin.H{
			"customer_id": customerID,
			"carts":       response,
		})
	}
}
