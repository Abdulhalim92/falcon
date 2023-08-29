package productsuc

import (
	"falconapi/domain/entities"
)

type ProductsDataStores interface {
	GetAll() []entities.Product
	Create(product *entities.Product) error
}
