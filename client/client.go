package client

import (
	"fmt"
	"io"
	"net/http"

	"distributed-cache/hash"
)

type Client struct {
	ring *hash.HashRing
}

func NewClient(nodes []string) *Client {
	r := hash.NewHashRing(3)
	r.Add(nodes...)
	return &Client{ring: r}
}

func (c *Client) Get(key string) (string, error) {
	node := c.ring.Get(key)
	if node == "" {
		return "", fmt.Errorf("no nodes available")
	}
	url := fmt.Sprintf("http://%s/get?key=%s", node, key)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("key not found")
	}

	// Simple parsing, assuming JSON {"value":"..."}
	if len(body) > 10 {
		return string(body[10 : len(body)-2]), nil // crude extraction
	}
	return string(body), nil
}

func (c *Client) Set(key, value string) error {
	return c.SetWithTTL(key, value, 0)
}

func (c *Client) SetWithTTL(key, value string, ttl int64) error {
	node := c.ring.Get(key)
	if node == "" {
		return fmt.Errorf("no nodes available")
	}
	url := fmt.Sprintf("http://%s/set?key=%s&value=%s", node, key, value)
	if ttl > 0 {
		url += fmt.Sprintf("&ttl=%d", ttl)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}