-- name: GetApplication :one
SELECT id FROM Applications WHERE name = $1;

-- name: InsertApplication :one
INSERT INTO Applications (name) VALUES ($1) RETURNING id;

-- name: GetPackage :one
SELECT id FROM Packages WHERE name = $1;

-- name: InsertPackage :one
INSERT INTO Packages (name) VALUES ($1) RETURNING id;

-- name: GetVersion :one
SELECT id FROM Versions WHERE package_id = $1 AND version = $2;

-- name: InsertVersion :exec
INSERT INTO Versions (package_id, version) VALUES ($1, $2);

-- name: GetLicense :one
SELECT id FROM Licenses WHERE name = $1;

-- name: InsertLicense :one
INSERT INTO Licenses (name) VALUES ($1) RETURNING id;

-- name: GetApplicationPackage :one
SELECT 1 FROM ApplicationPackages WHERE application_id = $1 AND package_id = $2 AND license_id = $3;

-- name: InsertApplicationPackage :exec
INSERT INTO ApplicationPackages (application_id, package_id, license_id) VALUES ($1, $2, $3);