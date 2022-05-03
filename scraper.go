package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const (
	URL     = "https://api.playphoenix.ca/rankings/%d"
	Retries = 3
)

type Scraper struct {
	MaxPage  int
	Schedule *time.Ticker
	Logger   *log.Logger
}

func MakeScraper(interval, maxpage int, logger *log.Logger) *Scraper {
	duration := time.Duration(interval) * time.Second
	return &Scraper{
		MaxPage:  maxpage,
		Schedule: time.NewTicker(duration),
		Logger:   logger,
	}
}

// Start runs the scraper in blocking mode
// returns error if failed to initialize
func (sc *Scraper) Start() error {
	sc.Logger.Println("Scraper Started")
	for range sc.Schedule.C {
		sc.Logger.Println("Scraping...")
		sc.scrapeAll()
	}
	return nil
}

// Stop stops the scraper
func (sc *Scraper) Stop() {
	sc.Schedule.Stop()
}

// scrapePage gets information about 1 page
func (sc *Scraper) scrapePage(client *http.Client, idx int) (*Page, error) {
	var page *Page

	req, _ := http.NewRequest("GET", fmt.Sprintf(URL, idx), nil)
	req.Header = http.Header{
		"User-Agent":   []string{"SHABot/0.1.0"},
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"*/*"},
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(responseData, &page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

// scrapeAll gets results from the leaderboards
func (sc *Scraper) scrapeAll() {
	var errors int
	var pages = make([]*Page, 0)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: time.Second * 15}
	for i := 1; i <= sc.MaxPage; i++ {
		page, err := sc.scrapePage(client, i)
		if err != nil {
			errors++
			sc.Logger.Printf("error scraping page %d: %s", i, err)
			if errors >= 3 {
				sc.Logger.Printf("error limit reached, aborting scrape...")
				return
			}
		}
		pages = append(pages, page)
	}
	sc.Logger.Printf("Scrape completed, %d pages downloaded", len(pages))
}
