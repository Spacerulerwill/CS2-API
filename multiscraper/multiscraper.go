package multiscraper

import (
	"net/http"
	"sync"
	"time"

	"gocasesapi/log"

	"github.com/PuerkitoBio/goquery"
)

// Multithreaded spawning of goroutines for scraping urls using  a fixed rate of documents
// per second
func MultiScrape[T any](urls []string, result map[string]T, perSecond int, callback func(*sync.Mutex, *goquery.Document, map[string]T)) {
	var wg sync.WaitGroup
	var scrapeWg sync.WaitGroup
	var mtx sync.Mutex

	requestsLeft := len(urls)
	requestsMade := 0
	var responses []*http.Response

	go continuallyScrapePages(&mtx, &scrapeWg, &responses, result, requestsLeft, callback)

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
			go makeConcurrentRequest(&wg, &mtx, urls[requestsMade+i], &requestsMade, &requestsLeft, &responses)
		}

		wg.Wait()
		elapsed := time.Since(start)
		time.Sleep((1000 - time.Duration(elapsed.Milliseconds())) * time.Millisecond)
	}

	scrapeWg.Wait()
}

func Http2Request(webUrl string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", webUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/117.0")

	res, err := client.Do(req)

	if err != nil || res.StatusCode != 200 {
		return nil, err
	}
	return res, nil
}

func makeConcurrentRequest(wg *sync.WaitGroup, mtx *sync.Mutex, webURL string, requestsMade *int, requestsLeft *int, outputResponses *[]*http.Response) {
	defer wg.Done()

	res, err := Http2Request(webURL)

	if err != nil {
		log.Error.Printf("%s", err.Error())
		mtx.Lock()
		defer mtx.Unlock()
		*requestsMade += 1
		*requestsLeft -= 1
	} else {
		mtx.Lock()
		defer mtx.Unlock()
		*requestsMade += 1
		*requestsLeft -= 1
		*outputResponses = append(*outputResponses, res)
	}
}

// Goroutine that will continually scrape from an array http responses with a callback until it has scraped a certain amount of times
func continuallyScrapePages[T any](mtx *sync.Mutex, scrapeWg *sync.WaitGroup, responses *[]*http.Response, result map[string]T, totalToScrape int, callback func(*sync.Mutex, *goquery.Document, map[string]T)) {
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
					go func() {
						defer scrapeWg.Done()
						callback(mtx, doc, result)
					}()
				}
				response.Body.Close()
				mtx.Unlock()
			}
			amountScraped = lenResp
		}
	}
}
