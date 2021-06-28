module github.com/emar-kar/urlshortener

// +heroku goVersion go1.16 install ./cmd/main.go ./bin/url-shortener
go 1.16

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/go-redis/redis/v8 v8.8.3
	github.com/teris-io/shortid v0.0.0-20201117134242-e59966efd125
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)
