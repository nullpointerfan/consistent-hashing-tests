package main

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// ConsistentHash represents the consistent hashing ring.
type ConsistentHash struct {
	hash     func(data []byte) uint32
	replicas int            // Number of virtual nodes per real node
	keys     []int          // Sorted list of hash values on the ring
	hashMap  map[int]string // Maps hash value to node name
}

// NewConsistentHash creates a new ConsistentHash instance.
func NewConsistentHash(replicas int, fn func(data []byte) uint32) *ConsistentHash {
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}
	return &ConsistentHash{
		hash:     fn,
		replicas: replicas,
		keys:     []int{},
		hashMap:  make(map[int]string),
	}
}

// Add adds a list of nodes to the hash ring.
func (c *ConsistentHash) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < c.replicas; i++ {
			hash := int(c.hash([]byte(node + strconv.Itoa(i))))
			c.keys = append(c.keys, hash)
			c.hashMap[hash] = node
		}
	}
	sort.Ints(c.keys)
}

// Get returns the node responsible for the given key.
func (c *ConsistentHash) Get(key string) string {
	if len(c.keys) == 0 {
		return ""
	}

	hash := int(c.hash([]byte(key)))

	// Find the first node on the ring that is greater than or equal to the key's hash.
	idx := sort.SearchInts(c.keys, hash)

	// If no such node is found, wrap around to the first node.
	if idx == len(c.keys) {
		idx = 0
	}
	return c.hashMap[c.keys[idx]]
}
