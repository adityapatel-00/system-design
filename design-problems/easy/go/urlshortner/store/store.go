package store

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrNotFound      = errors.New("URL not found")
	ErrAlreadyExists = errors.New("Short code already exists")
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type URLData struct {
	ShortCode   string
	OriginalURL string
	VisitCount  int
	CreatedAt   time.Time
}

type URLStore struct {
	mu   sync.RWMutex
	urls map[string]*URLData
}

func NewURLStore() *URLStore {
	return &URLStore{
		mu:   sync.RWMutex{},
		urls: make(map[string]*URLData),
	}
}

func (s *URLStore) SaveURL(shortCode, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.urls[shortCode]; exists {
		return ErrAlreadyExists
	}
	s.urls[shortCode] = &URLData{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		VisitCount:  0,
		CreatedAt:   time.Now(),
	}
	return nil
}

func (s *URLStore) GetURL(shortCode string) (*URLData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	urlData, exists := s.urls[shortCode]
	if !exists {
		return nil, ErrNotFound
	}
	return &URLData{
		ShortCode:   urlData.ShortCode,
		OriginalURL: urlData.OriginalURL,
		VisitCount:  urlData.VisitCount,
	}, nil
}

func (s *URLStore) IncrementVisitCount(shortCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	urlData, exists := s.urls[shortCode]
	if !exists {
		return ErrNotFound
	}
	s.urls[shortCode] = &URLData{
		ShortCode:   urlData.ShortCode,
		OriginalURL: urlData.OriginalURL,
		VisitCount:  urlData.VisitCount + 1,
	}
	return nil
}

func (s *URLStore) CheckIfShortCodeExists(shortCode string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.urls[shortCode]
	return exists
}

func GenerateShortCode(url string) string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
