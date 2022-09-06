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
	select {
	case <-ctx.Done():
		w.q.CloseConnections()
	default:
		for d := range msgs {
			log.Printf("Received a message: %s", d)
			w.d.Download(d)
		}
	}

	return nil
}
