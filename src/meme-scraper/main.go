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
	datafile = flag.String("data", "data.db", "The file to store the data in")
	url      = flag.String("url", BaseURL, "the url for know your meme")
	rate     = flag.Uint64("rate", 1500, "Number of milliseoncds between scrapes calls") // Keep this low to prevent violation of TOS
	page     = flag.Int("page", 1, "The page to start scraping at")
)

func main() {

	flag.Parse()

	saveChan = make(chan meme, 100)
	if err := initDB(*datafile); err != nil {
		log.Fatal(err)
	}
	go scrape(*page)
	go persist()

	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
