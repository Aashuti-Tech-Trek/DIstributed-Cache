package hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashRing struct {
	replicas int
	keys     []int
	hashMap  map[int]string
}

func NewHashRing(replicas int) *HashRing {
	return &HashRing{
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

func (h *HashRing) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < h.replicas; i++ {
			hash := int(crc32.ChecksumIEEE([]byte(node + strconv.Itoa(i))))
			h.keys = append(h.keys, hash)
			h.hashMap[hash] = node
		}
	}
	sort.Ints(h.keys)
}

func (h *HashRing) Get(key string) string {
	if len(h.keys) == 0 {
		return ""
	}
	hash := int(crc32.ChecksumIEEE([]byte(key)))
	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	})
	if idx == len(h.keys) {
		idx = 0
	}
	return h.hashMap[h.keys[idx]]
}

func (h *HashRing) GetNodes(key string, count int) []string {
	if len(h.keys) == 0 {
		return nil
	}
	hash := int(crc32.ChecksumIEEE([]byte(key)))
	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	})
	if idx == len(h.keys) {
		idx = 0
	}

	nodes := make([]string, 0, count)
	seen := make(map[string]bool)
	visited := 0

	for len(nodes) < count && visited < len(h.keys) {
		node := h.hashMap[h.keys[idx]]
		if !seen[node] {
			seen[node] = true
			nodes = append(nodes, node)
		}
		idx++
		if idx == len(h.keys) {
			idx = 0
		}
		visited++
	}
	return nodes
}
