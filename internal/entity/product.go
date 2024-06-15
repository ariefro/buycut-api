package entity

import "time"

type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;unique" json:"name"`
	Slug      string    `gorm:"not null;unique" json:"slug"`
	CompanyID uint      `gorm:"not null" json:"-"`
	Company   *Company  `gorm:"foreignKey:CompanyID" json:"company"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
