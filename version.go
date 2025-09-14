package migrate_packages

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (a *Version) String() string {
	return fmt.Sprintf("%v.%v.%v", a.Major, a.Minor, a.Patch)
}

func (v *Version) IsUninitialised() bool {
	return v.Major == -1 && v.Minor == -1 && v.Patch == -1
}

func ToVersion(a string) *Version {
	segments := strings.Split(a, ".")
	if len(segments) != 3 {
		return nil
	}
	c := []int{}
	for _, segment := range segments {
		n, err := strconv.Atoi(segment)
		if err != nil {
			return nil
		}
		c = append(c, n)
	}
	return &Version{
		Major: c[0],
		Minor: c[1],
		Patch: c[2],
	}
}

func (a *Version) Equals(b *Version) bool {
	return a.Major == b.Major && a.Minor == b.Minor && a.Patch == b.Patch
}
func (a *Version) LaterThan(b *Version) bool {
	return b.EarlierThan(a)
}
func (a *Version) EarlierThan(b *Version) bool {
	if a.Equals(b) {
		return false
	}

	if a.Major > b.Major {
		return false
	} else if a.Major < b.Major {
		return true
	}

	if a.Minor > b.Minor {
		return false
	} else if a.Minor < b.Minor {
		return true
	}

	if a.Patch > b.Patch {
		return false
	}
	return true
}

func (a *Version) MigrateTo(targetVersion *Version) packageManagerChoice {
	return &data{
		fromVersion: a,
		toVersion:   targetVersion,
	}
}
