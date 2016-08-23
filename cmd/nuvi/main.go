package main

import (
	"flag"
	"log"
	"net/http"

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

	scraper := &nuvi.Scraper{
		Downloader: nuvi.HTTPDownloader(http.Get),
		Extractor: &nuvi.LinkExtractor{
			FileExt: ".zip",
		},
		ArchiveWalker: &nuvi.ZIPWalker{
			FileExt: ".xml",
		},
		Cacher: &nuvi.RedisCacher{
			Key:    "NEWS_XML",
			Client: client,
		},
	}

	if err = scraper.Scrape(url); err != nil {
		log.Fatal(err)
	}
}
