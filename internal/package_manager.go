package migrate_packages_internal

import "github.com/digiconvent/migrate_packages/db"

type PackageManager interface {
	GetPackages() ([]string, error)
	Versions() ([]string, error)
	VersionsToMigrate() ([]string, error)
	GetPackageMigration(pkg, version string) (string, error)
	// <dir> is the folder where the data of the packages is located.
	// e.g., if package 'iam' has data, the database would be located at <dir>/iam/iam.db
	// so the <dir> usually is a home directory
	MigrateDatabasesIn(dir string) (map[string]db.DatabaseInterface, error)
}
