package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
)

func main() {
	url := "https://check.torproject.org/"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".content h1").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		fmt.Println(name)
	})
}
