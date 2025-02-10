package sbom

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/NissesSenap/sbom-api/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CycloneDXParser struct{}

func (p *CycloneDXParser) Parse(filePath string) (interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SBOM file: %w", err)
	}

	var bom cyclonedx.BOM
	err = json.Unmarshal(data, &bom)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SBOM file: %w", err)
	}
	return bom, nil
}

func (p *CycloneDXParser) Store(ctx context.Context, dbpool *pgxpool.Pool, bom interface{}, sbomURL string) error {
	cdxBOM, ok := bom.(cyclonedx.BOM)
	if !ok {
		return fmt.Errorf("invalid BOM type")
	}

	q := db.New(dbpool)

	applicationID, err := getOrInsertApplication(ctx, q, cdxBOM.Metadata.Component.Name, sbomURL)
	if err != nil {
		return err
	}

	for _, component := range *cdxBOM.Components {
		if err := storeComponent(ctx, q, applicationID, component); err != nil {
			return err
		}
	}

	return nil
}

func getOrInsertApplication(ctx context.Context, q *db.Queries, appName, sbomURL string) (int32, error) {
	applicationID, err := q.GetApplication(ctx, appName)
	if err != nil {
		applicationID, err = q.InsertApplication(ctx, db.InsertApplicationParams{
			Name:    appName,
			SbomUrl: sbomURL,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to insert application: %w", err)
		}
	}
	return applicationID, nil
}

func storeComponent(ctx context.Context, q *db.Queries, applicationID int32, component cyclonedx.Component) error {
	packageID, err := getOrInsertPackage(ctx, q, component.Name)
	if err != nil {
		return err
	}

	if err := q.InsertVersion(ctx, db.InsertVersionParams{
		PackageID: packageID,
		Version:   component.Version,
	}); err != nil {
		return fmt.Errorf("failed to insert version: %w", err)
	}

	for _, license := range *component.Licenses {
		if err := storeLicense(ctx, q, applicationID, packageID, license.License.Name); err != nil {
			return err
		}
	}

	return nil
}

func getOrInsertPackage(ctx context.Context, q *db.Queries, packageName string) (int32, error) {
	packageID, err := q.GetPackage(ctx, packageName)
	if err != nil {
		packageID, err = q.InsertPackage(ctx, packageName)
		if err != nil {
			return 0, fmt.Errorf("failed to insert package: %w", err)
		}
	}
	return packageID, nil
}

func storeLicense(ctx context.Context, q *db.Queries, applicationID, packageID int32, licenseName string) error {
	licenseID, err := getOrInsertLicense(ctx, q, licenseName)
	if err != nil {
		return err
	}

	if _, err := q.GetApplicationPackage(ctx, db.GetApplicationPackageParams{
		ApplicationID: applicationID,
		PackageID:     packageID,
		LicenseID:     licenseID,
	}); err == nil {
		if err := q.InsertApplicationPackage(ctx, db.InsertApplicationPackageParams{
			ApplicationID: applicationID,
			PackageID:     packageID,
			LicenseID:     licenseID,
		}); err != nil {
			return fmt.Errorf("failed to insert application package: %w", err)
		}
	}

	return nil
}

func getOrInsertLicense(ctx context.Context, q *db.Queries, licenseName string) (int32, error) {
	licenseID, err := q.GetLicense(ctx, licenseName)
	if err != nil {
		licenseID, err = q.InsertLicense(ctx, licenseName)
		if err != nil {
			return 0, fmt.Errorf("failed to insert license: %w", err)
		}
	}
	return licenseID, nil
}
