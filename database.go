package migrate_packages

import (
	"fmt"
	"os"
	"path"

	"github.com/digiconvent/migrate_packages/db"
)

func (d *data) MigrateDatabasesIn(dir string) (map[string]db.DatabaseInterface, error) {
	if _, err := os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return nil, fmt.Errorf("could not create required data folder at %s: %s", dir, err)
		}
	}

	d.dataDir = dir
	packages, err := d.GetPackages()
	if err != nil {
		return nil, err
	}

	versions, err := d.VersionsToMigrate()
	if err != nil {
		return nil, err
	}

	var databases map[string]db.DatabaseInterface = make(map[string]db.DatabaseInterface)
	for _, pkg := range packages {
		migrationScript := ""
		for _, version := range versions {
			script, _ := d.GetPackageMigration(pkg, version)
			migrationScript += script
		}

		pkgDir := path.Join(dir, pkg)
		os.MkdirAll(pkgDir, 0700)
		dbUri := path.Join(pkgDir, pkg+".db")
		dbConn, err := db.New(dbUri)
		if err != nil {
			return nil, err
		}

		_, err = dbConn.Exec(migrationScript)
		if err != nil {
			return nil, err
		}
		databases[pkg] = dbConn
	}
	return databases, nil
}
