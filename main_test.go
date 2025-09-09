package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	consistent_hashing "github.com/ArchishmanSengupta/consistent-hashing"
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
)

var (
	hosts = []string{"host 1", "host 2", "host 3", "host 4", "host 5", "host 6", "host 7", "host 8"}
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func TestConsistentHash(t *testing.T) {
	keys := make([]string, 0, 10000)
	for range 10000 {
		keys = append(keys, generateRandomString(64))
	}

	ch := NewConsistentHash(3, nil)

	if ch.Get("123") != "" {
		t.Error()
	}

	for _, host := range hosts {
		ch.Add(host)
	}

	hostMap := make(map[string]string, 0)

	for _, key := range keys {
		host := ch.Get(key)
		hostMap[key] = host
	}

	for _, key := range keys {
		host := ch.Get(key)
		if host != hostMap[key] {
			t.Errorf("%v != %v", host, hostMap[key])
		}
	}

	loadMap := make(map[string]int, 4)
	for _, host := range hostMap {
		loadMap[host] += 1
	}

	for key, value := range loadMap {
		t.Logf("Host: %v, Load: %v", key, value)
	}
}

func TestArchishmanSengConsistentHash(t *testing.T) {
	ch, err := consistent_hashing.NewWithConfig(consistent_hashing.Config{
		ReplicationFactor: 3,
		LoadFactor:        1.25,
	})
	if err != nil {
		log.Fatal(err)
	}

	keys := make([]string, 0, 10000)
	for range 10000 {
		keys = append(keys, generateRandomString(64))
	}

	// Add hosts to the consistent hash ring
	ctx := context.TODO()
	for _, host := range hosts {
		err := ch.Add(ctx, host)
		if err != nil {
			log.Printf("Error adding host %s: %v", host, err)
		}
	}

	hostMap := make(map[string]string, 0)
	for _, key := range keys {
		host, err := ch.GetLeast(ctx, key)
		if err != nil {
			t.Error(err)
		}
		ch.IncreaseLoad(ctx, host)
		hostMap[key] = host
	}

	loads := ch.GetLoads()
	fmt.Println("Current loads:")
	for host, load := range loads {
		fmt.Printf("%s: %d\n", host, load)
	}
}

type myMember string

func (m myMember) String() string {
	return string(m)
}

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	// you should use a proper hash function for uniformity.
	return xxhash.Sum64(data)
}

func TestBuraksezerConsistent(t *testing.T) {
	keys := make([]string, 0, 10000)
	for range 10000 {
		keys = append(keys, generateRandomString(64))
	}

	// Create a new consistent instance
	cfg := consistent.Config{
		PartitionCount:    103,
		ReplicationFactor: 5,
		Load:              1.25,
		Hasher:            hasher{},
	}
	c := consistent.New(nil, cfg)

	for _, host := range hosts {
		c.Add(myMember(host))
	}

	hostMap := make(map[string]string, 0)

	for _, key := range keys {
		host := c.LocateKey([]byte(key))
		hostMap[key] = host.String()
	}

	for _, key := range keys {
		host := c.LocateKey([]byte(key))
		if host.String() != hostMap[key] {
			t.Errorf("%v != %v", host, hostMap[key])
		}
	}

	loadMap := make(map[string]int, 4)
	for _, host := range hostMap {
		loadMap[host] += 1
	}

	sum := 0
	for key, value := range loadMap {
		t.Logf("Host: %v, Load: %v", key, value)
		sum += value
	}
	t.Logf("Sum: %v", sum)

	c.Remove(hosts[2])
	c.Remove(hosts[5])

	for _, key := range keys {
		host := c.LocateKey([]byte(key))
		hostMap[key] = host.String()
	}

	for _, key := range keys {
		host := c.LocateKey([]byte(key))
		if host.String() != hostMap[key] {
			t.Errorf("%v != %v", host, hostMap[key])
		}
	}

	loadMap = make(map[string]int, 4)
	for _, host := range hostMap {
		loadMap[host] += 1
	}

	sum = 0
	t.Log("After remove:")
	for key, value := range loadMap {
		t.Logf("Host: %v, Load: %v", key, value)
		sum += value
	}
	t.Logf("Sum: %v", sum)
}

func BenchmarkConsistentHash(b *testing.B) {
	ch := NewConsistentHash(3, nil)

	for _, host := range hosts {
		ch.Add(host)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ch.Get(strconv.FormatInt(int64(i), 10))
	}
}

func BenchmarkArchishmanSengConsistentHash(b *testing.B) {
	ch, err := consistent_hashing.NewWithConfig(consistent_hashing.Config{
		ReplicationFactor: 3,
		HashFunction:      fnv.New64a,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Add hosts to the consistent hash ring
	ctx := context.Background()
	for _, host := range hosts {
		err := ch.Add(ctx, host)
		if err != nil {
			log.Printf("Error adding host %s: %v", host, err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		host, err := ch.GetLeast(ctx, fmt.Sprintf("key%d", i))
		if err != nil {
			b.Error(err)
		}
		ch.IncreaseLoad(ctx, host)
	}
}
