package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"page-analyzer/internal/models/mysql"
)

type application struct {
	domains      *mysql.DomainModel
	domainChecks *mysql.DomainCheckModel
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	db, err := openDb(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	app := &application{
		domains: &mysql.DomainModel{
			DB: db,
		},
		domainChecks: &mysql.DomainCheckModel{
			DB: db,
		},
	}

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err = http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

func openDb(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
