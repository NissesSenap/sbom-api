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

func ParseSBOM(filePath string) (cyclonedx.BOM, error) {
	// Read the SBOM file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return cyclonedx.BOM{}, fmt.Errorf("failed to read SBOM file: %w", err)
	}

	// Parse the SBOM file
	var bom cyclonedx.BOM
	err = json.Unmarshal(data, &bom)
	if err != nil {
		return cyclonedx.BOM{}, fmt.Errorf("failed to parse SBOM file: %w", err)
	}
	return bom, nil
}

func StoreSBOM(ctx context.Context, dbpool *pgxpool.Pool, bom cyclonedx.BOM) error {
	// Insert application
	applicationID, err := db.GetOrInsertApplication(ctx, dbpool, bom.Metadata.Component.Name)
	if err != nil {
		return fmt.Errorf("failed to insert application: %w", err)
	}

	// Store the SBOM data in the database
	for _, component := range *bom.Components {
		// Insert package
		packageID, err := db.GetOrInsertPackage(ctx, dbpool, component.Name)
		if err != nil {
			return fmt.Errorf("failed to insert package: %w", err)
		}

		// Insert version
		err = db.GetOrInsertVersion(ctx, dbpool, packageID, component.Version)
		if err != nil {
			return fmt.Errorf("failed to insert version: %w", err)
		}

		// Insert licenses
		for _, license := range *component.Licenses {
			licenseID, err := db.GetOrInsertLicense(ctx, dbpool, license.License.ID)
			if err != nil {
				return fmt.Errorf("failed to insert license: %w", err)
			}

			// Insert application package
			err = db.GetOrInsertApplicationPackage(ctx, dbpool, applicationID, packageID, licenseID)
			if err != nil {
				return fmt.Errorf("failed to insert application package: %w", err)
			}
		}
	}

	return nil
}
