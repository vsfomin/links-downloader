package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Downloader struct{}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Download(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: 6 * time.Second,
	}
	//Get filename
	sl := strings.Split(url, "/")
	filename := sl[len(sl)-1]
	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("error while create file: %w", err)
	}
	defer out.Close()

	//Get the data

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error geting request: %w", err)
	}

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	//Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while io copy: %w", err)
	}

	return resp, nil
}
