package config

import (
	"os"
	"strconv"
)

//MongoDbURI returns db URI
func MongoDbURI() string {
	return os.Getenv("MONGODB_URI")
}

//MongoDbDatabase returns db database
func MongoDbDatabase() string {
	return os.Getenv("MONGODB_DATABASE")
}

//MongoDbMaxPoolSize returns max pool size
func MongoDbMaxPoolSize() uint64 {
	m := os.Getenv("MONGODB_MAX_POOL_SIZE")
	if m == "" {
		return 100
	}
	n, _ := strconv.Atoi(m)
	return uint64(n)
}

//MongoDbMaxRetries returns db max retries
func MongoDbMaxRetries() int {
	m := os.Getenv("MONGODB_MAX_RETRIES")
	if m == "" {
		return 10
	}
	n, _ := strconv.Atoi(m)
	return int(n)
}
