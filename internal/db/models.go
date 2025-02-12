// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

type Application struct {
	ID   int32
	Name string
}

type Applicationpackage struct {
	ApplicationID int32
	PackageID     int32
	LicenseID     int32
}

type Applicationversion struct {
	ID            int32
	ApplicationID int32
	Version       string
	SbomUrl       string
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
