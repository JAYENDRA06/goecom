package main

import (
	"database/sql"
	"log"

	"github.com/JAYENDRA06/apiproject/cmd/api"
	"github.com/JAYENDRA06/apiproject/config"
	"github.com/JAYENDRA06/apiproject/db"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		log.Println("Continuing with system environment variables...")
	}
	envs := config.InitConfig()
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 envs.DBUser,
		Passwd:               envs.DBPassword,
		Addr:                 envs.DBAddress,
		DBName:               envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	intiStorage(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func intiStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("DB: successfully connected")
}
