package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Image struct {
	Id          int64
	Source      string `xorm: "unique"`
	Path        string
	Name        string `xorm: "index"`
	Uploaded_to string
	Created_at  time.Time `xorm: "created"`
}

var (
	x     *xorm.Engine
	imgbi = "https://img.bi"
)

func InsertImage(image *Image) error {
	var hasImage = Image{Source: image.Source}
	has, err := x.Get(&hasImage)
	if !has {
		_, err = x.Insert(image)
	}
	return err
}

func ScrapeRgHost(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".main-column ul li a").Each(func(i int, s *goquery.Selection) {
		ScrapeFileLink(s)
	})
}

func ScrapeFileLink(s *goquery.Selection) {
	name := s.Text()
	href, _ := s.Attr("href")
	if IsImageType(name) {
		url := GetDirectLink(href)
		downloaded, fileName := DownloadImage(url, name)
		if downloaded {
			err := InsertImage(&Image{Source: url, Path: fileName, Name: name})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func IsImageType(name string) (result bool) {
	split_name := strings.Split(name, ".")
	ext := split_name[len(split_name)-1]
	if ext == "jpg" || ext == "png" || ext == "gif" || ext == "jpeg" {
		return true
	} else {
		return false
	}
}

func GetDirectLink(name string) (dl string) {
	return fmt.Sprintf("http://rghost.ru%s/image.png", name)
}

func DownloadImage(url string, name string) (downloaded bool, fileName string) {
	tokens := strings.Split(url, "/")
	fileName = fmt.Sprintf("images/%s-%s", tokens[len(tokens)-2], name)

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return false, fileName
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return false, fileName
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return false, fileName
	}

	return true, fileName
}

func main() {
	url := "http://rghost.net/main"
	var err error
	x, err = xorm.NewEngine("sqlite3", "./images.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer x.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	x.Sync(new(Image))
	for {
		ScrapeRgHost(url)
		time.Sleep(10 * time.Second)
	}

}
