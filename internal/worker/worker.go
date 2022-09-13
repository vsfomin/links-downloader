package worker

import (
	"context"
	"fmt"
	"net/http"

	"errors"
)

type Queue interface {
	TakeMessage() (<-chan string, error)
}

type Download interface {
	Download(url string) (*http.Response, error)
}

type Worker struct {
	queue    Queue
	download Download
}

func NewWorker(queue Queue, download Download) *Worker {
	newWorker := Worker{}
	newWorker.queue = queue
	newWorker.download = download

	return &newWorker
}

func (w *Worker) StartReceiveMessages(ctx context.Context) error {
	msgs, err := w.queue.TakeMessage()
	if err != nil {
		return fmt.Errorf("error while consume queue: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			return errors.New("Exit due to SIGNIN")
		case msg, ok := <-msgs:
			if !ok {
				return errors.New("Channel was closed")
			}
			fmt.Println(msg)
			_, err := w.download.Download(msg)
			if err != nil {
				return err
			}

		}

	}
}
