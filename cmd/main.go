package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/fomik2/links-downloader/internal/downloader"
	"github.com/fomik2/links-downloader/internal/rabbitmq"
	"github.com/fomik2/links-downloader/internal/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	log.Info().
		Str("method", "NewConfig").
		Msgf("Openning config file...")
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
		log.Info().
			Str("method", "waitSignal").
			Msgf("Close connection due to SIGTERM...")
		os.Exit(0)
	case syscall.SIGINT:
		log.Info().
			Str("method", "waitSignal").
			Msgf("Close connection due to SIGNIN...")
		cancel()
	}
}

func main() {
	//Global logging severity, change it if you don't want to see some logging ltvtl messages
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	var wg sync.WaitGroup

	cfg, err := NewConfig()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	workers, err := strconv.Atoi(cfg.Workers)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	rabbitmqAddr := cfg.RabbitmqAddr
	r, err := rabbitmq.NewRabbitMQ(rabbitmqAddr)
	defer r.CloseConnections()

	d := downloader.NewDownloader()

	w := worker.NewWorker(r, d)

	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go waitSignal(cancel, signalChannel)
	errCh := make(chan error)
	for i := 0; i <= workers; i++ {
		wg.Add(1)
		go func(i int) {
			log.Info().
				Str("method", "Worker").
				Msgf("Start worker: %v", i)
			errCh <- w.StartReceiveMessages(ctx)
			wg.Done()
		}(i)
	}
	go func() {
		for err := range errCh {
			log.Error().Err(err).Msg("")
		}
	}()
	wg.Wait()
}
