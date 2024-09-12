package models

type Product struct {
	ID          int     `json:"id" binding:"-"` // Exclude ID from binding, as it's not part of the request body
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gte=0"`
}
