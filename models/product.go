package models

type Product struct {
	ID          string  `json:"id" binding:"-"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}
