package config

import (
	"database/sql"
	"fmt"
	mysqlSession "github.com/go-session/mysql"
	"github.com/go-session/session"
	"log"
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

	InitDatabaseConnection()

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
	db.AutoMigrate(&entity.User{}, &entity.Request{}, &entity.Race{}, &entity.Lobby{}, &entity.Player{}, &entity.Transaction{})
	return db
}

// InitDatabaseConnection is initial script
func InitDatabaseConnection() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbRootPass := os.Getenv("MYSQL_ROOT_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8&parseTime=True&loc=Local", "root", dbRootPass, dbHost, dbPort)

	// Open a connection to the database.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// Ensure the database connection is available.
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	var dbExists sql.NullString
	db.QueryRow("SELECT table_schema AS db FROM information_schema.TABLES WHERE table_schema = '" + dbName + "' GROUP BY table_schema").Scan(&dbExists)

	if dbExists.Valid {
		return
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.Fatal("Error creating database: ", err)
	}

	_, err = db.Exec("CREATE USER IF NOT EXISTS '" + dbUser + "'@'%' IDENTIFIED BY '" + dbPass + "'")
	if err != nil {
		log.Fatal("Error user "+dbUser+" for database: ", err)
	}

	_, err = db.Exec("GRANT ALL PRIVILEGES ON " + dbName + ".* TO '" + dbUser + "'@'%'")
	if err != nil {
		log.Fatal("Error grant all privileges "+dbUser+" database: ", err)
	}

	_, err = db.Exec("FLUSH PRIVILEGES")

	if err != nil {
		log.Fatal("Error flushing privileges: ", err)
	}

	fmt.Println("Database created successfully!")
}

// CloseDatabaseConnection Close database connection
func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to close connection from database")
	}
	dbSQL.Close()
}
