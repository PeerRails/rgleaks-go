package rgleaksgo

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/daddye/vips"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"io"
	"io/ioutil"
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
	Name        string `xorm: "index not null"`
	Uploaded_to string `xorm: default 'none'`
	Created_at  time.Time
	Updated_at  time.Time
	File_type   string `xorm: "index"`
	Archived    bool   `xorm: "index default false"`
	Thumbnail   string
}

var (
	x       *xorm.Engine
	psql    string = "dbname=images_test user=lenny password=123456 sslmode=disable"
	img_dir string = "images"
)

func (image *Images) InsertImage() error {
	var hasImage = Images{Source: image.Source}
	has, err := x.Get(&hasImage)
	if !has {
		_, err = x.Insert(image)
	}
	return err
}

func torcheck() {
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
	i := Images{Created_at: time.Now(), Updated_at: time.Now()}
	i.Name = s.Text()
	href, _ := s.Attr("href")
	isImage := i.IsImageType()
	if isImage {
		i.Source = fmt.Sprintf("http://rghost.ru%s/image.png", href)
		downloaded := i.DownloadImage()
		if downloaded {
			i.Uploaded_to = "yes"
			i.Archived = false
			err := i.InsertImage()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (image *Images) IsImageType() (result bool) {
	split_name := strings.Split(image.Name, ".")
	image.File_type = split_name[len(split_name)-1]
	if image.File_type == "jpg" || image.File_type == "png" || image.File_type == "gif" || image.File_type == "jpeg" {
		return true
	} else {
		return false
	}
}

func (image *Images) DownloadImage() (downloaded bool) {
	tokens := strings.Split(image.Source, "/")
	nowTime := time.Now().Format("20060102")
	dirname := fmt.Sprintf("%s/%s", img_dir, nowTime)
	err := dirExists(dirname)
	if err != nil {
		fmt.Println("Error while looking for", dirname, "-", err)
		return false
	}

	thumb_dirname := fmt.Sprintf("%s/%s/thumb", img_dir, nowTime)
	err = dirExists(thumb_dirname)
	if err != nil {
		fmt.Println("Error while looking for", dirname, "-", err)
		return false
	}

	fileName := fmt.Sprintf("%s/%s.%s", dirname, tokens[len(tokens)-2], image.File_type)

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating original", fileName, "-", err)
		return false
	}
	defer output.Close()

	response, err := http.Get(image.Source)
	if err != nil {
		fmt.Println("Error while downloading original", image.Source, "-", err)
		return false
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while saving original", image.Source, "-", err)
		return false
	}

	image.Thumbnail = fmt.Sprintf("images/%s/thumb/%s.%s", nowTime, tokens[len(tokens)-2], image.File_type)
	image.Path = fmt.Sprintf("images/%s/%s.%s", nowTime, tokens[len(tokens)-2], image.File_type)
	err = image.CreateThumbnail(fmt.Sprintf("%s/%s.%s", thumb_dirname, tokens[len(tokens)-2], image.File_type))
	if err != nil {
		ferr := os.Remove(image.Path)
		if ferr != nil {
			fmt.Println("Error while removing", image.Path, "-", ferr)
		}
		return false
	}

	return true

}

func (image *Images) CreateThumbnail(thumb string) error {
	options := vips.Options{
		Width:        400,
		Height:       300,
		Crop:         false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      95,
	}

	f, _ := os.Open(image.Path)
	inBuf, _ := ioutil.ReadAll(f)
	defer f.Close()
	buf, err := vips.Resize(inBuf, options)
	if err != nil {
		fmt.Println("Error while opening original", image.Path, "-", err)
		return err
	}

	file, err := os.Create(thumb)
	if err != nil {
		fmt.Println("Error while creating thumbnail", thumb, "-", err)
		return err
	}
	defer file.Close()
	err = ioutil.WriteFile(thumb, buf, 777)
	if err != nil {
		fmt.Println("Error while saving thumbnail", thumb, "-", err)
		return err
	}
	return nil
}

func dirExists(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	DB_URL := os.Getenv("DB_URL")
	IMG_DIR := os.Getenv("IMG_DIR")

	if DB_URL != "" {
		psql = DB_URL
	}
	if IMG_DIR != "" {
		img_dir = IMG_DIR
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "torcheck" {
			torcheck()
		} else {
			fmt.Println("Known args: torcheck")
		}
		os.Exit(1)
	}
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
}
