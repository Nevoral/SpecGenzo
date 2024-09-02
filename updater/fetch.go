package updater

import (
	"io"
	"log"
	"net/http"
)

func fetchPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error fetching page: status code %d", resp.StatusCode)
	}
	page, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return page, nil
}
