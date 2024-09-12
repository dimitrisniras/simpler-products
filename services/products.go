package services

type ProductsServiceInterface interface {
	GetAllProducts()
	GetProductById()
	AddProduct()
	UpdateProduct()
	PatchProduct()
	DeleteProduct()
}

type ProductsService struct{}

func (ps *ProductsService) GetAllProducts() {
}

func (ps *ProductsService) GetProductById() {
}

func (ps *ProductsService) AddProduct() {
}

func (ps *ProductsService) UpdateProduct() {
}

func (ps *ProductsService) PatchProduct() {
}

func (ps *ProductsService) DeleteProduct() {
}
