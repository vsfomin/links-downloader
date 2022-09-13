package worker

import (
	"fmt"
)

type Queue interface {
	TakeMessage() (<-chan string, error)
}

type Download interface {
	Download(url string) error
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

func (w *Worker) Worker() error {
	msgs, err := w.queue.TakeMessage()
	if err != nil {
		return fmt.Errorf("error while consume queue: %w", err)
	}
	for {
		msg, ok := <-msgs
		if !ok {
			return nil
		}
		fmt.Println(msg)
		w.download.Download(msg)
	}
}
