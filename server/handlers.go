package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func HandleSearch(res http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user_id")
	chatID := req.URL.Query().Get("chat_id")
	query := req.URL.Query().Get("query")

	if userID == "" || chatID == "" || query == "" {
		log.Printf("[server][search] ❌ Missing required parameters: user_id=%s, chat_id=%s query=%s", userID, chatID, query)
		http.Error(res, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Split query into tags and trim spaces
	rawTags := strings.Split(strings.ToLower(query), " ")
	tags := make([]string, 0)
	for _, tag := range rawTags {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}

	chatIdInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		log.Printf("[server][search] ❌ Invalid chat ID: %v", err)
		http.Error(res, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	// Get images by tags from repository
	images, err := SearchService.GetImagesByChatIdByTags(context.Background(), chatIdInt, tags)
	if err != nil {
		log.Printf("[server][search] ❌ Failed to get images: %v", err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(images); err != nil {
		log.Printf("[server][search] ❌ Failed to encode response: %v", err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[server][search] ✅ Found %d images for query: %s", len(images), query)
}

func HandleChats(res http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user_id")
	if userID == "" {
		log.Printf("[server][chats] ❌ Missing required parameter: user_id=%s", userID)
		http.Error(res, "Missing required parameter", http.StatusBadRequest)
		return
	}

	// Get chats by user ID from repository
	chats, err := SearchService.GetChatsByUserID(context.Background(), userID)
	if err != nil {
		log.Printf("[server][chats] ❌ Failed to get chats: %v", err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(chats); err != nil {
		log.Printf("[server][chats] ❌ Failed to encode response: %v", err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[server][chats] ✅ Found %d chats for user ID: %s", len(chats), userID)
}
