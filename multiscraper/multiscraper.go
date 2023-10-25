package multiscraper

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

var (
	mtx sync.Mutex
)

// Multithreaded spawning of goroutines for scraping urls using  a fixed rate of documents
// per second
func MultiScrape[T callbackConstraint](urls []string, result T, perSecond int, callback func(*goquery.Document, T)) {
	var wg sync.WaitGroup

	requestsLeft := len(urls)
	requestsMade := 0
	var responses []*http.Response

	go continuallyScrapePages(&responses, result, requestsLeft, callback)

	for requestsLeft > 0 {
		start := time.Now()
		requestsToMake := 0

		if requestsLeft < perSecond {
			requestsToMake = requestsLeft
		} else {
			requestsToMake = perSecond
		}

		wg.Add(requestsToMake)

		for i := 0; i < requestsToMake; i++ {
			go makeConcurrentRequest(urls[requestsMade+i], &wg, &requestsMade, &requestsLeft, &responses)
		}

		wg.Wait()
		elapsed := time.Since(start)
		time.Sleep((1000 - time.Duration(elapsed.Milliseconds())) * time.Millisecond)
	}
}

func Http2Request(webUrl string) *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", webUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/117.0")

	res, err := client.Do(req)

	// If url could not be opened, we inform the channel chFailedUrls:
	if err != nil || res.StatusCode != 200 {
		log.Error().Msg(fmt.Sprintf("%s: Failed to HTTP request %s: %d", err.Error(), webUrl, res.StatusCode))
	}
	return res
}

func makeConcurrentRequest(webURL string, wg *sync.WaitGroup, requestsMade *int, requestsLeft *int, outputResponses *[]*http.Response) {
	defer wg.Done()

	res := Http2Request(webURL)

	mtx.Lock()
	*requestsMade += 1
	*requestsLeft -= 1
	*outputResponses = append(*outputResponses, res)
	mtx.Unlock()
}

// Goroutine that will continually scrape from an array http responses with a callback until it has scraped a certain amount of times
func continuallyScrapePages[T callbackConstraint](responses *[]*http.Response, result T, totalToScrape int, callback func(*goquery.Document, T)) {
	amountScraped := 0
	for amountScraped != totalToScrape {
		mtx.Lock()
		lenResp := len(*responses)
		mtx.Unlock()
		if lenResp > amountScraped {
			for i := 0; i < lenResp-amountScraped; i++ {
				response := (*responses)[amountScraped+i]
				if response.StatusCode == 200 {
					doc, err := goquery.NewDocumentFromReader(response.Body)
					if err != nil {
						log.Err(err)
					}
					go callback(doc, result)
				}
				response.Body.Close()
			}
			amountScraped = lenResp
		}
	}
}
