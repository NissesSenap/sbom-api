package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
