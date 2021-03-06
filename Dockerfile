FROM golang:1.16-alpine AS builder
RUN mkdir /url-shortener
COPY . /url-shortener
WORKDIR /url-shortener
RUN go mod tidy
RUN go build -o ./bin/main -v ./cmd/url-shortener/main.go

FROM alpine:latest
COPY --from=builder /url-shortener/bin /url-shortener/bin
COPY --from=builder /url-shortener/web /url-shortener/web
COPY --from=builder /url-shortener/README.md /url-shortener/
COPY --from=builder /url-shortener/docs /url-shortener/docs
WORKDIR /url-shortener
CMD ["./bin/main"]
