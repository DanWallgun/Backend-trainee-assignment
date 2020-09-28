FROM golang:1.15

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go build -mod vendor -o main cmd/urlshortener/main.go
CMD [ "/app/main" ]