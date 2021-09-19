package github

import (
	"encoding/base64"
	"errors"
)

var ErrorUnknownContentType = errors.New("Unknown content type, cannot decode file")

type FileResults struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int    `json:"size"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"Content"`
	Encoding    string `json:"Encoding"`
}

func (fr *FileResults) DecodedContent() ([]byte, error) {
	switch fr.Encoding {
	case "base64":
		return base64.StdEncoding.DecodeString(fr.Content)
	}
	return nil, ErrorUnknownContentType
}
