package controller

import (
	"errors"
	"go-final-cross/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CartController(router *gin.Engine, db *gorm.DB) {
	carts := router.Group("/carts")
	{
		carts.POST("/add-item", addItemToCart(db))
		// ... routes อื่นๆ ...
	}
}

func addItemToCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			CustomerID uint   `json:"customer_id" binding:"required"`
			CartName   string `json:"cart_name" binding:"required"`
			ProductID  uint   `json:"product_id" binding:"required"`
			Quantity   int    `json:"quantity" binding:"required,min=1"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// เริ่ม transaction
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// 1. ค้นหาหรือสร้างรถเข็น
		var cart model.Cart
		if err := tx.Where("customer_id = ? AND cart_name = ?", request.CustomerID, request.CartName).
			First(&cart).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// สร้างรถเข็นใหม่
				cart = model.Cart{
					CustomerID: request.CustomerID,
					CartName:   request.CartName,
				}
				if err := tx.Create(&cart).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
					return
				}
			} else {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// 2. ตรวจสอบว่ามีสินค้านี้ในรถเข็นแล้วหรือไม่
		var existingItem model.CartItem
		if err := tx.Where("cart_id = ? AND product_id = ?", cart.CartID, request.ProductID).
			First(&existingItem).Error; err == nil {

			// ถ้ามีอยู่แล้ว ให้อัปเดตจำนวน
			existingItem.Quantity += request.Quantity
			if err := tx.Save(&existingItem).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
				return
			}
		} else {
			// ถ้ายังไม่มี ให้สร้างใหม่
			newItem := model.CartItem{
				CartID:    cart.CartID,
				ProductID: request.ProductID,
				Quantity:  request.Quantity,
			}
			if err := tx.Create(&newItem).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
				return
			}
		}

		// 3. อัปเดตข้อมูลสินค้า (ลดสต็อก)
		if err := tx.Model(&model.Product{}).
			Where("product_id = ?", request.ProductID).
			Update("stock_quantity", gorm.Expr("stock_quantity - ?", request.Quantity)).Error; err != nil {

			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
			return
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// โหลดข้อมูลรถเข็นใหม่พร้อมรายการสินค้า
		var updatedCart model.Cart
		if err := db.Preload("CartItems.Product").
			First(&updatedCart, cart.CartID).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedCart)
	}
}
