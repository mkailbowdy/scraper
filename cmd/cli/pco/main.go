package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strings"
)

func main() {
	/* Flags (Configuration) */
	file := os.Args[len(os.Args)-1]
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	// Convert data []byte to string
	myString := string(data)
	// Split myString into []string
	codes := strings.Split(myString, "\n")

	for _, code := range codes {
		// Instantiate default collector
		c := colly.NewCollector(
			// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
			colly.AllowedDomains("www.pokemoncenter-online.com"),
		)
		c.OnHTML("input#availability", func(e *colly.HTMLElement) {
			stock := e.Attr("value")
			fmt.Printf("%s in stock? %s\n", code, stock)
		})

		c.Visit("https://www.pokemoncenter-online.com/" + code + ".html")
	}
}
