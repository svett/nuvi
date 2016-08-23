package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"gopkg.in/redis.v4"
)

// RedisPort is the default port
func RedisPort() int {
	port := os.Getenv("REDIS_SERVER_PORT")
	if port != "" {
		if portNo, err := strconv.Atoi(port); err == nil {
			return portNo
		}
	}
	return 6379
}

// NewRedisClient creates a new redis client
func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%d", RedisPort()),
		Password: "",
		DB:       0,
	})
}

// RedisRunner runs redis server
type RedisRunner struct {
	process *os.Process
	dir     string
}

// Start starts the server
func (runner *RedisRunner) Start(redisArgs ...string) error {
	command := exec.Command("redis-server", redisArgs...)
	dir, err := ioutil.TempDir("", "redis-client-test")
	if err != nil {
		return err
	}
	runner.dir = dir

	err = command.Start()
	if err != nil {
		return err
	}

	runner.process = command.Process
	port := RedisPort()
	if ok := IsPortOpen(port); !ok {
		return fmt.Errorf("Redis port %d is not open", port)
	}

	return nil
}

// Stop stops the server
func (runner *RedisRunner) Stop() error {
	err := runner.process.Kill()
	if err != nil {
		return err
	}

	client := NewRedisClient()
	defer client.Close()

	_, err = client.Ping().Result()
	if err == nil {
		return fmt.Errorf("Redis server is still running")
	}

	return os.RemoveAll(runner.dir)
}
