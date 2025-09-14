package migrate_packages_test

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// remove all the test databases
	err := os.RemoveAll(DataFolder)
	if err != nil {
		fmt.Println(err)
	}
	m.Run()
}
