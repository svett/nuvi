package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/redis.v4"
)

// RedisPort is the default port
const RedisPort = 6379

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
	if ok := IsPortOpen(RedisPort); !ok {
		return fmt.Errorf("Redis port %d is not open", RedisPort)
	}

	return nil
}

// Stop stops the server
func (runner *RedisRunner) Stop() error {
	err := runner.process.Kill()
	if err != nil {
		return err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%d", RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = client.Ping().Result()
	if err == nil {
		return fmt.Errorf("Redis server is still running")
	}

	return os.RemoveAll(runner.dir)
}
