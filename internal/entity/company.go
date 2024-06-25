package entity

import (
	"time"

	"github.com/lib/pq"
)

type Company struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;unique" json:"name"`
	Slug        string         `gorm:"not null;unique" json:"slug"`
	Description string         `gorm:"not null" json:"description"`
	ImageURL    string         `gorm:"not null;type:varchar(255)" json:"image_url"`
	Proof       pq.StringArray `gorm:"not null;type:text[]" json:"proof"`
	Brands      []Brand        `gorm:"foreignKey:CompanyID" json:"brands,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"-"`
}
