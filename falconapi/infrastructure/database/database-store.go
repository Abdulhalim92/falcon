package database

import (
	"falconapi/domain/entities"
	"gorm.io/gorm"
)

type productDatabase struct {
	db *gorm.DB
}

func NewProductDatabase(db *gorm.DB) *productDatabase {
	return &productDatabase{
		db: db,
	}
}

func (ds *productDatabase) Create(product *entities.Product) error {
	err := ds.db.Create(product).Error

	return err
}

func (ds *productDatabase) GetAll() []entities.Product {
	var (
		all []entities.Product
	)

	err := ds.db.Find(&all).Error
	if err != nil {
		return nil
	}

	return all
}
