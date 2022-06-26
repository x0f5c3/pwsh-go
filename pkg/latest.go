package pkg

import (
	"encoding/json"
	"github.com/coreos/go-semver/semver"
	"time"
)

type LatestJson struct {
	ReleaseDate time.Time `json:"ReleaseDate"`
	BlobName    string    `json:"BlobName"`
	ReleaseTag  string    `json:"ReleaseTag"`
}

type Latest struct {
	ReleaseDate time.Time
	Version     *semver.Version
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

func LatestGithub() (*PWSHRelease, error) {
	resp, err := Get("https://api.github.com/repos/PowerShell/PowerShell/releases/latest")
	if err != nil {
		return nil, err
	}
	var rel Release
	err = json.Unmarshal(resp, &rel)
	if err != nil {
		return nil, err
	}
	return rel.Parse()
}
