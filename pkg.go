package migrate_packages

import (
	"os"
	"path"
	"slices"

	migrate_packages_internal "github.com/digiconvent/migrate_packages/internal"
)

type data struct {
	fromVersion *Version
	toVersion   *Version
	pkgDir      string
}

func (d *data) GetPackages() ([]string, error) {
	entries, err := os.ReadDir(d.pkgDir)
	if err != nil {
		return nil, err
	}
	var packages = []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			packages = append(packages, entry.Name())
		}
	}
	return packages, nil
}

func (d *data) Versions() ([]string, error) {
	var versions = []string{}
	packages, err := d.GetPackages()
	if err != nil {
		return nil, err
	}

	for _, pkg := range packages {
		entries, err := os.ReadDir(path.Join(d.pkgDir, pkg, "db"))
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if ToVersion(entry.Name()) == nil {
				continue
			}

			if !slices.Contains(versions, entry.Name()) {
				versions = append(versions, entry.Name())
			}
		}
	}
	return versions, nil
}

func (d *data) To(version *Version) packageManagerChoice {
	d.toVersion = version
	return d
}

func (d *data) ToVersion(ma int, mi int, pa int) packageManagerChoice {
	d.toVersion = &Version{Major: ma, Minor: mi, Patch: pa}
	return d
}

func (d *data) WithLocalFilesAt(projectRoot, pkgDir string) (packageManager, error) {
	d.pkgDir = path.Join(projectRoot, pkgDir)
	return d, nil
}

func (d *data) WithPrivateRepository(username string, repository string, token string) (repoPackageManager, error) {
	err := migrate_packages_internal.DownloadRepoZip(username, repository, token)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (d *data) WithPublicRepository(username string, repository string) (repoPackageManager, error) {
	err := migrate_packages_internal.DownloadRepoZip(username, repository, "")
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type packageManager interface {
	GetPackages() ([]string, error)
	Versions() ([]string, error)
}

func From(version *Version) migrateToVersion {
	return &data{
		fromVersion: version,
	}
}

// if fresh, pass on a -1.-1.-1
func FromSemVer(ma, mi, pa int) migrateToVersion {
	return From(&Version{Major: ma, Minor: mi, Patch: pa})
}

type migrateToVersion interface {
	// to migrate to latest, pass on a -1,-1,-1
	ToVersion(ma, mi, pa int) packageManagerChoice
	// to migrate to latest, pass on a nil
	To(version *Version) packageManagerChoice
}

type packageManagerChoice interface {
	WithPublicRepository(username, repository string) (repoPackageManager, error)
	WithPrivateRepository(username, repository, token string) (repoPackageManager, error)
	WithLocalFilesAt(projectRoot, pkgDir string) (packageManager, error)
}

type repoPackageManager interface {
	WithPkgDir(dir string) packageManager
}
