package pkg

import (
	"encoding/json"
	"time"
)

type Latest struct {
	ReleaseDate time.Time `json:"ReleaseDate"`
	BlobName    string    `json:"BlobName"`
	ReleaseTag  string    `json:"ReleaseTag"`
}

func GetLatest() (*Latest, error) {
	resp, err := Get("https://aka.ms/pwsh-buildinfo-stable")
	if err != nil {
		return nil, err
	}
	var res Latest
	err = json.Unmarshal(resp, &res)
	return &res, err
}
