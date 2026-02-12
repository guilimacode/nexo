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

	createSchema()
}
func createSchema() {
	schema := `
	CREATE TABLE IF NOT EXISTS organizations (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		document VARCHAR(20),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS addresses (
		id SERIAL PRIMARY KEY,
		street VARCHAR(100) NOT NULL,
		complement VARCHAR(50),
		number INTEGER NOT NULL,
		neighborhood VARCHAR(50) NOT NULL,
		city VARCHAR(50) NOT NULL,
		state CHAR(2) NOT NULL,
		zip_code VARCHAR(9) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS establishments (
		id SERIAL PRIMARY KEY,
		organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
		address_id INTEGER REFERENCES addresses(id) ON DELETE SET NULL,
		name VARCHAR(255) NOT NULL,
		nickname VARCHAR(50),
		document VARCHAR(20),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
		establishment_id INTEGER REFERENCES establishments(id) ON DELETE SET NULL,
		full_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL CHECK (role IN ('owner', 'manager', 'seller')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		establishment_id INTEGER NOT NULL REFERENCES establishments(id) ON DELETE CASCADE,
		name VARCHAR(150) NOT NULL,
		sku VARCHAR(50) NOT NULL,
		description TEXT,
		price BIGINT NOT NULL,
		quantity INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,		
		UNIQUE(establishment_id, sku)
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
	`

	DB.MustExec(schema)
}
