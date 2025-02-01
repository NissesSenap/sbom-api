package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := "postgres://youruser:yourpassword@localhost:5432/yourdbname"
	dbpool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	createTables(dbpool)

	fmt.Println("Database connected and tables created successfully!")
}

func createTables(dbpool *pgxpool.Pool) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS Applications (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS Packages (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS Versions (
            id SERIAL PRIMARY KEY,
            package_id INT NOT NULL,
            version VARCHAR(50) NOT NULL,
            FOREIGN KEY (package_id) REFERENCES Packages(id)
        )`,
		`CREATE TABLE IF NOT EXISTS Licenses (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS ApplicationPackages (
            application_id INT NOT NULL,
            package_id INT NOT NULL,
            license_id INT NOT NULL,
            FOREIGN KEY (application_id) REFERENCES Applications(id),
            FOREIGN KEY (package_id) REFERENCES Packages(id),
            FOREIGN KEY (license_id) REFERENCES Licenses(id),
            PRIMARY KEY (application_id, package_id)
        )`,
	}

	for _, query := range queries {
		_, err := dbpool.Exec(context.Background(), query)
		if err != nil {
			log.Fatalf("Failed to execute query: %v\n", err)
		}
	}
}
