package migrate_packages

import (
	"errors"
	"fmt"
	"os"
	"path"
	"slices"

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
		pragmas := []string{
			// persistent pragmas
			"journal_mode = WAL",
			"foreign_keys = ON",
			// non-persistent pragmas
			"synchronous = NORMAL",
			"busy_timeout = 1000",
			"cache_size = -500",
			"mmap_size = 16777216",
		}
		dbUri := path.Join(pkgDir, pkg)

		dbConn, err := db.New(dbUri)
		if err != nil {
			return nil, err
		}
		databases[pkg] = dbConn

		for _, pragma := range pragmas {
			_, err := dbConn.Exec("PRAGMA " + pragma)
			if err != nil {
				return nil, fmt.Errorf("%s could not execute pragma command %s", dbUri, pragma)
			}
		}

		for _, version := range versions {
			pkgv := pkg + ":" + version
			script, err := d.GetPackageMigration(pkg, version)
			if err != nil && d.verbose {
				fmt.Println("[Migrate Packages] No migration for " + pkgv)
			}

			keys := make([]string, 0, len(script))
			for k := range script {
				keys = append(keys, k)
			}
			slices.Sort(keys)

			for _, s := range keys {
				_, err = dbConn.Exec(script[s])
				if err != nil {
					if d.verbose {
						fmt.Println("[Migrate Packages] ❌" + pkgv + "(" + s + ")")
					}
					return nil, errors.New("Could not migrate " + pkgv + " (" + s + "): \n" + err.Error())
				} else {
					if d.verbose {
						fmt.Println("[Migrate Packages] ✅" + pkgv + "(" + s + ")")
					}
				}
			}
		}
	}
	return databases, nil
}
