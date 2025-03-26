package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Customer struct {
	CustomerID  int       `gorm:"primaryKey;autoIncrement;column:customer_id" json:"customer_id"`
	FirstName   string    `gorm:"column:first_name;not null" json:"first_name"`
	LastName    string    `gorm:"column:last_name;not null" json:"last_name"`
	Email       string    `gorm:"column:email;unique;not null" json:"email"`
	PhoneNumber string    `gorm:"column:phone_number" json:"phone_number"`
	Address     string    `gorm:"column:address" json:"address"`
	Password    string    `gorm:"column:password;not null" json:"-"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// ความสัมพันธ์
	Carts []Cart `gorm:"foreignKey:CustomerID" json:"-"`
}

func (m *Customer) TableName() string {
	return "customer"
}

// BeforeSave - Hook สำหรับ hash password ก่อนบันทึก (ทั้งสร้างใหม่และอัปเดต)
func (c *Customer) BeforeSave(tx *gorm.DB) error {
	if len(c.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		c.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword ตรวจสอบรหัสผ่าน
func (c *Customer) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
	return err == nil
}

// CustomerResponse โครงสร้างสำหรับส่งข้อมูลลูกค้าโดยไม่รวมรหัสผ่าน
type CustomerResponse struct {
	CustomerID  int       `json:"customer_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Address     string    `json:"address,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
