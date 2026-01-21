package handler

import (
	"encoding/json"
	"net/http"

	"github.com/adityapatel-00/system-design/design-problems/easy/go/urlshortner/store"
)

type ShortenURLRequest struct {
	Url string `json:"url"`
}

type ShortenURLResponse struct {
	ShortenedUrl string `json:"shortened_url"`
	ShortCode    string `json:"short_code"`
	RedirectUrl  string `json:"redirect_url"`
}

type AnalyticsResponse struct {
	Url        string `json:"url"`
	VisitCount int    `json:"visit_count"`
	CreatedAt  string `json:"created_at"`
}

func SaveNewURL(s *store.URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation goes here
		var req ShortenURLRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Url == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		if req.Url[:4] != "http" && req.Url[:5] != "https" {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		shortCode := store.GenerateShortCode(req.Url)

		if s.CheckIfShortCodeExists(shortCode) {
			shortCode = store.GenerateShortCode(req.Url)
		}

		if err := s.SaveURL(shortCode, req.Url); err != nil {
			http.Error(w, "Failed to save URL mapping", http.StatusInternalServerError)
			return
		}

		resp := ShortenURLResponse{
			ShortenedUrl: req.Url,
			ShortCode:    shortCode,
			RedirectUrl:  "http://localhost:8080/" + shortCode,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func RedirectUrl(s *store.URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation goes here
		shortCode := r.URL.Path[len("/"):]
		urlData, err := s.GetURL(shortCode)
		if err != nil {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		s.IncrementVisitCount(shortCode)
		http.Redirect(w, r, urlData.OriginalURL, http.StatusTemporaryRedirect)
	}
}

func GetAnalytics(s *store.URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation goes here
		shortCode := r.URL.Path[len("/analytics/"):]
		urlData, err := s.GetURL(shortCode)
		if err != nil {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		resp := AnalyticsResponse{
			Url:        urlData.OriginalURL,
			VisitCount: urlData.VisitCount,
			CreatedAt:  urlData.CreatedAt.String(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
