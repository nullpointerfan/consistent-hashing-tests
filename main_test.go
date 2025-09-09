package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"strconv"
	"testing"

	consistent_hashing "github.com/ArchishmanSengupta/consistent-hashing"
)

func TestConsistentHash(t *testing.T) {
	ch := NewConsistentHash(3, nil)

	if ch.Get("123") != "" {
		t.Error()
	}

	hosts := []string{"host 1", "host 2", "host 3", "host 4"}
	for _, host := range hosts {
		ch.Add(host)
	}

	hostMap := make(map[string]string, 0)
	for i := range 10000 {
		key := fmt.Sprintf("key%d", i)
		host := ch.Get(key)
		t.Log(host, i)
		hostMap[key] = host
	}

	for i := range 10000 {
		key := fmt.Sprintf("key%d", i)
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
	})
	if err != nil {
		log.Fatal(err)
	}

	// Add hosts to the consistent hash ring
	hosts := []string{"host 1", "host 2", "host 3", "host 4"}
	ctx := context.TODO()
	for _, host := range hosts {
		err := ch.Add(ctx, host)
		if err != nil {
			log.Printf("Error adding host %s: %v", host, err)
		}
	}

	hostMap := make(map[string]string, 0)
	for i := range 10000 {
		key := fmt.Sprintf("key%d", i)
		host, err := ch.GetLeast(ctx, key)
		if err != nil {
			t.Error(err)
		}
		ch.IncreaseLoad(ctx, host)
		t.Log(host, i)
		hostMap[key] = host
	}

	for i := range 10000 {
		key := fmt.Sprintf("key%d", i)
		host, err := ch.Get(ctx, key)
		if err != nil {
			t.Error(err)
		}
		if host != hostMap[key] {
			t.Errorf("%v != %v", host, hostMap[key])
		}
	}

	loads := ch.GetLoads()
	fmt.Println("Current loads:")
	for host, load := range loads {
		fmt.Printf("%s: %d\n", host, load)
	}
}

func BenchmarkConsistentHash(b *testing.B) {
	ch := NewConsistentHash(3, nil)

	hosts := []string{"host 1", "host 2", "host 3", "host 4"}
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
	hosts := []string{"host 1", "host 2", "host 3", "host 4"}
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
