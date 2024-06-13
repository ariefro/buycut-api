package entity

import "time"

type Company struct {
	ID        uint      `gorm:"primaryKey;type:serial" json:"id"`
	Name      string    `gorm:"not null;unique" json:"name"`
	Slug      string    `gorm:"not null;unique" json:"slug"`
	Products  []Product `gorm:"foreignKey:CompanyID" json:"products,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
