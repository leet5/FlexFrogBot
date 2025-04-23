package tg_bot_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// GetUpdatesChan starts a long-polling goroutine and returns a channel of updates.
func GetUpdatesChan() <-chan *Update {
	updates := make(chan *Update)

	go func() {
		var lastUpdateID int
		log.Println("[tg_bot_api] Starting long polling...")

		for {
			resp, err := fetchUpdates(lastUpdateID + 1)
			if err != nil {
				log.Printf("[tg_bot_api] ⚠️ Error fetching updates: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			updateResp, err := decodeUpdates(resp.Body)
			if err != nil {
				log.Printf("[tg_bot_api] ❌ Error decoding updates: %v", err)
				continue
			}

			if len(updateResp.Result) > 0 {
				log.Printf("[tg_bot_api] 📥 Received %d update(s)", len(updateResp.Result))
			}

			for _, update := range updateResp.Result {
				lastUpdateID = update.UpdateID
				updates <- &update
			}
		}
	}()

	return updates
}

// SendPayloadMessage sends a message payload to users and returns an error
func SendPayloadMessage(payload MessagePayload) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("http.Post: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("[tg_bot_api] <UNK> Error closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram returned status: %s", resp.Status)
	}

	return nil
}

// SendTextMessage sends a text message to users and returns an error
func SendTextMessage(chatID int64, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("http.Post: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[tg_bot_api] ⚠️ Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram returned status: %s", resp.Status)
	}

	log.Printf("[tg_bot_api] ✅ Text message sent to chat_id=%d", chatID)
	return nil
}

// DownloadFile downloads a file from telegram servers
func DownloadFile(saveDir, fileID, filename string) {
	fileInfo, err := getFileInfo(fileID)
	if err != nil {
		log.Printf("[photo] ❌ Failed to get file path: %v", err)
		return
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", token, fileInfo.FilePath)
	dest := filepath.Join(saveDir, filename)

	err = downloadFile(fileURL, dest)
	if err != nil {
		log.Printf("[tg_bot_api] ❌ Failed to download file: %v", err)
		return
	}

	log.Printf("[tg_bot_api] ✅ Saved file to %s", dest)
}

// getFileInfo gets file information from telegram servers by it ID
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

// IsUserAdmin checks whether is user with this ID is Administrator of chat with chatID
func IsUserAdmin(chatID, userID int64) (bool, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getChatAdministrators?chat_id=%d", token, chatID)

	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("http.Get: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("[tg_bot_api] Failed to close response body: %v", err)
		}
	}(resp.Body)

	var result struct {
		OK     bool `json:"ok"`
		Result []struct {
			User struct {
				ID int64 `json:"id"`
			} `json:"user"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("json.Decode: %w", err)
	}

	for _, admin := range result.Result {
		if admin.User.ID == userID {
			return true, nil
		}
	}
	return false, nil
}
