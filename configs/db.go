package configs

import (
	"fmt"
	"go-jwt/models"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBRepo struct {
	DB *gorm.DB
}

var once sync.Once
var DBInstance *DBRepo

func GetDBInstance() *DBRepo {
	once.Do(func() {
		DBInstance = &DBRepo{}
		DBInstance.ConnectDb()
	})
	return DBInstance
}

func (t *DBRepo) ConnectDb() {
	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Singapore",
		PG_HOST, PG_USER, PG_PASSWORD, PG_DB, PG_PORT,
	)

	t.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	log.Println("Connected to database successfully!")

	// Sync database models automatically
	t.syncDb()
}

func (t *DBRepo) registerModels() []interface{} {
	return []interface{}{
		&models.User{}, // Add more models here
	}
}

func (t *DBRepo) syncDb() {
	if err := t.DB.AutoMigrate(t.registerModels()...); err != nil {
		log.Fatal("Failed to migrate database models:", err)
	}
	log.Println("Database migrated successfully!")
}
