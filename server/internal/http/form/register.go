package form

import (
	"time"
)

type Register struct {
	Email     string    `json:"email" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	Birth     time.Time `json:"birth_date" binding:"required"`
	Sex       string    `json:"sex" binding:"required"`
	Interest  string    `json:"interest"`
	City      string    `json:"city"`
}
