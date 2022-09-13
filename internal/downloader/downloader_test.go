package downloader

import (
	"io"
	"log"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		"https://example.com",
		httpmock.NewStringResponder(200, "resp string"),
	)

	dwnd := NewDownloader()
	resp, err := dwnd.Download("https://example.com")
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, "resp string", string(body))
}
