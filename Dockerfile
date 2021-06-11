FROM ubuntu:focal
RUN apt-get update \
     && apt-get install -y \
	redis-tools
RUN mkdir /url-shortener
COPY ./bin /url-shortener
COPY ./web /url-shortener/web
CMD ["/url-shortener/main"]
