package migrate_packages

import (
	"os"
	"path"
	"slices"
	"strings"

	"github.com/digiconvent/migrate_packages/db"
	migrate_packages_internal "github.com/digiconvent/migrate_packages/internal"
)

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

type data struct {
	fromVersion *Version
	toVersion   *Version
	pkgDir      string
	dataDir     string
}

func (d *data) GetPackageMigration(pkg string, version string) (string, error) {
	versionInPkgDir := path.Join(d.pkgDir, pkg, "db", version)
	files, err := os.ReadDir(versionInPkgDir)
	if err != nil {
		return "", err
	}

	var migrationScript string = ""
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		// these files cannot be too big to handle so reading the entire contents of an sql file should not be a problem
		contents, err := os.ReadFile(path.Join(versionInPkgDir, file.Name()))
		if err != nil {
			continue
		}
		migrationScript += "\n-- " + file.Name() + "\n\n" + string(contents)
	}
	return migrationScript, nil
}

func (d *data) WithPkgDir(dir string) PackageManager {
	downloadFolder := path.Join(os.TempDir(), "migrate_packages")
	d.pkgDir = path.Join(downloadFolder, dir)

	segments := strings.Split(dir, "/")
	toKeep := ""
	for i := range len(segments) {
		toScan := downloadFolder + toKeep
		entries, _ := os.ReadDir(toScan)
		toKeep += "/" + segments[i]
		for _, entry := range entries {
			uri := path.Join(toScan, entry.Name())
			keep := strings.HasSuffix(uri, toKeep)
			if !keep {
				os.RemoveAll(uri)
			}
		}
	}

	packages, err := d.GetPackages()
	if err != nil {
		return nil
	}

	for _, pkg := range packages {
		pkgDir := path.Join(downloadFolder, "", dir, pkg)
		entries, _ := os.ReadDir(pkgDir)
		for _, entry := range entries {
			toRemove := path.Join(pkgDir, entry.Name())
			if entry.Name() != "db" {
				os.RemoveAll(toRemove)
			}
		}
	}

	return d
}

func (d *data) GetPackages() ([]string, error) {
	entries, err := os.ReadDir(d.pkgDir)
	if err != nil {
		return nil, err
	}
	var packages = []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			// only append if the folder is not empty
			dirEntries, err := os.ReadDir(path.Join(d.pkgDir, entry.Name()))
			if err != nil {
				return nil, err
			}
			if len(dirEntries) > 0 {
				packages = append(packages, entry.Name())
			}
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
			continue // if db doesn't exist, don't do anything
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

func (d *data) WithLocalFilesAt(projectRoot, pkgDir string) (PackageManager, error) {
	d.pkgDir = path.Join(projectRoot, pkgDir)
	return d, nil
}

func (d *data) WithPrivateRepository(username string, repository string, token string) (repoPackageManager, error) {
	err := migrate_packages_internal.DownloadExtractDeleteZip(username, repository, token)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *data) WithPublicRepository(username string, repository string) (repoPackageManager, error) {
	err := migrate_packages_internal.DownloadExtractDeleteZip(username, repository, "")
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *data) VersionsToMigrate() ([]string, error) {
	versions, err := d.Versions()
	if err != nil {
		return nil, err
	}
	var versionsToMigrate = []string{}
	for _, v := range versions {
		version := ToVersion(v)
		if version.EarlierThan(d.fromVersion) {
			continue
		}
		if d.fromVersion.Equals(version) {
			continue
		}
		if version.LaterThan(d.toVersion) {
			continue
		}
		versionsToMigrate = append(versionsToMigrate, version.String())
	}
	return versionsToMigrate, nil
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
	WithLocalFilesAt(projectRoot, pkgDir string) (PackageManager, error)
}

type repoPackageManager interface {
	// dir is the directory, relative to the project root, where the packages are
	WithPkgDir(dir string) PackageManager
}
