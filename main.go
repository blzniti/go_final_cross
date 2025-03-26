package main

import (
	"fmt"
	"go-final-cross/controller"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.Get("mysql.dsn"))
	dsn := viper.GetString("mysql.dsn")

	dialactor := mysql.Open(dsn)
	db, err := gorm.Open(dialactor)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection successful")

	router := gin.Default()

	controller.AuthController(router, db)
	controller.CustomerController(router, db)
	controller.ProductController(router, db)
	controller.CartController(router, db)

	// รันเซิร์ฟเวอร์
	router.Run(":8080")
}
