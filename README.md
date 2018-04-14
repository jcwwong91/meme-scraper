# meme-scraper
This is a scraper for knowyourmeme for a job interview.  In order for this to
work, please make sure you have access to the internet

System requirements:
- Linux
- gcc installed
- SQLite3 installed (usually installed by default)

TLDR:
1. make tools


Querying:
There is only one endpoint to query the information, a HTTP get request to the
"/meme" endpoint.  By default, this will return all of the memes that are
currently stored in the database
