// Code generated by sql2gorm. DO NOT EDIT.
package model

import (
)
// cart_item.go
type CartItem struct {
    CartItemID uint    `gorm:"primaryKey" json:"cart_item_id"`
    CartID     uint    `gorm:"not null" json:"cart_id"`
    ProductID  uint    `gorm:"not null" json:"product_id"`
    Quantity   int     `gorm:"not null" json:"quantity"`
    
    Cart    Cart    `gorm:"foreignKey:CartID" json:"-"`
    Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (m *CartItem) TableName() string {
	return "cart_item"
}

