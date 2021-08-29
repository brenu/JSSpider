package main

import (
	"bufio"
	"flag"
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

func verifyDomain(waitGroup *sync.WaitGroup, domain string, file *os.File) bool {
	var hasOutputFile bool = false

	if file != nil {
		hasOutputFile = true
	}

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

		if hasOutputFile == true {
			file.WriteString(path + "\n")
		}
	})

	c.Visit(finalURL)
	c.Wait()

	defer waitGroup.Done()
	return true
}

func main() {
	var outputPath string

	flag.StringVar(&outputPath, "o", "", "Path to the output file. Optional, and it uses an append approach, so that whenever you choose a file with content inside, it will not erase it.")
	flag.Parse()

	var file *os.File
	var err error

	if outputPath != "" {
		file, err = os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

		if err != nil {
			log.Fatal(err)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	var waitGroup sync.WaitGroup

	for scanner.Scan() {
		waitGroup.Add(1)
		go verifyDomain(&waitGroup, scanner.Text(), file)
	}

	waitGroup.Wait()

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
