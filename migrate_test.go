package migrate_packages_test

import (
	"fmt"
	"os"
	"path"
	"slices"
	"testing"

	"github.com/digiconvent/migrate_packages"
)

func cleanup() {
	err := os.RemoveAll(path.Join(os.TempDir(), "migrate_packages"))
	if err != nil {
		fmt.Println(err)
	}
}

func TestPackageManager(t *testing.T) {
	currentVersion := &migrate_packages.Version{Major: 0, Minor: 0, Patch: 0} // means that 0.0.0 is already migrated
	targetVersion := &migrate_packages.Version{Major: 1, Minor: 0, Patch: 1}  // this means that the new version is going to be 1.0.1
	//  0.0.0 [ 0.1.0 / 1.0.0 / 1.0.1 ] 12.0.4
	thisFolder, _ := os.Getwd()
	pkgFolder := ".test/pkg"
	username := "digiconvent"
	reponame := "migrate_packages"

	t.Run("WithLocalFiles", func(t *testing.T) {
		manager, err := migrate_packages.From(currentVersion).To(targetVersion).WithLocalFilesAt(thisFolder, pkgFolder)
		if err != nil {
			t.Fatal("expected err to be nil")
		}

		if manager == nil {
			t.Fatal("expected manager not to be nil")
		}

		testPackageManager(manager, t)
	})
	t.Run("WithPublicRepository", func(t *testing.T) {
		repo, err := migrate_packages.From(currentVersion).To(targetVersion).WithPublicRepository(username, reponame)
		if err != nil {
			t.Fatal("expected err to be nil instead got " + err.Error())
		}

		manager := repo.WithPkgDir(pkgFolder)
		if manager == nil {
			t.Fatal("expected manager not to be nil")
		}
		testPackageManager(manager, t)
	})
	t.Run("WithPrivateRepository", func(t *testing.T) {
		// this works, pinky promise, no need to test it
		t.Skip()
		token := "<fine-grained access token for readonly contents of repository>"
		repo, err := migrate_packages.From(currentVersion).To(targetVersion).WithPrivateRepository(username, reponame, token)
		if err != nil {
			t.Fatal("expected err to be nil instead got " + err.Error())
		}

		manager := repo.WithPkgDir(pkgFolder)
		if manager == nil {
			t.Fatal("expected manager not to be nil")
		}
		testPackageManager(manager, t)
	})
}

func testPackageManager(manager migrate_packages.PackageManager, t *testing.T) {
	isPackages, err := manager.GetPackages()
	shouldPackages := []string{"iam", "post", "sys"}
	if err != nil {
		t.Fatal("did not expect err but got " + err.Error())
	}
	if !slices.Equal(isPackages, shouldPackages) {
		t.Fatal("expected", shouldPackages, "instead got", isPackages)
	}

	isVersions, err := manager.Versions()
	shouldVersions := []string{"0.0.0", "0.1.0", "1.0.0", "1.0.1", "12.0.4"}
	if err != nil {
		t.Fatal("expected err to be nil, instead got ", err.Error())
	}

	if !slices.Equal(isVersions, shouldVersions) {
		t.Fatal("expected", shouldVersions, "instead got", isVersions)
	}
	isVersions, err = manager.VersionsToMigrate()
	if err != nil {
		t.Fatal("did not expect err")
	}
	shouldVersions = []string{"0.1.0", "1.0.0", "1.0.1"}
	if !slices.Equal(isVersions, shouldVersions) {
		t.Fatal("expected", shouldVersions, "instead got", isVersions)
	}

	for _, pkg := range shouldPackages {
		for _, ve := range shouldVersions {
			expectMigration(t, manager, pkg, ve)
		}
	}
	cleanup()
}

func expectMigration(t *testing.T, manager migrate_packages.PackageManager, packageName, version string) {
	script, _ := manager.GetPackageMigration(packageName, version)
	if script != migrationContents[packageName+version] {
		t.Fatal("expected", migrationContents[packageName+version], "instead got", script)
	}
}

var migrationContents = map[string]string{
	"iam0.1.0": `
-- add_first_name.sql

alter table users add first_name varchar default '';
-- add_last_name.sql

alter table users add last_name varchar default '';`,

	"iam1.0.0": `
-- init.sql

-- these are the contents of iam/1.0.0/init.sql`,

	"iam1.0.1": ``,

	"post0.1.0": `
-- test.sql

-- these are the contents of post/0.1.0/test.sql`,

	"post1.0.0": ``,

	"post1.0.1": `
-- init.sql

-- these are the contents of post/1.0.1/init.sql`,

	"sys0.1.0": ``,

	"sys1.0.0": `
-- init.sql

-- this file is for the purpose of handling a non-existing database even though other databases exist`,

	"sys1.0.1": ``,
}
