package downloader

import (
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
		"http://example.com",
		httpmock.NewStringResponder(200, "resp string"),
	)

	dwnd := NewDownloader()
	body, err := dwnd.Download("http://example.com")
	//	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		log.Println(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, "resp string", string(body))
}
