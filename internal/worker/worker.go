package worker

import (
	"context"
	"fmt"
	"log"
)

type Queue interface {
	TakeMessage() (<-chan string, error)
	CloseConnections()
}

type Download interface {
	Download(url string) error
}

type Worker struct {
	q Queue
	d Download
}

func NewWorker(queue Queue, download Download) *Worker {
	newWorker := Worker{}
	newWorker.q = queue
	newWorker.d = download

	return &newWorker
}

func (w *Worker) Worker(ctx context.Context) error {
	msgs, err := w.q.TakeMessage()
	if err != nil {
		return fmt.Errorf("error while consume queue: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			w.q.CloseConnections()
			return nil
		case msg := <-msgs:
			log.Println(msg)
			w.d.Download(msg)
		}
	}
}
