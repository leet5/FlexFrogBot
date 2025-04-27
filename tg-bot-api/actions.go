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
func DownloadFile(saveDir, fileID, filename string) (string, error) {
	fileInfo, err := getFileInfo(fileID)
	if err != nil {
		log.Printf("[photo] ❌ Failed to get file path: %v", err)
		return "", err
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", token, fileInfo.FilePath)
	dest := filepath.Join(saveDir, filename)

	err = downloadFile(fileURL, dest)
	if err != nil {
		log.Printf("[tg_bot_api] ❌ Failed to download file: %v", err)
		return "", err
	}

	log.Printf("[tg_bot_api] ✅ Saved file to %s", dest)
	return dest, nil
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

// GetUserProfilePhoto retrieves the profile photo of a user by their userID.
func GetUserProfilePhoto(saveDir, filename string, userID int64) (string, error) {
	fileID, err := getUserPhotoFileID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user profile photo fileID: %w", err)
	}

	dest, err := DownloadFile(saveDir, fileID, filename)
	if err != nil {
		return "", fmt.Errorf("failed to download user profile photo: %w", err)
	}
	return dest, nil
}

// GetChatProfilePhoto retrieves the profile photo of a chat by its chatID.
func GetChatProfilePhoto(saveDir, filename string, chatID int64) (string, error) {
	fileID, err := getChatPhotoFileID(chatID)
	if err != nil {
		return "", fmt.Errorf("failed to get chat profile photo fileID: %w", err)
	}

	dest, err := DownloadFile(saveDir, fileID, filename)
	if err != nil {
		return "", fmt.Errorf("failed to download chat profile photo: %w", err)
	}
	return dest, nil
}
