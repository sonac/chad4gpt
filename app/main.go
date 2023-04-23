package main

import (
	"os"
	"os/signal"
	"syscall"

	"chad4gpt/app/tg"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err)
	}
	telegram := tg.NewTelegram()
	telegram.Start()
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go telegram.Start()

	log.Print("Press CTRL + C to stop programm")
	<-sigCh
	log.Print("Shutting down")
	os.Exit(0)
}
