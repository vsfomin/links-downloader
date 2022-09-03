package main

import (
	"log"

	"github.com/fomik2/links-downloader/internal/downloader"
	"github.com/fomik2/links-downloader/internal/rabbitmq"
)

func main() {
	r, err := rabbitmq.NewRabbitMQ()
	defer r.CloseConnections()
	if err != nil {
		log.Println(err)
		return
	}
	downloader := downloader.Downloader{}

	msgs, err := r.DeliverMessages()
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			downloader.DownloadFile(string(d.Body))
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-r.Forever

}
