package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-xorm/xorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Images struct {
	Id          int64
	Source      string `xorm: "unique"`
	Path        string
	Name        string `xorm: "index"`
	Uploaded_to string
	Created_at  time.Time `xorm: "created"`
	Updated_at  time.Time `xorm: "created"`
}

var (
	x       *xorm.Engine
	imgbi   string = "https://img.bi"
	sqlite  string = "./images.db"
	psql    string = "dbname=images_test user=lenny password=123456 sslmode=disable"
	img_dir string = "images/"
)

func InsertImage(image *Images) error {
	var hasImage = Images{Source: image.Source}
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
			err := InsertImage(&Images{Source: url, Path: fileName, Name: name})
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
	dirname := fmt.Sprintf("%s/%s", img_dir, time.Now().Format("20060102"))
	dirCreated, err := dirExists(dirname)
	if err != nil {
		fmt.Println("Error while creating", dirname, "-", err)
		return false, fileName
	}
	if !dirCreated {
		err := os.Mkdir(dirname, 0777)
		if err != nil {
			fmt.Println("Error while creating", dirname, "-", err)
			return false, dirname
		}
	}
	fileName = fmt.Sprintf("%s/%s-%s", dirname, tokens[len(tokens)-2], name)

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

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DB_URL := os.Getenv("DB_URL")
	IMG_DIR := os.Getenv("IMG_DIR")
	if DB_URL != "" {
		psql = DB_URL
	}
	if IMG_DIR != "" {
		img_dir = IMG_DIR
	}
}

func main() {
	url := "http://rghost.ru/main"
	var err error
	x, err = xorm.NewEngine("postgres", psql)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer x.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	x.Sync(new(Images))
	for {
		ScrapeRgHost(url)
		time.Sleep(10 * time.Second)
	}

}
