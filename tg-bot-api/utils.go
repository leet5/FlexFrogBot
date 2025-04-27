package tg_bot_api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	filepathpkg "path/filepath"
)

// downloadFile downloads a file from the given URL and saves it to the specified filepath.
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

// fetchUpdates performs the HTTP GET request to Telegram's getUpdates API.
func fetchUpdates(offset int) (*http.Response, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=60", token, offset)
	log.Printf("[tg_bot_api] 🔄 Fetching updates with offset=%d", offset)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %w", err)
	}

	return resp, nil
}

// getFileInfo gets file information from telegram servers by its ID
func getFileInfo(fileID string) (*FileInfo, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", token, fileID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("[tg_bot_api] Error closing body: %v", err)
		}
	}(resp.Body)

	var result struct {
		OK     bool     `json:"ok"`
		Result FileInfo `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}
	return &result.Result, nil
}

// decodeUpdates reads and parses the response body into a GetUpdatesResponse.
func decodeUpdates(body io.ReadCloser) (*GetUpdatesResponse, error) {
	defer func() {
		if err := body.Close(); err != nil {
			log.Printf("[tg_bot_api] ⚠️ Failed to close response body: %v", err)
		}
	}()

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	var result GetUpdatesResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	if !result.OK {
		log.Println("[tg_bot_api] ❗ Telegram API returned !OK response")
	}

	return &result, nil
}

func getUserPhotoFileID(userID int64) (string, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUserProfilePhotos?user_id=%d&limit=1", token, userID)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("http.Get: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("[tg_bot_api] ⚠️ Error closing body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("telegram returned status: %s", resp.Status)
	}

	var userPhotos UserProfilePhotos
	if err := json.NewDecoder(resp.Body).Decode(&userPhotos); err != nil {
		return "", fmt.Errorf("json.Decode: %w", err)
	}

	if !userPhotos.OK || userPhotos.Result.TotalCount == 0 || len(userPhotos.Result.Photos) == 0 || len(userPhotos.Result.Photos[0]) == 0 {
		return "", fmt.Errorf("no profile photo found")
	}

	lastIndex := len(userPhotos.Result.Photos[0]) - 1
	photo := userPhotos.Result.Photos[0][lastIndex]

	if photo.FileID == "" {
		return "", fmt.Errorf("profile photo has no file_id")
	}

	return photo.FileID, nil
}

func getChatPhotoFileID(chatID int64) (string, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getChat?chat_id=%d", token, chatID)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("http.Get: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("[tg_bot_api] ⚠️ Error closing body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("telegram returned status: %s", resp.Status)
	}

	var result ChatProfilePhoto
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("json.Decode: %w", err)
	}

	if !result.OK || result.Result.Photo == nil || result.Result.Photo.BigFileID == "" {
		return "", fmt.Errorf("no chat photo found")
	}

	return result.Result.Photo.BigFileID, nil
}
