package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.Default()
	db, _ := MakeDatabase("./phoenix.db", logger)
	scraper := MakeScraper(120, 1, db, logger)
	fmt.Println(db.Names)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go scraper.Start()
	<-done
	scraper.Stop()
}
