package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gocolly/colly"
	// /"github.com/lunny/html2md"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func lastString(ss []string) string {
	return ss[len(ss)-2]
}
func writeIntoMD(md string, name string) {
	//	a := strconv.Itoa(rand.Intn(100))
	f, err := os.Create(name + ".md")

	check(err)
	defer f.Close()
	// n3, err := f.WriteString(md)
	// fmt.Printf("wrote %d bytes\n", n3)
	// f.Sync()

	w := bufio.NewWriter(f)
	n4, err := w.WriteString(md)
	fmt.Printf("wrote %d bytes\n", n4)
	w.Flush()
}
func main() {
	fileN := "index"
	c := colly.NewCollector(
		colly.AllowedDomains("golangbot.com"),
	)

	// // Callback for when a scraped page contains an article element
	// c.OnHTML("article", func(e *colly.HTMLElement) {
	// 	isEmojiPage := false

	// 	// Extract meta tags from the document
	// 	metaTags := e.DOM.ParentsUntil("~").Find("meta")
	// 	metaTags.Each(func(_ int, s *goquery.Selection) {
	// 		// Search for og:type meta tags
	// 		property, _ := s.Attr("property")
	// 		if strings.EqualFold(property, "og:type") {
	// 			content, _ := s.Attr("content")

	// 			// Emoji pages have "article" as their og:type
	// 			isEmojiPage = strings.EqualFold(content, "article")
	// 		}
	// 	})

	// 	if isEmojiPage {
	// 		// Find the emoji page title
	// 		fmt.Println("Emoji: ", e.DOM.Find("h1").Text())
	// 		// Grab all the text from the emoji's description
	// 		fmt.Println(
	// 			"Description: ",
	// 			e.DOM.Find(".description").Find("p").Text())
	// 	}
	// })

	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println(e.Text + "\n")
	})

	c.OnHTML("article ", func(e *colly.HTMLElement) {
		//fmt.Println(e.Text + "\n")
		//fmt.Println("Description: ", e.DOM.Find(".post-title").Find("h1").Text())
		//md := html2md.Convert(e.Text)
		converter := md.NewConverter("", true, nil)

		//html = `<strong>Important</strong>`

		md, err := converter.ConvertString(e.Text)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(md)
		go writeIntoMD(md, fileN)
	})

	// Callback for links on scraped pages
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Extract the linked URL from the anchor tag
		link := e.Attr("href")
		// Have our crawler visit the linked URL
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		fileN = lastString(strings.Split(r.URL.String(), "/"))
		//fmt.Println("---------------", lastString(strings.Split(r.URL.String(), "/")))
	})

	c.Visit("https://golangbot.com/learn-golang-series/")
}
