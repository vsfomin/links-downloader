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

func (d MockDownloader) Download(msg string) ([]byte, error) {
	d.OnDownload(msg)
	return nil, nil
}

func TestReceiveMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newQueue := &MockQueue{
		Ch: make(chan string, 5),
	}
	actual := "hello world"
	newQueue.AddMessage(actual)
	newDownload := &MockDownloader{}
	newWorker := Worker{newQueue, newDownload}
	newDownload.OnDownload = func(actual string) {
		close(newQueue.Ch)
		if actual != "hello world" {
			t.Errorf("\"%v\" not equal \"%v\"", actual, "hello world")
		}
	}
	newWorker.StartReceiveMessages(ctx)
}

func TestReceiveMessagetoCloseChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newQueue := &MockQueue{
		Ch: make(chan string, 5),
	}
	close(newQueue.Ch)
	newDownload := &MockDownloader{}
	newWorker := Worker{newQueue, newDownload}
	err := newWorker.StartReceiveMessages(ctx)
	assert.Nil(t, err)
}

func TestCloseContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newQueue := &MockQueue{
		Ch: make(chan string),
	}
	newDownload := &MockDownloader{}
	newWorker := Worker{newQueue, newDownload}
	resCh := make(chan error)
	go func() {
		resCh <- newWorker.StartReceiveMessages(ctx)
	}()
	cancel()
	err := <-resCh
	expectedError := "Exit due to SIGNIN"
	assert.EqualErrorf(t, err, expectedError, "Error shoud be %v, but got %v", expectedError, err)
}
