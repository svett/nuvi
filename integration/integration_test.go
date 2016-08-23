package integration_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/redis.v4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/svett/nuvi/integration/utils"

	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Integration", func() {
	var (
		httpServer  *ghttp.Server
		redisClient *redis.Client
	)

	BeforeEach(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:%d", utils.RedisPort),
			Password: "",
			DB:       0,
		})

		httpServer = ghttp.NewServer()
		httpServer.AllowUnhandledRequests = true

		page, err := os.Open("../assets/index.html")
		Expect(err).NotTo(HaveOccurred())

		data, err := ioutil.ReadAll(page)
		Expect(err).NotTo(HaveOccurred())

		httpServer.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/"),
			ghttp.RespondWith(http.StatusOK, string(data)),
		))

		zip, err := os.Open("../assets/info.zip")
		Expect(err).NotTo(HaveOccurred())

		zipData, err := ioutil.ReadAll(zip)
		Expect(err).NotTo(HaveOccurred())

		httpServer.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/info.zip"),
			ghttp.RespondWith(http.StatusOK, string(zipData)),
		))
	})

	AfterEach(func() {
		httpServer.Close()
		redisClient.Del("NEWS_XML")
	})

	It("scrapes and caches the dired files", func() {
		session, err := runScraper(fmt.Sprintf("-url=%s", httpServer.URL()))
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		len := redisClient.LLen("NEWS_XML")
		Expect(len.Err()).NotTo(HaveOccurred())
		Expect(len.Val()).To(Equal(int64(3)))
	})

	Context("when the server is not available", func() {
		It("returns an error", func() {
			session, err := runScraper("-redisAddr=localhost:1234", "-url=www.example.com")
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
		})
	})

	Context("when the url is not provided", func() {
		It("returns an error", func() {
			session, err := runScraper()
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
		})
	})

	Context("when the password is wrong", func() {
		It("returns an error", func() {
			session, err := runScraper("-redisPassword=wrong", fmt.Sprintf("-url=%s", httpServer.URL()))
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
		})
	})

	Context("when the url is not accessible", func() {
		It("returns an error", func() {
			session, err := runScraper("-url=www.example.com")
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
		})
	})
})
