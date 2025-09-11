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
	t.Run("WithLocalFiles", func(t *testing.T) {
		thisFolder, _ := os.Getwd()
		pkgFolder := ".test/pkg"
		currentVersion := &migrate_packages.Version{Major: -1, Minor: -1, Patch: -1}
		targetVersion := &migrate_packages.Version{Major: -1, Minor: -1, Patch: -1}

		manager, err := migrate_packages.From(currentVersion).To(targetVersion).WithLocalFilesAt(thisFolder, pkgFolder)
		if err != nil {
			t.Fatal("expected err to be nil")
		}

		if manager == nil {
			t.Fatal("expected manager not to be nil")
		}

		packages, err := manager.GetPackages()
		if err != nil {
			t.Fatal("did not expect err but got " + err.Error())
		}
		if !slices.Contains(packages, "iam") || !slices.Contains(packages, "post") {
			t.Fatal("expected both iam and post to be part of the packages, instead got", packages)
		}
		t.Log(manager.Versions())
		cleanup()
	})
	t.Run("WithPublicRepository", func(t *testing.T) {
		username := "digiconvent"
		reponame := "migrate_packages"
		pkgFolder := ".test/pkg"
		currentVersion := &migrate_packages.Version{Major: -1, Minor: -1, Patch: -1}
		targetVersion := &migrate_packages.Version{Major: -1, Minor: -1, Patch: -1}

		repo, err := migrate_packages.From(currentVersion).To(targetVersion).WithPublicRepository(username, reponame)
		if err != nil {
			t.Fatal("expected err to be nil instead got " + err.Error())
		}

		manager := repo.WithPkgDir(pkgFolder)
		if manager == nil {
			t.Fatal("expected manager not to be nil")
		}

		packages, err := manager.GetPackages()
		if err != nil {
			t.Fatal("did not expect err but got " + err.Error())
		}
		if !slices.Contains(packages, "iam") || !slices.Contains(packages, "post") {
			t.Fatal("expected both iam and post to be part of the packages, instead got", packages)
		}

		t.Log(manager.Versions())
		cleanup()
	})
}
