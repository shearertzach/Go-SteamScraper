package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/gocolly/colly"
)

type SteamGame struct {
	Name        string `json:"name"`
	Link        string `json:"link"`
	Price       int    `json:"price"`
	Image       string `json:"image"`
	ReleaseDate string `json:"releaseDate"`
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector()
	d := color.New(color.FgGreen)
	f, _ := os.Create("data.json")

	f.WriteString("[")

	c.OnHTML("div#search_resultsRows > a", func(e *colly.HTMLElement) {
		// Find link using an attribute selector
		// Matches any element that includes href=""
		name := e.ChildText("div.responsive_search_name_combined > div.search_name > span")
		link := e.Attr("href")
		price, _ := strconv.Atoi(e.ChildAttr("div.responsive_search_name_combined > div.col.search_price_discount_combined.responsive_secondrow", "data-price-final"))
		image := e.ChildAttr("div.col.search_capsule > img", "src")
		release_date := e.ChildText("div.responsive_search_name_combined > div.search_released")

		steam_game := &SteamGame{
			Name:        name,
			Link:        link,
			Price:       price,
			Image:       image,
			ReleaseDate: release_date,
		}

		pathJson, _ := json.Marshal(steam_game)
		f.WriteString(string(pathJson) + ",")
		d.Println("Added:", name)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		stat, _ := f.Stat()
		f.WriteAt([]byte("]"), stat.Size()-1)
		f.Close()
	})

	c.Visit("https://store.steampowered.com/search/?filter=topsellers")
}
