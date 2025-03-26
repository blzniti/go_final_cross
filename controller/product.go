package controller

import (
	"go-final-cross/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProductController(router *gin.Engine, db *gorm.DB) {
	products := router.Group("/products")
	{
		products.GET("/search", searchProducts(db))
		// ... routes อื่นๆ ...
	}
}

func searchProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// รับ query parameters
		keyword := c.Query("keyword")
		minPrice := c.Query("min_price")
		maxPrice := c.Query("max_price")

		query := db.Model(&model.Product{})

		if keyword != "" {
			query = query.Where("product_name LIKE ? OR description LIKE ?",
				"%"+keyword+"%", "%"+keyword+"%")
		}

		if minPrice != "" {
			query = query.Where("price >= ?", minPrice)
		}

		if maxPrice != "" {
			query = query.Where("price <= ?", maxPrice)
		}

		var products []model.Product
		if err := query.Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
