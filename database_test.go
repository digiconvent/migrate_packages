package migrate_packages_test

import (
	"slices"
	"testing"

	"github.com/digiconvent/migrate_packages"
)

var currentVersion = migrate_packages.ToVersion("-1.0.0")
var targetVersion = migrate_packages.ToVersion("1.0.1")

// [ 0.1.0 / 1.0.0 / 1.0.1 ]
var pkgFolder = ".test/pkg"
var username = "digiconvent"
var reponame = "migrate_packages"
var DataFolder = "/home/digiconvent/data"

func TestDatabaseHandling(t *testing.T) {
	repo, err := migrate_packages.From(currentVersion).To(targetVersion).WithPublicRepository(username, reponame)
	if err != nil {
		t.Fatal(err)
	}
	databases, err := repo.WithPkgDir(pkgFolder).MigrateDatabasesIn(DataFolder)
	if err != nil {
		t.Fatal(err)
	}

	isKeys := make([]string, 0, len(databases))
	for k := range databases {
		isKeys = append(isKeys, k)
	}

	shouldKeys := []string{"iam", "post", "sys"}
	if !slices.Equal(shouldKeys, isKeys) {
		t.Fatal("expected", shouldKeys, "instead got", isKeys)
	}
	// this package promotes separate sqlite databases for every package under /pkg
	// so every pkg needs to handle their own database
	// create live/test
	// - this needs a folder where the database would be (/home/user/data/<pkgName>/database)
	// migrate
	// - database exists, migration script should work on it
	//	delete
	// - test databases should be deleted after every run so there will be a fresh start for the next time
	// - production databases should NOT be deleted once they are created

}
