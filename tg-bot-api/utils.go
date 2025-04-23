package tg_bot_api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	filepathpkg "path/filepath"
)

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http.Get: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("[photo] Failed to close response body: %v", err)
		}
	}(resp.Body)

	dir := filepathpkg.Dir(filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Printf("[photo] Failed to close file: %v", err)
		}
	}(out)

	_, err = io.Copy(out, resp.Body)
	return err
}
