package integration_test

import (
	"strings"

	"gopkg.in/redis.v4"

	"github.com/svett/nuvi"
	"github.com/svett/nuvi/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cacher", func() {
	var (
		redisClient *redis.Client
		cacher      *nuvi.RedisCacher
	)

	BeforeEach(func() {
		redisClient = newRedisClient()

		cacher = &nuvi.RedisCacher{
			Key:    "NEWS_XML",
			Client: redisClient,
			Logger: &fakes.FakeLogger{},
		}
	})

	AfterEach(func() {
		redisClient.Close()
	})

	It("caches the file only once", func() {
		cacher.Cache(strings.NewReader("dummy content"))

		len, err := redisClient.LLen("NEWS_XML").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(len).To(Equal(int64(1)))

		cacher.Cache(strings.NewReader("dummy content"))
		cacher.Cache(strings.NewReader("breaking news"))

		len, err = redisClient.LLen("NEWS_XML").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(len).To(Equal(int64(2)))

		text, err := redisClient.LPop("NEWS_XML").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(text).To(Equal("breaking news"))

		text, err = redisClient.LPop("NEWS_XML").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(text).To(Equal("dummy content"))
	})
})
