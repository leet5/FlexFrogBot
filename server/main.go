package server

import (
	"context"
	"errors"
	"flex-frog-bot/services/interfaces"
	"log"
	"net/http"
	"time"
)

var (
	SearchService interfaces.SearchService
)

func RunServer(ctx context.Context, searchService interfaces.SearchService) {
	SearchService = searchService

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/search", WithCORS(HandleSearch))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: LogRequestsMiddleware(mux),
	}

	go func() {
		log.Println("[server] 🕑 Starting server on :8080 ...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[server] ⚠️ Error starting server: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("[server] ⚠️ Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[server] ❌ Error during server shutdown: %v", err)
	}
}

func LogRequestsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		log.Printf("[server][request] ⬇️ %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)

		log.Printf("[server][response] ⬆️ %s %s -> %d (%v)", r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
