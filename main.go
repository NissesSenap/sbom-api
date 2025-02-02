package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := "postgres://sbom:sbom@localhost:5432/sbom"
	dbpool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	createTables(dbpool)

	err = parseAndStoreSBOM(dbpool, "go-bom.json")
	if err != nil {
		log.Fatalf("Failed to parse and store SBOM: %v\n", err)
	}

	fmt.Println("Database connected, tables created, and SBOM data stored successfully!")
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

func parseAndStoreSBOM(dbpool *pgxpool.Pool, filePath string) error {
	// Read the SBOM file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SBOM file: %w", err)
	}

	// Parse the SBOM file
	var bom cyclonedx.BOM
	err = json.Unmarshal(data, &bom)
	if err != nil {
		return fmt.Errorf("failed to parse SBOM file: %w", err)
	}

	// Store the SBOM data in the database
	for _, component := range *bom.Components {
		// Insert package
		var packageID int
		err = dbpool.QueryRow(context.Background(), `INSERT INTO Packages (name) VALUES ($1) RETURNING id`, component.Name).Scan(&packageID)
		if err != nil {
			return fmt.Errorf("failed to insert package: %w", err)
		}

		// Insert version
		_, err = dbpool.Exec(context.Background(), `INSERT INTO Versions (package_id, version) VALUES ($1, $2)`, packageID, component.Version)
		if err != nil {
			return fmt.Errorf("failed to insert version: %w", err)
		}

		// Insert licenses
		for _, license := range component.Licenses {
			var licenseID int
			err = dbpool.QueryRow(context.Background(), `INSERT INTO Licenses (name) VALUES ($1) RETURNING id`, license.License.ID).Scan(&licenseID)
			if err != nil {
				return fmt.Errorf("failed to insert license: %w", err)
			}

			// Insert application package
			_, err = dbpool.Exec(context.Background(), `INSERT INTO ApplicationPackages (application_id, package_id, license_id) VALUES ($1, $2, $3)`, 1, packageID, licenseID)
			if err != nil {
				return fmt.Errorf("failed to insert application package: %w", err)
			}
		}
	}

	return nil
}
