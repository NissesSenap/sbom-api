// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Application struct {
	ID      int32
	Name    string
	SbomUrl pgtype.Text
}

type Applicationpackage struct {
	ApplicationID int32
	PackageID     int32
	LicenseID     int32
}

type License struct {
	ID   int32
	Name string
}

type Package struct {
	ID   int32
	Name string
}

type Version struct {
	ID        int32
	PackageID int32
	Version   string
}
