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
	datafile = flag.String("data", "data.csv", "The file to store the data in")
	url      = flag.String("url", BaseURL, "the url for know your meme")
	rate     = flag.Uint64("rate", 1500, "Number of milliseoncds between scrapes calls") // Keep this low to prevent violation of TOS
)

func main() {

	flag.Parse()

	go scrape()

	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
