package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fomik2/links-downloader/internal/downloader"
	"github.com/fomik2/links-downloader/internal/rabbitmq"
	"github.com/fomik2/links-downloader/internal/worker"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Workers      string `yaml:"workers"`
		RabbitmqAddr string `yaml:"rabbitmqAddr"`
	}
)

func NewConfig() (Config, error) {
	cfg := Config{}
	data, err := os.Open("./config/config.yaml")
	if err != nil {
		return cfg, fmt.Errorf("open config file: %w", err)
	}
	defer data.Close()
	byteData, err := ioutil.ReadAll(data)
	if err != nil {
		return cfg, fmt.Errorf("read config file: %w", err)
	}
	err = yaml.Unmarshal(byteData, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("unmarshal config file: %w", err)
	}
	return cfg, err
}

func waitSignal(cancel context.CancelFunc, signalCh chan os.Signal) {
	sig := <-signalCh
	switch sig {
	case syscall.SIGKILL:
		os.Exit(0)
	case syscall.SIGINT:
		cancel()
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := NewConfig()
	fmt.Println(cfg)
	if err != nil {
		log.Println(err)
		return
	}

	workers, err := strconv.Atoi(cfg.Workers)
	if err != nil {
		log.Println(err)
		return
	}
	rabbitmqAddr := cfg.RabbitmqAddr
	fmt.Println(rabbitmqAddr)
	r, err := rabbitmq.NewRabbitMQ(rabbitmqAddr)
	d := downloader.NewDownloader()
	w := worker.NewWorker(r, d)

	if err != nil {
		log.Println(err)
		return
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	for i := 0; i <= workers; i++ {
		go w.Worker(ctx)
	}

	go waitSignal(cancel, signalChannel)

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-r.Forever
}
