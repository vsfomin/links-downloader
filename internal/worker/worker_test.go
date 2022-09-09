package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockQueue struct {
	mock.Mock
}

type MockDownloader struct {
	mock.Mock
}

func (m *MockQueue) TakeMessage() (<-chan string, error) {
	//args := m.Called()
	strCh := make(chan string)
	strCh <- "some_url/some.txt"
	return strCh, nil
}

func (d *MockDownloader) Download(url string) error {
	if url == "some_url/some.txt" {
		return nil
	} else {
		return errors.New(url)
	}
	// args := d.Called(url)
	// return args.Error(1)
}

func TestWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	newQueue := &MockQueue{}
	newDownload := &MockDownloader{}
	//newDownload.On("Download").Return(nil)
	newWorker := Worker{newQueue, newDownload}
	err := newWorker.Worker(ctx)
	if err == nil {
		t.Errorf("test pass")
	} else {
		t.Errorf("test fail, want \"some_url/some.txt\", got: ")
	}

}
