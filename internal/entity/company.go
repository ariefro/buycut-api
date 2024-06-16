package entity

import "time"

type Company struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null;unique" json:"name"`
	Slug        string    `gorm:"not null;unique" json:"slug"`
	Description string    `gorm:"not null" json:"description"`
	ImageURL    string    `gorm:"type:varchar(255)" json:"image_url"`
	Products    []Product `gorm:"foreignKey:CompanyID" json:"products,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"-"`
}
