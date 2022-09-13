package downloader

import (
	"io"
	"log"
	"net/http"
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
	resp1, _ := dwnd.Download("https://example.com")
	body1, err := io.ReadAll(resp1.Body)
	resp1Body := resp1.Body
	log.Println(resp1Body)
	if err != nil {
		log.Println(err)
	}
	log.Println("RESP1: ", string(body1))

	resp, _ := http.Get("https://example.com")
	respBody := resp1.Body
	log.Println(respBody)
	body, _ := io.ReadAll(resp.Body)
	log.Println("RESP: ", string(body))

	//assert.Nil(t, err)
	assert.Equal(t, "resp string", string(body))
}
