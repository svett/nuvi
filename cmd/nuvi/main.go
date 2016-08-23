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
)

func init() {
	flag.StringVar(&url, "url", "", "Web URL")
	flag.StringVar(&redisAddr, "redisAddr", "localhost:6379", "Redis Addr")
	flag.StringVar(&redisPassword, "redisPassword", "", "Redis Password")
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

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(os.Stderr, "scraper: ", log.LstdFlags)
	scraper := &nuvi.Scraper{
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
