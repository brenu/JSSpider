package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/gocolly/colly/v2"
)

func getFinalURL(domain string) string {
	resp, err := http.Get(fmt.Sprintf("http://%s", domain))

	if err != nil {
		return ""
	}

	return resp.Request.URL.String()
}

func verifyDomain(waitGroup *sync.WaitGroup, domain string) bool {
	finalURL := getFinalURL(domain)

	if finalURL == "" {
		defer waitGroup.Done()
		return false
	}

	c := colly.NewCollector(
		colly.AllowedDomains(domain),
		colly.MaxDepth(2),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		re := regexp.MustCompile(`^https?://`)

		path := e.Attr("src")

		if re.MatchString(path) == false {
			path = fmt.Sprintf("%s", e.Request.AbsoluteURL(e.Attr("src")))
		}

		fmt.Println(path)
	})

	c.Visit(finalURL)
	c.Wait()

	defer waitGroup.Done()
	return true
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	var waitGroup sync.WaitGroup

	for scanner.Scan() {
		waitGroup.Add(1)
		go verifyDomain(&waitGroup, scanner.Text())
	}

	waitGroup.Wait()

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
