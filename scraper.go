package main

import (
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
	MaxPage   int
	Schedule  *time.Ticker
	DBHandler *Database
	Logger    *log.Logger
}

func MakeScraper(interval, maxpage int, db *Database, logger *log.Logger) *Scraper {
	duration := time.Duration(interval) * time.Second
	return &Scraper{
		MaxPage:   maxpage,
		Schedule:  time.NewTicker(duration),
		Logger:    logger,
		DBHandler: db,
	}
}

// Start runs the scraper in blocking mode
// returns error if failed to initialize
func (sc *Scraper) Start() error {
	var snap *Snapshot
	sc.Logger.Println("[Scraper] Scraper Started.")
	sc.Logger.Println("[Scraper] Scraping snapshot for Database.")
	snap = sc.scrapeAll()
	sc.DBHandler.GetSnapshots(snap)
	for range sc.Schedule.C {
		sc.Logger.Println("[Scraper] Scraping...")
		snap = sc.scrapeAll()
		sc.Logger.Printf("[Scraper] Completed scraping: %d players scraped", len(snap.Players))
		sc.DBHandler.SaveSnapshot(snap)
	}
	return nil
}

// Stop stops the scraper
func (sc *Scraper) Stop() {
	sc.Schedule.Stop()
}

// scrapePage gets information about 1 page
func (sc *Scraper) scrapePage(client *http.Client, idx int) (*Page, error) {
	var page *Page = &Page{}

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
	err = page.Parse(responseData)
	if err != nil {
		return nil, err
	}
	return page, nil
}

// scrapeAll gets results from the leaderboards
func (sc *Scraper) scrapeAll() *Snapshot {
	var errors int
	var snap = &Snapshot{
		Timestamp: time.Now().Unix(),
		Players:   make(map[int]*Datapoint),
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: time.Second * 15}
	for i := 1; i <= sc.MaxPage; i++ {
		page, err := sc.scrapePage(client, i)
		if err != nil {
			errors++
			sc.Logger.Printf("[Scraper] Error scraping page %d: %s", i, err)
			if errors >= 3 {
				sc.Logger.Printf("[Scraper] Error limit reached, aborting scrape...")
				return nil
			}
			continue
		}
		for _, dp := range page.Data {
			snap.Players[dp.Id] = dp
		}
	}
	return snap
}
