package entity

import "time"

type Product struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Stock       int       `json:"stock"`
	Created     time.Time `json:"created"`
}

func (Product) TableName() string {
	return "product"
}
