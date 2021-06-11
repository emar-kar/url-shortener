FROM golang:1.16-alpine AS builder
RUN mkdir /url-shortener
COPY . /url-shortener
WORKDIR /url-shortener
RUN go mod tidy
RUN go build -o ./bin/main -v ./cmd/main.go

FROM alpine:latest
RUN apk --update add redis 
COPY --from=builder /url-shortener /url-shortener
WORKDIR /url-shortener
CMD ["./bin/main"]
