package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const MB_TO_BYTES = 1048576

func createRedisClient(redisHost string, redisPort int, useTLS bool) *redis.Client {
	addr := fmt.Sprintf("%s:%d", redisHost, redisPort)
	options := &redis.Options{
		Addr: addr,
	}
	if useTLS {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		log.Println("Using secured TLS connection to Redis.")
	}
	return redis.NewClient(options)
}

func populateData(redisHost string, redisPort int, numConnections int, initialKeySize int, delta int, useTLS bool) {
	rdb := createRedisClient(redisHost, redisPort, useTLS)
	defer rdb.Close()

	var wg sync.WaitGroup
	for i := 1; i <= numConnections; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", i)
			valueSize := (initialKeySize + (i-1)*delta) * MB_TO_BYTES
			value := strings.Repeat("x", valueSize)

			err := rdb.Set(context.Background(), key, value, 0).Err()
			if err != nil {
				log.Printf("Error setting key %s: %v", key, err)
			} else {
				log.Printf("Set key: %s with size: %d bytes", key, valueSize)
			}
		}(i)
	}
	wg.Wait()
	log.Println("All connections closed after populating data.")
}

func fetchDataSlowly(redisHost string, redisPort int, numConnections int, sleepTime int, useTLS bool) {
	var wg sync.WaitGroup
	for i := 1; i <= numConnections; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var conn net.Conn
			var err error
			if useTLS {
				conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", redisHost, redisPort), &tls.Config{
					InsecureSkipVerify: true,
				})
				log.Printf("Using secured TLS connection for key_%d", i)
			} else {
				conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", redisHost, redisPort))
			}
			if err != nil {
				log.Printf("Error connecting to Redis for key_%d: %v", i, err)
				return
			}
			defer conn.Close()

			key := fmt.Sprintf("key_%d\r\n", i)
			command := fmt.Sprintf("GET %s", key)
			_, err = conn.Write([]byte(command))
			if err != nil {
				log.Printf("Error sending command for key_%d: %v", i, err)
				return
			}

			log.Printf("Sent GET command for: %s", strings.TrimSpace(key))
			for j := 0; j < sleepTime*10; j++ {
				time.Sleep(100 * time.Millisecond)
				log.Printf("Sleeping for key_%d", i)
			}
		}(i)
	}
	wg.Wait()
	log.Println("Finished fetching data slowly.")
}

func main() {
	if len(os.Args) < 8 {
		log.Fatalf("Usage: %s <redis_host> <redis_port> <num_connections> <initial_key_size_MB> <delta_MB> <sleep_time_seconds> <noflush> <use_tls>", os.Args[0])
	}

	redisHost := os.Args[1]
	redisPort, _ := strconv.Atoi(os.Args[2])
	numConnections, _ := strconv.Atoi(os.Args[3])
	initialKeySize, _ := strconv.Atoi(os.Args[4])
	delta, _ := strconv.Atoi(os.Args[5])
	sleepTime, _ := strconv.Atoi(os.Args[6])
	noflush := os.Args[7] == "true"
	useTLS := os.Args[8] == "true"

	if !noflush {
		rdb := createRedisClient(redisHost, redisPort, useTLS)
		defer rdb.Close()
		_, err := rdb.FlushAll(context.Background()).Result()
		if err != nil {
			log.Fatalf("Error flushing Redis: %v", err)
		}
		log.Println("Flushed all Redis databases.")
	}

	log.Println("Starting population stage...")
	populateData(redisHost, redisPort, numConnections, initialKeySize, delta, useTLS)

	log.Println("Starting fetch stage...")
	fetchDataSlowly(redisHost, redisPort, numConnections, sleepTime, useTLS)
}
