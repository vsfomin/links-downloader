package worker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockQueue struct {
	Ch chan string
}

func (q MockQueue) AddMessage(msg string) {
	q.Ch <- msg
}

func (q MockQueue) TakeMessage() (<-chan string, error) {
	return q.Ch, nil
}

type MockDownloader struct {
	OnDownload func(msg string)
}

func (d MockDownloader) Download(msg string) error {
	d.OnDownload(msg)
	return nil
}

func TestWorkerReceiveMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newQueue := &MockQueue{
		Ch: make(chan string),
	}
	newDownload := &MockDownloader{}
	newWorker := Worker{newQueue, newDownload}

	go func() {
		newWorker.Worker(ctx)
	}()
	expected := "hello world"
	newDownload.OnDownload = func(actual string) {
		if actual != expected {
			t.Errorf("\"%v\" not equal \"%v\"", actual, expected)
		}
	}
	newQueue.AddMessage(expected)

}

func TestWorkerCloseContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newQueue := &MockQueue{}
	newDownload := &MockDownloader{}
	newWorker := Worker{newQueue, newDownload}
	resCh := make(chan error)
	go func() {
		resCh <- newWorker.Worker(ctx)
	}()
	cancel()

	err := <-resCh
	assert.Nil(t, err)
}
