package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := "postgres://sbom:sbom@localhost:5432/sbom"
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	createTables(ctx, dbpool)

	// get all data from the database
	rows, err := dbpool.Query(ctx, "SELECT * FROM Applications")
	if err != nil {
		log.Fatalf("Failed to query database: %v\n", err)
	}
	defer rows.Close()
	// iterate over the rows and print them
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatalf("Failed to scan row: %v\n", err)
		}
		fmt.Printf("Application: %d %s\n", id, name)
	}

	err = parseAndStoreSBOM(ctx, dbpool, "go-bom.json")
	if err != nil {
		log.Fatalf("Failed to parse and store SBOM: %v\n", err)
	}

	fmt.Println("Database connected, tables created, and SBOM data stored successfully!")
}

func createTables(ctx context.Context, dbpool *pgxpool.Pool) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS Applications (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE
        )`,
		`CREATE TABLE IF NOT EXISTS Packages (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE
        )`,
		`CREATE TABLE IF NOT EXISTS Versions (
            id SERIAL PRIMARY KEY,
            package_id INT NOT NULL,
            version VARCHAR(50) NOT NULL,
            FOREIGN KEY (package_id) REFERENCES Packages(id),
            UNIQUE (package_id, version)
        )`,
		`CREATE TABLE IF NOT EXISTS Licenses (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE
        )`,
		`CREATE TABLE IF NOT EXISTS ApplicationPackages (
            application_id INT NOT NULL,
            package_id INT NOT NULL,
            license_id INT NOT NULL,
            FOREIGN KEY (application_id) REFERENCES Applications(id),
            FOREIGN KEY (package_id) REFERENCES Packages(id),
            FOREIGN KEY (license_id) REFERENCES Licenses(id),
            PRIMARY KEY (application_id, package_id, license_id)
        )`,
	}

	for _, query := range queries {
		_, err := dbpool.Exec(ctx, query)
		if err != nil {
			log.Fatalf("Failed to execute query: %v\n", err)
		}
	}
}

func parseAndStoreSBOM(ctx context.Context, dbpool *pgxpool.Pool, filePath string) error {
	// Read the SBOM file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SBOM file: %w", err)
	}

	// Parse the SBOM file
	var bom cyclonedx.BOM
	err = json.Unmarshal(data, &bom)
	if err != nil {
		return fmt.Errorf("failed to parse SBOM file: %w", err)
	}

	// Insert application
	applicationID, err := getOrInsertApplication(ctx, dbpool, bom.Metadata.Component.Name)
	if err != nil {
		return fmt.Errorf("failed to insert application: %w", err)
	}

	// Store the SBOM data in the database
	for _, component := range *bom.Components {
		// Insert package
		packageID, err := getOrInsertPackage(ctx, dbpool, component.Name)
		if err != nil {
			return fmt.Errorf("failed to insert package: %w", err)
		}

		// Insert version
		err = getOrInsertVersion(ctx, dbpool, packageID, component.Version)
		if err != nil {
			return fmt.Errorf("failed to insert version: %w", err)
		}

		// Insert licenses
		for _, license := range *component.Licenses {
			licenseID, err := getOrInsertLicense(ctx, dbpool, license.License.ID)
			if err != nil {
				return fmt.Errorf("failed to insert license: %w", err)
			}

			// Insert application package
			err = getOrInsertApplicationPackage(ctx, dbpool, applicationID, packageID, licenseID)
			if err != nil {
				return fmt.Errorf("failed to insert application package: %w", err)
			}
		}
	}

	return nil
}

func getApplication(ctx context.Context, dbpool *pgxpool.Pool, name string) (int, error) {
	var id int
	err := dbpool.QueryRow(ctx, `SELECT id FROM Applications WHERE name = $1`, name).Scan(&id)
	return id, err
}

func insertApplication(ctx context.Context, dbpool *pgxpool.Pool, name string) (int, error) {
	var id int
	err := dbpool.QueryRow(ctx, `INSERT INTO Applications (name) VALUES ($1) RETURNING id`, name).Scan(&id)
	return id, err
}

func getOrInsertApplication(ctx context.Context, dbpool *pgxpool.Pool, name string) (int, error) {
	id, err := getApplication(ctx, dbpool, name)
	if err == nil {
		return id, nil
	}
	return insertApplication(ctx, dbpool, name)
}

func getPackage(ctx context.Context, dbpool *pgxpool.Pool, name string) (int, error) {
	var id int
	err := dbpool.QueryRow(ctx, `SELECT id FROM Packages WHERE name = $1`, name).Scan(&id)
	return id, err
}

func insertPackage(ctx context.Context, dbpool *pgxpool.Pool, name string) (int, error) {
	var id int
	err := dbpool.QueryRow(ctx, `INSERT INTO Packages (name) VALUES ($1) RETURNING id`, name).Scan(&id)
	return id, err
}

func getOrInsertPackage(ctx context.Context, dbpool *pgxpool.Pool, name string) (int, error) {
	id, err := getPackage(ctx, dbpool, name)
	if err == nil {
		return id, nil
	}
	return insertPackage(ctx, dbpool, name)
}

func getVersion(ctx context.Context, dbpool *pgxpool.Pool, packageID int, version string) (int, error) {
	var id int
	err := dbpool.QueryRow(ctx, `SELECT id FROM Versions WHERE package_id = $1 AND version = $2`, packageID, version).Scan(&id)
	return id, err
}

func insertVersion(ctx context.Context, dbpool *pgxpool.Pool, packageID int, version string) error {
	_, err := dbpool.Exec(ctx, `INSERT INTO Versions (package_id, version) VALUES ($1, $2)`, packageID, version)
	return err
}

func getOrInsertVersion(ctx context.Context, dbpool *pgxpool.Pool, packageID int, version string) error {
	_, err := getVersion(ctx, dbpool, packageID, version)
	if err == nil {
		return nil
	}
	return insertVersion(ctx, dbpool, packageID, version)
}

func getLicense(ctx context.Context, dbpool *pgxpool.Pool, name string) (int, error) {
	var id int
	err := dbpool.QueryRow(ctx, `SELECT id FROM Licenses WHERE name = $1`, name).Scan(&id)
	return id, err
}

func insertLicense(ctx context.Context, dbpool *pgxpool.Pool, name string) (int, error) {
	var id int
	err := dbpool.QueryRow(ctx, `INSERT INTO Licenses (name) VALUES ($1) RETURNING id`, name).Scan(&id)
	return id, err
}

func getOrInsertLicense(ctx context.Context, dbpool *pgxpool.Pool, name string) (int, error) {
	id, err := getLicense(ctx, dbpool, name)
	if err == nil {
		return id, nil
	}
	return insertLicense(ctx, dbpool, name)
}

func getApplicationPackage(ctx context.Context, dbpool *pgxpool.Pool, applicationID, packageID, licenseID int) (int, error) {
	var id int
	err := dbpool.QueryRow(ctx, `SELECT 1 FROM ApplicationPackages WHERE application_id = $1 AND package_id = $2 AND license_id = $3`, applicationID, packageID, licenseID).Scan(&id)
	return id, err
}

func insertApplicationPackage(ctx context.Context, dbpool *pgxpool.Pool, applicationID, packageID, licenseID int) error {
	_, err := dbpool.Exec(ctx, `INSERT INTO ApplicationPackages (application_id, package_id, license_id) VALUES ($1, $2, $3)`, applicationID, packageID, licenseID)
	return err
}

func getOrInsertApplicationPackage(ctx context.Context, dbpool *pgxpool.Pool, applicationID, packageID, licenseID int) error {
	_, err := getApplicationPackage(ctx, dbpool, applicationID, packageID, licenseID)
	if err == nil {
		fmt.Println("This is already in the database")
		return nil
	}
	return insertApplicationPackage(ctx, dbpool, applicationID, packageID, licenseID)
}
