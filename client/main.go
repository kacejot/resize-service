package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

func main() {
	// create a client (safe to share across requests)
	client := graphql.NewClient("http://127.0.0.1:8080/query")

	// resize(client)
	resizeExisting(client)
	// list(client)
}

func resize(client *graphql.Client) {
	file, err := os.Open("eye.png")
	if err != nil {
		log.Fatal(err)
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// make a request
	req := graphql.NewRequest(fmt.Sprintf(`
		mutation {
			resize(image: {filename: "file.png", contents: "%s"}, width: 600, height: 500) {
				id,
				original {
					imageLink
				},
				resized {
					imageLink
				}
			}
		}`, base64.StdEncoding.EncodeToString(buf)))

	// set header fields
	req.Header.Set("Authorization", "sampleuser")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", respData)
}

func resizeExisting(client *graphql.Client) {
	// make a request
	req := graphql.NewRequest(`
		mutation {
			resizeExisting(id: 9431, width: 600, height: 500) {
				id,
				original {
					imageLink
				},
				resized {
					imageLink
				}
			}
		}`)

	// set header fields
	req.Header.Set("Authorization", "sampleuser")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", respData)
}

func list(client *graphql.Client) {
	// make a request
	req := graphql.NewRequest(`
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
		}`)

	// set header fields
	req.Header.Set("Authorization", "sampleuser")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", respData)
}
