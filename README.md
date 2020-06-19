# resize-service
Ready to launch resize service with Dropbox storage

## Get

```sh
$ go get github.com/kacejot/resize-service
```

## Build
### By yourself
Make sure that `$GOPATH/bin` is in your `$PATH`.
Navigate to project directory.
Download dependencies:
```sh
$ go mod download
```

`go mod` does not install binaries to `$GOPATH/bin`, so we need to use `go get`:
```sh
$ go get github.com/golang/mock
$ go get github.com/99designs/gqlgen
```

Now we need to generate mocks and GraphQL code:
```sh
$ go generate ./...
```

Ready to build:
```sh
$ go build -o server pkg/api
```

### Using docker
```sh
$ docker build .
```

## Run
### On machine run
The instance of ArangoDB is required for work. Access to DB and cloud is configured through environment variables.
* `ARANGO_USER`
* `ARANGO_PASSWORD`
* `ARANGO_ENDPOINT`
* `DROPBOX_KEY`

After preparing working ArangoDB instance and setting environment variables properly, run the binary you built before by previous instructions:
```sh
$ ./server
```

### Containerized run
You can use `docker-compose` to run service. It will also run ArangoDB. You just need to configure the same environment variables for storage access through `docker-compose.yml`:
```yaml
...
services:
  api:
    depends_on:
      - arangodb
    build: ./
    environment:
      ARANGO_USER: root
      ARANGO_PASSWORD: rootpassword
      ARANGO_ENDPOINT: "http://arangodb:8529"
      DROPBOX_KEY: "<YOUR KEY>"
...
```

After this just execute:
```sh
$ docker-compose up
```

## How to use
After service is run you can send it HTTP request with GraphQL queries in body (that is how it works low level). In fact, you can generate client for any language using GraphQL scheme located in `pkg/api/graph`.
Query example:
```graphql
query {
    listProcessedImages {
    id,
    original {
        imageLink
        expiresAt
        width
        height
    },
    resized {
        imageLink
        expiresAt
        width
        height
    }
}
```
### GraphQL motivation
GraphQL has one giant feature: it allows you receive just what you want. For example, let us modfiy a query above:
```graphql
query {
    listProcessedImages {
    id,
    id,
    original {
        imageLink
        width
        height
    },
    resized {
        width
        height
    }
    id,
}
```
It still works, but returns 3 ids instead of one, because we asked this format. Also it returns not full info about images in case we don't need it.

## Authorization
It is not implemented yet, but all operations are bound to the user that is present in every GraphQL requst in `Authorization` header.
This is made for easy embedding authorization in future.
It is simple to add JWT instead of just user.

## What else to do
* API versioning
* Authorization
* Logging
* Beautiful errors
