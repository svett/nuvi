# nuvi

A web scraper for zip files. This utility inspect an html pages and downloads
any files that is an zip archive.  Note that it consider a file as a zip
archive if its anchor hyperlink contains `.zip` extension.

### Prerequisites

You should have the following dependencies installed and configured:

- `Golang 1.6`
- `Redis`

### Installation

```Golang
go get github.com/svett/nuvi/cmd/nuvi
```

### Usage

The `navi` binary can be executed with the following arguments:

- `url` the page address that will be inspected for zip files. **required**
- `redis-addr` the address of redis server. **optional**
- `redis-password` the password of redis server that the app is connecting to. **optional**
- `max-parallel-download-conn` the number of files downloaded in parallel. **optional**

```bash
$ nuvi -url=http_url_to_desired_page \
       -redis-addr=redis_server_host_and_port \
       -redis-password=redis_server_password \
       -max-parallel-download-conn \
```

#### Example

```bash
$ nuvi -url=http://feed.omgili.com/5Rh5AMTrc4Pv/mainstream/posts/
```

### Contribution

Getting the sources and all dependencies with the following git commands:

```bash
$ git clone https://github.com/svett/nuvi
$ git submodule update --init --recursive
```

In order to start contributing to the project, you should install
[ginkgo](http://github.com/onsi/ginkgo) and
[gomega](http://github.com/ons/gomega) package that are used in unit and
integration tests:

```bash
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega
```

You can run all unit and integration tests by executing the following script:

Note that you need `redis-server` installed. Every integration tests starts and
stops the server. Therefore, you should not have it running as a deamon.

The `redis-server` is running on port `6379`. If your instance is configured to
run on different port, you should set the environment variable `REDIS_SERVER_PORT`
before you execute the tests.

```bash
$ ./scripts/run_tests.sh
```

Also you can use `ginkgo` binary directly to execute the tests:

```bash
# Running the integration tests
$ ginkgo integration/
# Running the unit tests
$ ginkgo .
```

Presently the test coverage is **91.7%**.

### License

**MIT License**
