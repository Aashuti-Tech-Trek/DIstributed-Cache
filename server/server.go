package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"distributed-cache/cache"
)

type Server struct {
	cache    *cache.Cache
	replicas []string
}

func NewServer(capacity int, replicas []string) *Server {
	return &Server{
		cache:    cache.NewCache(capacity),
		replicas: replicas,
	}
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value, ok := s.cache.Get(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"value": value})
}

func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	ttl := int64(0)
	if t := r.URL.Query().Get("ttl"); t != "" {
		fmt.Sscanf(t, "%d", &ttl)
	}

	s.cache.SetWithTTL(key, value, ttl)

	// replicate asynchronously
	for _, replica := range s.replicas {
		replicaURL := fmt.Sprintf("http://%s/set?key=%s&value=%s", replica, url.QueryEscape(key), url.QueryEscape(value))
		if ttl > 0 {
			replicaURL += fmt.Sprintf("&ttl=%d", ttl)
		}
		go http.Get(replicaURL)
	}

	w.Write([]byte("OK"))
}

func (s *Server) SaveHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		filename = "cache.json"
	}
	if err := s.cache.SaveToDisk(filename); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("saved"))
}

func (s *Server) LoadHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		filename = "cache.json"
	}
	if err := s.cache.LoadFromDisk(filename); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("loaded"))
}
