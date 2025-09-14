package migrate_packages_test

import (
	"testing"

	"github.com/digiconvent/migrate_packages"
)

func v(a string) *migrate_packages.Version {
	return migrate_packages.ToVersion(a)
}

func TestVersions(t *testing.T) {
	// these versions are sorted from earliest to latest so index can be used to verify that they are later or earlier than one another
	vers := []*migrate_packages.Version{v("0.0.0"), v("0.0.1"), v("0.2.0"), v("0.3.4"), v("5.0.0"), v("6.0.7"), v("8.9.0"), v("10.11.12")}

	for i := range vers {
		for j := range vers {
			if i < j && (!vers[i].EarlierThan(vers[j]) || vers[i].LaterThan(vers[j])) {
				t.Fatal("expected", vers[i], "to be earlier than", vers[j])
			}
			if j < i && (vers[i].EarlierThan(vers[j]) || !vers[i].LaterThan(vers[j])) {
				t.Fatal("expected", vers[i], "not to be later than", vers[j])
			}
			if i == j && !vers[i].Equals(vers[j]) {
				t.Fatal("expected i to equal j")
			}
		}
	}
}
