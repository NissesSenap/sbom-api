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

func (p *CycloneDXParser) Store(ctx context.Context, dbpool *pgxpool.Pool, bom interface{}) error {
	cdxBOM, ok := bom.(cyclonedx.BOM)
	if !ok {
		return fmt.Errorf("invalid BOM type")
	}

	q := db.New(dbpool)

	// Insert application
	applicationID, err := q.GetApplication(ctx, cdxBOM.Metadata.Component.Name)
	if err != nil {
		applicationID, err = q.InsertApplication(ctx, cdxBOM.Metadata.Component.Name)
		if err != nil {
			return fmt.Errorf("failed to insert application: %w", err)
		}
	}

	// Store the SBOM data in the database
	for _, component := range *cdxBOM.Components {
		// Insert package
		packageID, err := q.GetPackage(ctx, component.Name)
		if err != nil {
			packageID, err = q.InsertPackage(ctx, component.Name)
			if err != nil {
				return fmt.Errorf("failed to insert package: %w", err)
			}
		}

		// Insert version
		err = q.InsertVersion(ctx, db.InsertVersionParams{
			PackageID: packageID,
			Version:   component.Version,
		})
		if err != nil {
			return fmt.Errorf("failed to insert version: %w", err)
		}

		// Insert license
		for _, license := range *component.Licenses {
			licenseID, err := q.GetLicense(ctx, license.License.Name)
			if err != nil {
				licenseID, err = q.InsertLicense(ctx, license.License.Name)
				if err != nil {
					return fmt.Errorf("failed to insert license: %w", err)
				}
			}

			_, err = q.GetApplicationPackage(ctx, db.GetApplicationPackageParams{
				ApplicationID: applicationID,
				PackageID:     packageID,
				LicenseID:     licenseID,
			})
			if err == nil {
				// Insert application package
				err = q.InsertApplicationPackage(ctx, db.InsertApplicationPackageParams{
					ApplicationID: applicationID,
					PackageID:     packageID,
					LicenseID:     licenseID,
				})
				if err != nil {
					return fmt.Errorf("failed to insert application package: %w", err)
				}
			}
		}
	}

	return nil
}
