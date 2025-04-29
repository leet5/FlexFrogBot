package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func HandleSearch(res http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user_id")
	query := req.URL.Query().Get("query")

	if userID == "" || query == "" {
		log.Printf("[server][search] ❌ Missing required parameters: user_id=%s, query=%s", userID, query)
		http.Error(res, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Split query into tags and trim spaces
	tags := strings.Split(strings.ToLower(query), ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}

	// Get images by tags from repository
	images, err := SearchService.GetImagesByTags(context.Background(), tags)
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
