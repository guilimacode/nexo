package store

import (
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func Connect(connStr string) {
	var err error

	DB, err = sqlx.Connect("pgx", connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar no Postgres: %v", err)
	}

	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err := DB.Ping(); err != nil {
		log.Fatalf("Banco não está respondendo: %v", err)
	}

	log.Println("Conectado ao Banco de Dados!")
}
