package config

import (
	"os"
	"strconv"
)

func MongoDbURI() string {
	return os.Getenv("MONGODB_URI")
}
func MongoDbDatabase() string {
	return os.Getenv("MONGODB_DATABASE")
}
func MongoDbMaxPoolSize() uint64 {
	m := os.Getenv("MONGODB_MAX_POOL_SIZE")
	if m == "" {
		return 100
	}
	n, _ := strconv.Atoi(m)
	return uint64(n)
}
func MongoDbMaxRetries() int {
	m := os.Getenv("MONGODB_MAX_RETRIES")
	if m == "" {
		return 10
	}
	n, _ := strconv.Atoi(m)
	return int(n)
}
