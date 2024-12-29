package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// downloadFile downloads a file from the given URL and saves it locally with the provided file name.
func downloadFile(url, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// unzipFile extracts the contents of a zip file to a specified destination directory.
func unzipFile(zipFileName, destDir string) error {
	zipReader, err := zip.OpenReader(zipFileName)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		// Strip the top-level directory
		pathParts := strings.SplitN(file.Name, "/", 2)
		if len(pathParts) < 2 {
			continue
		}
		relativePath := pathParts[1]
		outputFilePath := filepath.Join(destDir, relativePath)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(outputFilePath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directories: %w", err)
			}
		} else {
			if err := extractFile(file, outputFilePath); err != nil {
				return err
			}
		}
	}
	return nil
}

// extractFile extracts a single file from the zip archive to the specified output path.
func extractFile(file *zip.File, outputPath string) error {
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file inside zip: %w", err)
	}
	defer rc.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, rc)
	return err
}

// fetchLatestTag fetches the latest tag of a GitHub repository.
func fetchLatestTag(username, repoName string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", username, repoName)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch tags: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch tags: status code %d", resp.StatusCode)
	}

	var tags []struct{ Name string }
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found")
	}
	return tags[0].Name, nil
}

// FetchTagsFromGithub downloads and extracts the latest tagged release of a repository.
func FetchTagsFromGithub(username, repoName string) error {
	tag, err := fetchLatestTag(username, repoName)
	if err != nil {
		return err
	}
	fmt.Println("Latest Tag:", tag)

	zipURL := fmt.Sprintf("https://github.com/%s/%s/archive/refs/tags/%s.zip", username, repoName, tag)
	zipFileName := fmt.Sprintf("temp-%s.zip", tag)

	if err := downloadFile(zipURL, zipFileName); err != nil {
		return fmt.Errorf("failed to download release zip: %w", err)
	}
	defer os.Remove(zipFileName)

	exe, _ := os.Executable()
	destDir := filepath.Join(filepath.Dir(exe), "pkg", "@" + username, repoName)
	return unzipFile(zipFileName, destDir)
}

// FetchFromGithub downloads and extracts the main branch of a repository.
func FetchFromGithub(username, repoName string) error {
	zipURL := fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/main.zip", username, repoName)
	zipFileName := "temp.zip"

	if err := downloadFile(zipURL, zipFileName); err != nil {
		return fmt.Errorf("failed to download main branch zip: %w", err)
	}
	defer os.Remove(zipFileName)
	
	exe, _ := os.Executable()
	destDir := filepath.Join(filepath.Dir(exe), "pkg", "@" + username, repoName)
	return unzipFile(zipFileName, destDir)
}
