package gorm_db

import (
	"falconapi/domain/entities"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Dushanbe",
		viper.GetString("Database.Host"),
		viper.GetString("Database.Port"),
		viper.GetString("Database.Username"),
		viper.GetString("Database.Password"),
		viper.GetString("Database.Dbname"),
		viper.GetString("Database.SslMode"))

	log.Println(dsn)

	dbConn, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("Connection success-->> host:port -->> %s:%s", viper.GetString("Database.Host"), viper.GetString("Database.Port"))

	return dbConn, nil
}

func MigrateDatabase(db *gorm.DB) error {
	if !db.Migrator().HasTable(entities.Product{}) {
		return db.Migrator().AutoMigrate(entities.Product{})
	}

	return nil
}
