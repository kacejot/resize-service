// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Image struct {
	ImageLink string `json:"imageLink"`
	ExpiresAt string `json:"expiresAt"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type ImageInput struct {
	Filename string `json:"filename"`
	Contents string `json:"contents"`
}

type ResizeResult struct {
	ID       string `json:"id"`
	Original *Image `json:"original"`
	Resized  *Image `json:"resized"`
}
