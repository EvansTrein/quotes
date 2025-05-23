package main

import (
	"os"
	"os/signal"
	"quotes/config"
	"quotes/internal/repository"
	"quotes/internal/server"
	"quotes/pkg/logs"
	"sync"
	"syscall"
)

func main() {
	config.InitConfig()
	log := logs.InitLog(config.GetEnv(), false)

	var singltone sync.Once
	singltone.Do(func() {
		repository.InitQuoteRepo()
	})

	server := server.NewServer(&server.ServerDeps{
		Logger: log,
	})

	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	if err := server.Stop(); err != nil {
		log.Error("error occurred when stopping the application", "error", err)
		panic(err)
	}
}
