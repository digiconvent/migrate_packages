package migrate_packages

import (
	"errors"
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
		pkgDir := path.Join(dir, pkg)
		err = os.MkdirAll(pkgDir, 0700)
		if err != nil {
			return nil, err
		}
		dbUri := path.Join(pkgDir, pkg+".db")
		dbConn, err := db.New(dbUri)
		if err != nil {
			return nil, err
		}
		databases[pkg] = dbConn

		for _, version := range versions {
			script, err := d.GetPackageMigration(pkg, version)
			if err != nil {
				return nil, errors.New("Could not fetch " + pkg + ":" + version)
			}

			keys := make([]string, 0, len(script))
			for k := range script {
				keys = append(keys, k)
			}

			for _, s := range keys {
				_, err = dbConn.Exec(script[s])
				if err != nil {
					return nil, errors.New("Could not migrate " + pkg + ":" + version + " (" + s + "): \n" + err.Error())
				}
			}
		}
	}
	return databases, nil
}
