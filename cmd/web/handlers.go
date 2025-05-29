package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"html/template"
	"net/http"
	"os"
)

// Define a handler
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		//Shigotos: shigotos,
	}
	count := 8
	flag.Parse()
	// Instantiate default collector
	var products []byte

	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.goodsmile.com"),
	)

	var ranOnce bool
	c.OnHTML("div.c-top-product-list__unit", func(e *colly.HTMLElement) {
		if ranOnce {
			return
		}
		ranOnce = true

		e.ForEachWithBreak("a.c-top-product-list__item[href]", func(i int, h *colly.HTMLElement) bool {
			fmt.Println(i, count)
			if i >= count {
				return false
			}
			link := h.Attr("href")
			c.Visit(h.Request.AbsoluteURL(link))
			return true
		})
	})

	getDetails(c, "h1.b-product-info__title", &products)
	getDetails(c, "span.c-price__main", &products)
	getDetails(c, "p.b-product-info__note", &products)
	getDetails(c, "p[name]", &products)

	c.OnHTML("section.p-product__section a", func(e *colly.HTMLElement) {
		if e.Text == "パートナーショップ一覧" {
			fmt.Println("Parter? Yes!", e.Attr("href"))
		} else {
			fmt.Println("Partner? No!")
		}
	})

	c.OnHTML("div#specification", func(e *colly.HTMLElement) {
		e.ForEach("dl.b-outline-table__detail", func(_ int, dl *colly.HTMLElement) {
			term := dl.ChildText("dt h3")
			if term == "仕様" || term == "Specifications" {
				specText := dl.ChildText("dd p")
				products = append(products, []byte(specText+"\n"+"============================================\n\n")...)
				fmt.Printf("Specification:%s\n\n\n", specText)
				fmt.Print("============================================")
			}
		})
	})

	c.Visit("https://www.goodsmile.com/ja")
	c.Visit("https://www.goodsmile.com/en")
	err = os.WriteFile("ui/static/goodsmile_jp.txt", products, 0644)
	if err != nil {
		fmt.Println(err)
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}

}
func (app *application) gseng(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		//Shigotos: shigotos,
	}
	count := 8
	flag.Parse()
	// Instantiate default collector
	var products []byte

	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.goodsmile.com"),
	)

	var ranOnce bool
	c.OnHTML("div.c-top-product-list__unit", func(e *colly.HTMLElement) {
		if ranOnce {
			return
		}
		ranOnce = true

		e.ForEachWithBreak("a.c-top-product-list__item[href]", func(i int, h *colly.HTMLElement) bool {
			fmt.Println(i, count)
			if i >= count {
				return false
			}
			link := h.Attr("href")
			c.Visit(h.Request.AbsoluteURL(link))
			return true
		})
	})

	getDetails(c, "h1.b-product-info__title", &products)
	getDetails(c, "span.c-price__main", &products)
	getDetails(c, "p.b-product-info__note", &products)
	getDetails(c, "p[name]", &products)

	c.OnHTML("section.p-product__section a", func(e *colly.HTMLElement) {
		if e.Text == "パートナーショップ一覧" {
			fmt.Println("Parter? Yes!", e.Attr("href"))
		} else {
			fmt.Println("Partner? No!")
		}
	})

	c.OnHTML("div#specification", func(e *colly.HTMLElement) {
		e.ForEach("dl.b-outline-table__detail", func(_ int, dl *colly.HTMLElement) {
			term := dl.ChildText("dt h3")
			if term == "仕様" || term == "Specifications" {
				specText := dl.ChildText("dd p")
				products = append(products, []byte(specText+"\n"+"============================================\n\n")...)
				fmt.Printf("Specification:%s\n\n\n", specText)
				fmt.Print("============================================")
			}
		})
	})

	c.Visit("https://www.goodsmile.com/en")
	err = os.WriteFile("ui/static/goodsmile_en.txt", products, 0644)
	if err != nil {
		fmt.Println(err)
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}

}
func getDetails(c *colly.Collector, goquerySelector string, product *[]byte) {
	c.OnHTML(goquerySelector, func(e *colly.HTMLElement) {
		detail := e.DOM.Text()
		fmt.Printf("%s\n", detail)
		*product = append(*product, []byte(detail+"\n")...)
	})
}

//func (app *application) shigotoView(w http.ResponseWriter, r *http.Request) {
//	id, err := strconv.Atoi(r.PathValue("id"))
//	if err != nil || id < 1 {
//		http.NotFound(w, r)
//		return
//	}
//	shigoto, err := app.shigotos.Get(id)
//	if err != nil {
//		// Use our custom error models.ErrNoRecord. NOT sql.ErrorNoRows!
//		if errors.Is(err, models.ErrNoRecord) {
//			http.NotFound(w, r)
//		} else {
//			app.serverError(w, r, err)
//		}
//		return
//	}
//
//	files := []string{
//		"./ui/html/base.tmpl.html",
//		"./ui/html/partials/nav.tmpl.html",
//		"./ui/html/pages/view.tmpl.html",
//	}
//
//	ts, err := template.ParseFiles(files...)
//	if err != nil {
//		app.serverError(w, r, err)
//		return
//	}
//
//	// This templateData struct will hold all dynamic data we need for the html templates
//	data := templateData{
//		Shigoto: shigoto,
//	}
//
//	err = ts.ExecuteTemplate(w, "base", data)
//	if err != nil {
//		app.serverError(w, r, err)
//	}
//}
//
//func (app *application) shigotoCreate(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("Create job"))
//}
//
//func (app *application) shigotoCreatePost(w http.ResponseWriter, r *http.Request) {
//	companyName := "Samsung"
//	jobTitle := "QA Engineer"
//	category := "IT"
//	location := "Seoul, Korea"
//	employmentType := "Contract"
//	description := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
//	japaneseLevel := "N1"
//	englishLevel := "Native"
//	sponsorship := true
//
//	id, err := app.shigotos.Insert(companyName, jobTitle, category, location, employmentType, description, japaneseLevel, englishLevel, sponsorship)
//	if err != nil {
//		app.serverError(w, r, err)
//		return
//	}
//	http.Redirect(w, r, fmt.Sprintf("/shigoto/view/%d", id), http.StatusSeeOther)
//}
