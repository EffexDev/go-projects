package main

import (
	"fmt"
	"net/http"
	"sync"
	"golang.org/x/net/html"
)

/* Changes

Create a struct for results to print cleanly
Have main basically just print the results
Write a function to send the links to send links to channels and spawn workers
Write a function to parse the sites for links
Write a function to check the links

Requires complete refactor
*/


type Result struct {
	URL string
	Link string
	Status string
}

func main() {
	seedURLs := []string {
		"https://aussiebroadband.com.au",
	}

	results := Orchestrator(seedURLs)
	for _, r := range results {
		fmt.Printf("[%s] %s --> %s\n", r.URL, r.Link, r.Status)
	}
}

func Orchestrator(seedURLs []string) []Result {
	var results []Result
	scraperChannel := make(chan string)
	linksChannel := make(chan Result)

	var scraperWG, linksWG sync.WaitGroup

//	Create worker pool for concurrently checking multiple seed URLs to find any links
	for i := 0; i <=10; i++ {
		scraperWG.Add(1)
		go func() {
			defer scraperWG.Done()
			for _, site := range seedURLs {
//				Webscraper func needs to parse the html and return a []string of links
				links := webScraper(site)
//				After receiving the []string of links I need to range over it and then push the information to the 										linksChannel as instances of the struct so that the results can be collated nicely and printed
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
	for i := 0; i <= 10; i++ {
		linksWG.Add(1)
		go func() {
			defer linksWG.Done()
//			Ranges over a slice of structs so we need to grab specifically the link item of each struct and send that to the linkchecker 			function. That then needs to send a HEAD request to each of the links and return the STATUS code which we assign to the Status 			item of the struct, then append each struct in full to the slice of structs declared at the top
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
