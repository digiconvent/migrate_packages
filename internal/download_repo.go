package migrate_packages_internal

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func DownloadRepoZip(owner, name, token string) error {
	targetDir := path.Join(os.TempDir(), "migrate_packages")
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/", owner, name)

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
		return errors.New("repository " + owner + "/" + name + " does not exist")
	}

	outFile, err := os.Create(path.Join(targetDir, "migrate_packages.zip"))
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	openCloseReader, _ := zip.OpenReader(outFile.Name())
	defer openCloseReader.Close()
	for _, file := range openCloseReader.File {
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(file.Name, file.Mode())
			if err != nil {
				return err
			}
		} else {
			srcFile, err := file.Open()
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.OpenFile(path.Join(targetDir, file.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
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
	return nil
}
