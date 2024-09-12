package services

import (
	"database/sql"
	"simpler-products/models"
)

type ProductsServiceInterface interface {
	GetAllProducts()
	GetProductById()
	AddProduct(product *models.Product) error
	UpdateProduct()
	PatchProduct()
	DeleteProduct()
}

type ProductsService struct {
	DB *sql.DB
}

func (ps *ProductsService) GetAllProducts() {
}

func (ps *ProductsService) GetProductById() {
}

func (ps *ProductsService) AddProduct(product *models.Product) error {
	result, err := ps.DB.Exec("INSERT INTO PRODUCTS (name, description, price) VALUES (?, ?, ?)", product.Name, product.Description, product.Price)
	if err != nil {
		return err
	}

	lastInsertID, _ := result.LastInsertId()
	product.ID = int(lastInsertID)

	return nil
}

func (ps *ProductsService) UpdateProduct() {
}

func (ps *ProductsService) PatchProduct() {
}

func (ps *ProductsService) DeleteProduct() {
}
