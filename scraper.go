package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type item struct {
	Name string `json:"name"`
	Price string `json:"price"`
	ImageURL string `json:"imgurl"`
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("j2store.net"),
	)

	var items []item

	c.OnHTML("div.col-sm-9 div[itemprop=itemListElement]", func(h *colly.HTMLElement) {
		item := item{
			Name: h.ChildText("h2.product-title"),
			Price: h.ChildText("div.sale-price"),
			ImageURL: h.ChildAttr("img", "src"),
		}
		items = append(items, item)
	})

	c.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
		nextPage := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(nextPage)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("http://j2store.net/demo/index.php/shop")
	
	content, err := json.Marshal(items)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("products.json", content, 0644)
}