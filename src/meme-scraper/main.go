package main

import (
	"flag"
	"log"
	"net/http"
)

const (
	BaseURL = "http://knowyourmeme.com/"
)

var (
	datafile   = flag.String("data", "data.db", "The file to store the data in")
	scrapeBool = flag.Bool("scrape", true, "Set to scrape otherwise just serve requests")
	url        = flag.String("url", BaseURL, "the url for know your meme")
	rate       = flag.Uint64("rate", 30000, "Number of milliseoncds between scrapes calls") // Keep this value high to prevent violation of TOS
	random     = flag.Int64("random", 10000, "Number of milliseconds to randomly stall to not make this look so robotic")
	page       = flag.Int("page", 1, "The page to start scraping at")
)

func main() {

	flag.Parse()

	saveChan = make(chan meme, 100)
	if err := initDB(*datafile); err != nil {
		log.Fatal(err)
	}
	if *scrapeBool {
		go scrape(*page)
	}
	go persist()

	http.HandleFunc("/meme", memeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
