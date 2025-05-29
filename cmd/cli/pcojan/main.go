package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"strconv"
)

func main() {
	// Instantiate default collector
	var janCodes []byte
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.pokemoncenter-online.com"),
	)
	c.OnHTML("li[data-pid] ", func(e *colly.HTMLElement) {
		jan := e.Attr("data-pid")
		janCodes = append(janCodes, []byte(jan+"\n")...)
		fmt.Printf("jan code: %s\n", jan)
	})

	for i := 1; i < 10; i++ {
		err := c.Visit("https://www.pokemoncenter-online.com/search/?q=%E3%81%AC%E3%81%84%E3%81%90%E3%82%8B%E3%81%BF+Pokemon+fit+%E3%82%A2%E3%83%AD%E3%83%BC%E3%83%A9%E5%9C%B0%E6%96%B9&page=" + strconv.Itoa(i))
		if err != nil {
			break
		}
	}
	os.WriteFile("gen7.txt", janCodes, 0644)
}
