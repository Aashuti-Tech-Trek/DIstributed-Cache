package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"distributed-cache/cache"
	"distributed-cache/hash"
)

type Server struct {
	cache *cache.Cache
	self  string
	ring  *hash.HashRing
}

func NewServer(capacity int, self string, replicas []string) *Server {
	ring := hash.NewHashRing(50)
	allNodes := []string{self}
	allNodes = append(allNodes, replicas...)
	ring.Add(allNodes...)

	return &Server{
		cache: cache.NewCache(capacity),
		self:  self,
		ring:  ring,
	}
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nodes := s.ring.GetNodes(key, 2)
	isOwner := false
	for _, n := range nodes {
		if n == s.self {
			isOwner = true
			break
		}
	}

	if !isOwner && len(nodes) > 0 {
		owner := nodes[0]
		target := fmt.Sprintf("http://%s/get?key=%s", owner, url.QueryEscape(key))
		resp, err := http.Get(target)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		return
	}

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
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ttl := int64(0)
	if t := r.URL.Query().Get("ttl"); t != "" {
		fmt.Sscanf(t, "%d", &ttl)
	}

	isReplicaWrite := r.Header.Get("X-Is-Replica") == "true"
	if isReplicaWrite {
		s.cache.SetWithTTL(key, value, ttl)
		w.Write([]byte("OK"))
		return
	}

	nodes := s.ring.GetNodes(key, 2)
	if len(nodes) > 0 && nodes[0] != s.self {
		owner := nodes[0]
		target := fmt.Sprintf("http://%s/set?key=%s&value=%s", owner, url.QueryEscape(key), url.QueryEscape(value))
		if ttl > 0 {
			target += fmt.Sprintf("&ttl=%d", ttl)
		}
		resp, err := http.Get(target)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		return
	}

	s.cache.SetWithTTL(key, value, ttl)

	if len(nodes) > 1 && nodes[0] == s.self {
		for i := 1; i < len(nodes); i++ {
			replica := nodes[i]
			go func(repl string) {
				target := fmt.Sprintf("http://%s/set?key=%s&value=%s", repl, url.QueryEscape(key), url.QueryEscape(value))
				if ttl > 0 {
					target += fmt.Sprintf("&ttl=%d", ttl)
				}
				req, err := http.NewRequest("GET", target, nil)
				if err == nil {
					req.Header.Set("X-Is-Replica", "true")
					http.DefaultClient.Do(req)
				}
			}(replica)
		}
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
