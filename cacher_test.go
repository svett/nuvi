package nuvi_test

import (
	"fmt"
	"strings"

	"gopkg.in/redis.v4"

	"github.com/svett/nuvi"
	"github.com/svett/nuvi/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RedisCacher", func() {
	var (
		client *fakes.FakeRedisClient
		cacher *nuvi.RedisCacher
	)

	BeforeEach(func() {
		client = &fakes.FakeRedisClient{}
		cacher = &nuvi.RedisCacher{
			Key:    "NEWS_XML",
			Client: client,
			Logger: &fakes.FakeLogger{},
		}
	})

	It("appends the content to a list", func() {
		client.LLenReturns(redis.NewIntCmd())
		cacher.Cache(strings.NewReader("Financial Times"))
		Expect(client.LLenCallCount()).To(Equal(1))

		key, values := client.LPushArgsForCall(0)
		Expect(key).To(Equal("NEWS_XML"))
		Expect(values).To(HaveLen(1))
		Expect(values).To(ContainElement("Financial Times"))
	})

	Context("when io.Reader fails", func() {
		It("does not caches the content", func() {
			reader := &fakes.FakeReadCloser{}
			reader.ReadReturns(0, fmt.Errorf("Oh no!"))
			cacher.Cache(reader)
			Expect(client.LLenCallCount()).To(Equal(0))
		})
	})
})
