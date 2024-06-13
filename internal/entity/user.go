package entity

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;type:serial" json:"id"`
	Name      string    `gorm:"not null;unique" json:"name"`
	Email     string    `gorm:"not null;unique" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	CreatedAt time.Time `json:"-"`
}
