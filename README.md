# meme-scraper
This is a scraper for knowyourmeme for a job interview.  In order for this to
work, please make sure you have access to the internet

System requirements:
- Linux
- gcc installed
- SQLite3 installed (usually installed by default)

To build and run:
1. make tools
2. make install
3. make run

Querying:
There is only one endpoint to query the information, a HTTP get request to the
"/meme" endpoint.  This call will return a list of JSON objects describing the
different attributes of a meme.  The JSON objects will have the following format

{
  Name:        string             The name of the meme
  Src:         string             The URL that the information was scraped from
  Views:       int                The number of times that page was viewed
  Videos:      int                the number of videos related to the meme
  Images:      int                The number of images related to the meme
  Comments:    int                The number of comments the meme has
  Created:     RFC3339 timestamp  The time the meme was created
  LastUpdated: RFC3339 timestamp  The time the meme was last updated
}

In addition to getting all of the memes, the "/meme" call will also allow the
results to be filtered by their corresponding field.  This can be accomplished
with by attaching query parameters to the URL of the cal.  The possible query
parameters are:

  name:               string            meme's name must contain parameter value
  src:                string            meme's src must contain parameter value
  views :             int
  views_lt:           int
  views_gt :          int
  videos:             int
  videos_lt:          int
  videos_gt:          int
  images:             int
  images_lt:          int
  images_gt:          int
  comments:           int
  comments_lt:        int
  comments_gt:        int
  created:            RFC3339 timestamp
  created_before:     RFC3339 timestamp
  created_after:      RFC3339 timestamp
  lastUpdated:        RFC3339 timestamp
  lastUpdated_before: RFC3339 timestamp
  lastUpdated_after:  RFC3339 timestamp



If multiple parameters are specified, the filter will treat it as an "AND"
operation.  In other words an entry will only be returned if they match all
of the filter criteria

For example, if we
wanted to get all the memes that has more than 1000 views, the cal will be

  /meme?views_gt=1000

If I wanted all memes with less than 1000 views and was last updated before 2018
I would call

  /meme?views_lt=1000&lastUpdated_before=2018-01-01T00:00:00-00:00

Architecture:

This application has 3 goroutines that runs simultaneously.  The first goroutine
will crawl along the main page and then find any links related to a meme.  It
will then crawl and scrape information off those links and pass the information
to the second goroutine via a channel.  The second goroutine will read from the
channel and persist the meme information in a SQLite3 database.  The structure
of the table matches the structure of the JSON object in the previous section.
The thrid goroutine will listen on a port and satisfy any query requests
