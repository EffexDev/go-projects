package main

import (
	"fmt"
	"net/http"
	"sync"
	"github.com/PuerkitoBio/goquery"
	"net/url"
)

//Edit worker count to scale
const workerCount = 12

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
)

type Result struct {
	URL    string
	Link   string
	Status string
}

func main() {
	seedURLs := []string{
		"https://buddytelco.com.au",
	}

	results := Orchestrator(seedURLs)
	for _, r := range results {
		fmt.Printf("[%s] %s\n", r.Status, r.Link)
	}
}

func Orchestrator(seedURLs []string) []Result {
	var results []Result
	scraperChannel := make(chan string)
	linksChannel := make(chan Result)

	var scraperWG, linksWG sync.WaitGroup

	//	Create worker pool for concurrently checking multiple seed URLs to find any links
	for i := 0; i < workerCount; i++ {
		scraperWG.Add(1)
		go func() {
			defer scraperWG.Done()
			for site := range scraperChannel {
				//Webscraper func needs to parse the html and return a []string of links
				links := webScraper(site)
				//After receiving the []string of links I need to range over it and then push the information to the linksChannel as instances of the struct so that the results can be collated nicely and printed
				for _, link := range links {
					linksChannel <- Result{URL: site, Link: link}
				}
			}
		}()
	}

	go func() {
		for _, u := range seedURLs {
			scraperChannel <- u
		}
		close(scraperChannel)
	}()

	//	Create worker pool for checking the status of the links scraped by the previous worker group
	for i := 0; i < workerCount; i++ {
		linksWG.Add(1)
		go func() {
			defer linksWG.Done()
			//Ranges over a slice of structs so we need to grab specifically the link item of each struct and send that to the linkchecker function. That then needs to send a HEAD request to each of the links and return the STATUS code which we assign to the Status item of the struct, then append each struct in full to the slice of structs declared at the top
			for s := range linksChannel {
				status := linkChecker(s.Link)
				s.Status = status
				results = append(results, s)
			}
		}()
	}

	go func() {
		scraperWG.Wait()
		close(linksChannel)
	}()

	linksWG.Wait()
	return results
}

func webScraper(site string) []string {
	var links []string
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing HTML for %s: %v", site, err)
		return links
	}

	base, err := url.Parse(site)
	if err != nil {
		fmt.Printf("Invalid base URL %s: %v\n", site, err)
		return links
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			u, err := url.Parse(href)
			if err == nil {
				fullURL := base.ResolveReference(u).String()
				links = append(links, fullURL)
			}
		}
	})
	return links
}

func linkChecker(link string) string {
	resp, err := http.Head(link)
	if err != nil {
		return "Failed"
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode != 403 {
		return fmt.Sprintf(ColorRed+"Dead - %d"+ColorReset, resp.StatusCode)
	}

	if resp.StatusCode == 403 {
		return fmt.Sprintf(ColorYellow + "Unauth" + ColorReset)
	}
	return fmt.Sprintf(ColorGreen+"OK - %d"+ColorReset, resp.StatusCode)
}
