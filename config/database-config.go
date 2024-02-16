package config

import (
	"fmt"
	mysqlSession "github.com/go-session/mysql"
	"github.com/go-session/session"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/webjohny/cashflow-go/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SetupDatabaseConnection is creating a new connection to our database
func SetupDatabaseConnection() *gorm.DB {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	session.InitManager(
		session.SetStore(mysqlSession.NewStore(mysqlSession.NewConfig(dsn), "sessions", 0)),
	)

	if err != nil {
		panic("Failed to create connection to database")
	}

	//Isi model / table disini
	db.AutoMigrate(&entity.User{}, &entity.Race{})
	return db
}

// CloseDatabaseConnection Close database connection
func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to close connection from database")
	}
	dbSQL.Close()
}
