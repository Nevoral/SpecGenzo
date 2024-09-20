package mdn

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func fetchPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error fetching page: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func downloadData(filePath, url string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	body, err := fetchPage(url)
	if err != nil {
		return err
	}
	// Copy the response body to the file
	_, err = io.Copy(file, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

func readJsonFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
