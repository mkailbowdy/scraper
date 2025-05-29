package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Product struct {
	URL             string `json:"url"`
	Title           string `json:"title"`
	Image           string `json:"image"`
	Quantity        int    `json:"quantity"`
	Price           string `json:"price"`
	Color           string `json:"color"`
	Size            string `json:"size"`
	ItemPreferences string `json:"item_preferences"` // Notes field
	RateImage       string `json:"rate_image"`
}

func main() {
	/* Flags (Configuration) */
	website := os.Args[len(os.Args)-1]

	product := Product{}

	/* Make a *Collector for scraping */
	c := colly.NewCollector()

	/* Scraper instructions */
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	// Get Product's Title
	var titles []string
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		titles = append(titles, e.DOM.Text())
		product.Title = titles[0]
	})

	var prices []int
	// Get Product's Price
	c.OnHTML("span", func(e *colly.HTMLElement) {
		text := e.DOM.Text()
		match, err := regexp.MatchString("^[0-9,¥￥円$]+$", text)
		if err != nil {
			fmt.Printf(err.Error())
		}
		if !match {
			return
		}
		if strings.ContainsAny(text, ",¥￥円$") {
			text = strings.Trim(text, ",¥￥円$")
			removeComma := strings.Split(text, ",")
			text = strings.Join(removeComma, "")
		}
		price, err := strconv.Atoi(text)
		if err != nil {
			return
		}
		if price < 100 {
			return
		}
		prices = append(prices, price)
		product.Price = strconv.Itoa(prices[0])
	})

	c.OnHTML("main img", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	/* Run the scraper */
	c.Visit(website)
	fmt.Println(product)
}
