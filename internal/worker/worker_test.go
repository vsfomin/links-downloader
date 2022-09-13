package worker

import (
	"testing"
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
	newQueue := &MockQueue{
		Ch: make(chan string, 5),
	}
	expected := "hello world5"
	newQueue.AddMessage(expected)
	newDownload := &MockDownloader{}
	newWorker := Worker{newQueue, newDownload}
	newDownload.OnDownload = func(actual string) {
		close(newQueue.Ch)
		if actual != "hello world" {
			t.Errorf("\"%v\" not equal \"%v\"", actual, expected)
		}
	}
	newWorker.Worker()
}

func TestWorkerCloseContext(t *testing.T) {
	newQueue := &MockQueue{
		Ch: make(chan string, 5),
	}
	expected := "hello world"
	newQueue.AddMessage(expected)
	newDownload := &MockDownloader{}
	newWorker := Worker{newQueue, newDownload}
	newDownload.OnDownload = func(actual string) {
		t.Errorf("Message received")
	}
	close(newQueue.Ch)
	newWorker.Worker()
}
