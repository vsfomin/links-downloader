package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Downloader struct{}

func (d *Downloader) DownloadFile(url string) (err error) {
	//Get filename
	sl := strings.Split(url, "/")
	filename := sl[len(sl)-1]
	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
