package entity

import "time"

type Product struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Slug      string    `gorm:"not null" json:"slug"`
	CompanyID uint      `gorm:"not null" json:"company_id"`
	Company   *Company  `gorm:"foreignKey:CompanyID" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
