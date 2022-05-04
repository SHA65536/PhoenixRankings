package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.Default()
	db, _ := MakeDatabase("./phoenix.db", logger)
	scraper := MakeScraper(600, 500, db, logger)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go scraper.Start()
	<-done
	scraper.Stop()
}
