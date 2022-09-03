package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			downloader.DownloadFile(string(d.Body))
		}
	}()
	go func() {
		sig := <-signalChannel
		switch sig {
		case syscall.SIGKILL:
			os.Exit(0)
		case syscall.SIGTERM:
			r.CloseConnections()
			os.Exit(0)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-r.Forever
}
