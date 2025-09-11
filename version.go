package migrate_packages

import (
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
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

func (v *Version) EarlierThan(b *Version) bool {
	if b.Major > v.Major {
		return false
	}
	if b.Minor > v.Minor {
		return false
	}
	if b.Patch > v.Patch {
		return false
	}
	return true
}

func (v *Version) MigrateTo(targetVersion *Version) packageManagerChoice {
	return &data{
		fromVersion: v,
		toVersion:   targetVersion,
	}
}
