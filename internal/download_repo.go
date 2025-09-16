package migrate_packages_internal

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// no need to test this it works trust me bro
func DownloadExtractDeleteZip(owner, name, token string, verbose bool) error {
	targetDir := path.Join(os.TempDir(), "migrate_packages")
	targetZip := path.Join(os.TempDir(), "migrate_packages.zip")
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/main", owner, name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("could not reach " + url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("repository " + owner + "/" + name + " does not exist: " + string(body))
	}

	if verbose {
		fmt.Println("[Migrate Packages] Downloading " + url)
	}
	outFile, err := os.Create(targetZip)
	if err != nil {
		if verbose {
			fmt.Println("[Migrate Packages] Could not create " + targetZip + ": " + err.Error())
		}
		return err
	}
	defer outFile.Close()
	defer os.Remove(targetZip)

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	openCloseReader, _ := zip.OpenReader(outFile.Name())
	defer openCloseReader.Close()

	var prefix string = ""
	for _, file := range openCloseReader.File {
		if prefix == "" {
			prefix = file.Name
			if verbose {
				fmt.Println("[Migrate Packages] Setting prefix to " + prefix + " (we want to have the zipped files, not nested inside another folder)")
			}
		}
		fileName, _ := strings.CutPrefix(file.Name, prefix)
		var reset = "\033[0m"
		var red = "\033[31m"
		var green = "\033[32m"
		target := path.Join(targetDir, fileName)
		if verbose {
			fmt.Println("[Migrate Packages] " + red + prefix + green + targetDir + reset + fileName)
		}
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(target, file.Mode())
			if err != nil {
				return err
			}
		} else {
			srcFile, err := file.Open()
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}
		}
	}

	err = os.Remove(targetZip)

	if err != nil {
		if verbose {
			fmt.Println("[Migrate Packages] Could not clean up " + targetZip)
		}
	}
	return nil
}
