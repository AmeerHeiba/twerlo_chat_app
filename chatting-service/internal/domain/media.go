package domain

import "io"

type MediaUpload struct {
	File        io.Reader
	Filename    string
	ContentType string
	Size        int64
}

type MediaResponse struct {
	URL string `json:"url"`
}
