package otomoto

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/thegroobi/web-listing-scrapper/models"
)

var (
	listing           models.CarListing
	power             uint64
	mileage           uint64
	pageCount         = 1
	mutex             sync.Mutex
	currentPage       = 1
	anchorsFound      = 0
	totalAnchorsFound = 0
	notFoundOffers    = 0
)

func ScrapArticles(link string) {
	c, offerCollector := createCollectors()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL: ", r.URL.String())
	})

	c.OnHTML("article[data-id]", func(e *colly.HTMLElement) {
		articleHandler(e, offerCollector)
	})

	if currentPage == 1 {
		c.OnHTML("li[data-testid=pagination-list-item]", paginationHandler)
	}

	offerCollector.OnHTML("div[data-testid=summary-info-area]", offerSummaryHandler)
	offerCollector.OnHTML("div[data-testid=content-details-section]", offerDetailsHandler)

	offerCollector.OnScraped(func(e *colly.Response) {
		saveAllListings()
	})

	offerCollector.OnError(func(r *colly.Response, err error) {
		errorHandler(r, err)
		notFoundOffers++
		fmt.Printf("Offers missing %d\n", notFoundOffers)
	})
	c.OnError(errorHandler)

	c.OnScraped(func(r *colly.Response) {
		onScrapedHandler(r, link, offerCollector)
	})

	err := c.Visit(link)
	if err != nil {
		log.Panic("Error visiting the website: ", err.Error())
	}

	c.Wait()
	offerCollector.Wait()
}

func createCollectors() (*colly.Collector, *colly.Collector) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.AllowedDomains("otomoto.pl", "www.otomoto.pl"),
		colly.Async(true),
	)
	c.WithTransport(transport)

	offerCollector := c.Clone()
	offerCollector.WithTransport(transport)

	return c, offerCollector
}

func articleHandler(e *colly.HTMLElement, c *colly.Collector) {
	href := e.ChildAttr("h1 a[href]", "href")
	if href == "" {
		log.Println("Href not found, trying again")
		e.Request.Retry()
		return
	}

	mutex.Lock()
	totalAnchorsFound++
	anchorsFound++
	mutex.Unlock()

	err := c.Visit(href)
	if err != nil {
		log.Println("Error visiting the offer for details: ", err.Error())
	}
}

func paginationHandler(e *colly.HTMLElement) {
	p, err := strconv.ParseInt(e.ChildText("a span"), 10, 16)
	if err != nil {
		log.Panic("Error on pagination handler: ", err.Error())
	}

	mutex.Lock()
	pageCount = int(p)
	mutex.Unlock()
}

func offerSummaryHandler(e *colly.HTMLElement) {
	pStr := strings.ReplaceAll(e.ChildText("h3.offer-price__number"), " ", "")
	pStr = strings.ReplaceAll(pStr, ",", ".")
	price, err := strconv.ParseFloat(pStr, 64)
	if err != nil {
		log.Println("Error parsing price: ", err.Error())
	}

	listing := getOrCreateListing(e.Request.URL.String())
	listing.Title = e.ChildText("div h1.offer-title")
	listing.Price = price
}

func offerDetailsHandler(e *colly.HTMLElement) {
	details := make(map[string]string)

	e.ForEach("div[data-testid=advert-details-item]", func(_ int, d *colly.HTMLElement) {
		key := d.DOM.Find("p").First().Text()
		value := d.DOM.Find("p").Last().Text()

		if value == "" || value == key {
			value = d.DOM.Find("a").Last().Text()
		}

		if key != "" && value != "" {
			details[key] = value
		}
	})

	year, err := strconv.ParseUint(details["Rok produkcji"], 10, 16)
	if err != nil {
		log.Println("Error parsing year:", err.Error())
		return
	}

	r := regexp.MustCompile(`\D+`)

	mileageStr := r.ReplaceAllString(details["Przebieg"], "")
	if mileageStr != "" {
		mileage, err = strconv.ParseUint(mileageStr, 10, 32)
		if err != nil {
			log.Println("Error parsing power: ", err.Error())
		}
	} else {
		mileage = 0
	}

	powerStr := r.ReplaceAllString(details["Moc"], "")
	if powerStr != "" {
		power, err = strconv.ParseUint(powerStr, 10, 16)
		if err != nil {
			log.Println("Error parsing power: ", err.Error())
		}
	} else {
		power = 0
	}

	listing := getOrCreateListing(e.Request.URL.String())
	listing.BodyType = details["Typ nadwozia"]
	listing.Gearbox = details["Skrzynia biegów"]
	listing.FuelType = details["Rodzaj paliwa"]
	listing.Color = details["Kolor"]
	listing.Version = details["Wersja"]
	listing.Power = uint16(power)
	listing.Year = uint16(year)
	listing.Mileage = uint32(mileage)

	if listing.Gearbox == "" || listing.FuelType == "" || listing.Color == "" || listing.BodyType == "" {
		e.Request.Retry()
		log.Println("Getting listing data failed, retrying...")
		return
	}
}

func errorHandler(r *colly.Response, err error) {
	if r.StatusCode == 0 || r.StatusCode == 429 || r.StatusCode >= 500 || r.StatusCode == 408 {
		time.Sleep(5 * time.Second)
		r.Request.Retry()
		return
	} else if r.StatusCode == 404 {
		log.Println("Page not found: ", r.Request.URL)
		return
	}
	log.Panicf("Request URL: %s failed with statusCode: %d\nError: %s", r.Request.URL, r.StatusCode, err.Error())
}

func onScrapedHandler(r *colly.Response, link string, offerCollector *colly.Collector) {
	fmt.Println("\nFinished browsing URL: ", r.Request.URL)
	fmt.Println("Total Anchors found: ", totalAnchorsFound)

	if anchorsFound != 32 && currentPage < pageCount {
		log.Printf("Expected 32 links, but found %d. Retrying page %d\n", anchorsFound, currentPage)

		anchorsFound = 0

		err := r.Request.Retry()
		if err != nil {
			log.Println("Error retrying page ", currentPage)
			log.Println(err.Error())
		}
		return
	}

	anchorsFound = 0
	currentPage++
	if currentPage <= pageCount {
		nextPage := fmt.Sprintf("%s?page=%d", link, currentPage)
		err := r.Request.Visit(nextPage)
		if err != nil {
			log.Println("Error visiting page ", nextPage)
			log.Print(err.Error())
		}
	}
	if pageCount == currentPage {
		offerCollector.Wait()
	}
}
