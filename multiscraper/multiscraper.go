package multiscraper

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"gocasesapi/log"

	"github.com/PuerkitoBio/goquery"
)

var (
	mtx      sync.Mutex
	scrapeWg sync.WaitGroup
)

func WaitForCompletion() {
	scrapeWg.Wait()
}

// Multithreaded spawning of goroutines for scraping urls using  a fixed rate of documents
// per second
func MultiScrape[T callbackConstraint](urls []string, result T, perSecond int, callback func(*goquery.Document, T, *sync.WaitGroup)) {
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
		scrapeWg.Add(requestsToMake)

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
		log.Error.Println(fmt.Sprintf("%s: Failed to HTTP request %s: %d", err.Error(), webUrl, res.StatusCode))
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
func continuallyScrapePages[T callbackConstraint](responses *[]*http.Response, result T, totalToScrape int, callback func(*goquery.Document, T, *sync.WaitGroup)) {
	amountScraped := 0
	for amountScraped != totalToScrape {
		mtx.Lock()
		lenResp := len(*responses)
		mtx.Unlock()
		if lenResp > amountScraped {
			for i := 0; i < lenResp-amountScraped; i++ {
				mtx.Lock()
				response := (*responses)[amountScraped+i]

				if response.StatusCode == 200 {
					doc, err := goquery.NewDocumentFromReader(response.Body)
					if err != nil {
						log.Error.Println(err)
					}
					go callback(doc, result, &scrapeWg)
				}
				response.Body.Close()
				mtx.Unlock()
			}
			amountScraped = lenResp
		}
	}
}
