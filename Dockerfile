FROM golang:buster

LABEL maintainer="Maxim Goncharenko <kacejot@fex.net>"

WORKDIR /app

COPY . .

RUN go mod download && \
    go get github.com/golang/mock/mockgen && \
    go generate ./... && \
    go build -o server ./pkg/api

EXPOSE 8080
CMD ["./server"]
