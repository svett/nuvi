package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"gopkg.in/redis.v4"

	"github.com/svett/nuvi"
)

var (
	url           string
	redisAddr     string
	redisPassword string
	maxConn       int
)

func init() {
	flag.StringVar(&url, "url", "", "Web URL")
	flag.StringVar(&redisAddr, "redis-addr", "localhost:6379", "Redis Addr")
	flag.StringVar(&redisPassword, "redis-password", "", "Redis Password")
	flag.IntVar(&maxConn, "max-parallel-download-conn", 4, "Max download connections")
}

func main() {
	flag.Parse()

	if url == "" {
		log.Fatal("The web url is missing. Please provide -url flag.")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
	defer client.Close()

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(os.Stderr, "scraper: ", log.LstdFlags)
	scraper := &nuvi.Scraper{
		MaxConn:    maxConn,
		Downloader: nuvi.HTTPDownloader(http.Get),
		Extractor: &nuvi.LinkExtractor{
			FileExt: ".zip",
			Logger:  logger,
		},
		ArchiveWalker: &nuvi.ZIPWalker{
			FileExt: ".xml",
			Logger:  logger,
		},
		Cacher: &nuvi.RedisCacher{
			Key:    "NEWS_XML",
			Client: client,
			Logger: logger,
		},
		Logger: logger,
	}

	if err = scraper.Scrape(url); err != nil {
		log.Fatal(err)
	}
}
