package productsuc

import (
	"context"
	"falconapi/domain/entities"
)

type getProductsUseCase struct {
	dataStore ProductsDataStores
}

func NewGetProductsUseCase(ds ProductsDataStores) *getProductsUseCase {
	return &getProductsUseCase{
		dataStore: ds,
	}
}

func (uc *getProductsUseCase) GetProducts(ctx context.Context) []entities.Product {
	all := uc.dataStore.GetAll()
	return all
}
