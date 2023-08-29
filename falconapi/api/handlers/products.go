package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"falconapi/domain/entities"
	"falconapi/use_cases/productsuc"
)

type CreateProductUseCase interface {
	CreateProduct(ctx context.Context, request productsuc.CreateProductRequest) (*productsuc.CreateProductResponse, error)
}

type GetProductsUseCase interface {
	GetProducts(ctx context.Context) []entities.Product
}

func CreateProductHandler(useCase CreateProductUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx = c.Request.Context()

		var request = productsuc.CreateProductRequest{}

		err := c.BindJSON(&request)
		if err != nil {
			log.Println(err, "unable to parse incoming request")
			c.JSON(http.StatusBadRequest, "unable to parse incoming request")
			return
		}

		response, err := useCase.CreateProduct(ctx, request)
		if err != nil {
			log.Println(err, "unable to create product")
			c.JSON(http.StatusInternalServerError, "something went wrong")
		}

		c.JSON(http.StatusCreated, response)
	}
}

func GetProductsHandler(useCase GetProductsUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx = c.Request.Context()

		products := useCase.GetProducts(ctx)

		c.JSON(http.StatusOK, products)
	}
}
