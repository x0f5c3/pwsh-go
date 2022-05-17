package pkg

import (
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"io/ioutil"
	"net/http"
	"runtime"
)

var HttpClient = http.DefaultClient

func Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "GITHUB API")
	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

var FileExt = func() string {
	os, ext, err := func() (string, string, error) {
		switch runtime.GOOS {
		case "windows":
			if runtime.GOARCH != "amd64" || runtime.GOARCH != "386" {
				return "win", "zip", nil
			}
			return "win", "msi", nil
		case "darwin":
			return "osx", "pkg", nil
		case "linux":
			return "linux", "tar.gz", nil
		default:
			return "", "", errors.New("unsupported os")
		}
	}()
	pterm.Fatal.WithFatal(true).PrintOnError(err)
	arch, err := func() (string, error) {
		switch runtime.GOARCH {
		case "amd64":
			return "x64", nil
		case "386":
			return "x86", nil
		case "arm":
			return "arm32", nil
		case "arm64":
			return "arm64", nil
		default:
			return "", errors.New("unsupported arch")
		}
	}()
	pterm.Fatal.WithFatal(true).PrintOnError(err)
	return fmt.Sprintf("%s-%s.%s", os, arch, ext)
}()
