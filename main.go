package main

import (
	//"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-xorm/xorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	//"github.com/nfnt/resize"
	//"image"
	//"image/gif"
	//"image/jpeg"
	//"image/png"
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
	Name        string `xorm: "index not null"`
	Uploaded_to string `xorm: default 'none'`
	Created_at  time.Time
	Updated_at  time.Time
	File_type   string `xorm: "index"`
	Archived    bool
	Thumb_path  string
}

var (
	x       *xorm.Engine
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
	isImage, ext := IsImageType(name)
	if isImage {
		url := GetDirectLink(href)
		downloaded, fileName, thumbName := DownloadImage(url, name, ext)
		if downloaded {
			err := InsertImage(&Images{Source: url, Path: fileName, Thumb_path: thumbName, Name: name, Created_at: time.Now(), Updated_at: time.Now(), Uploaded_to: "yes", File_type: ext})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func IsImageType(name string) (result bool, ext_type string) {
	split_name := strings.Split(name, ".")
	ext := split_name[len(split_name)-1]
	if ext == "jpg" || ext == "png" || ext == "gif" || ext == "jpeg" {
		return true, ext
	} else {
		return false, ext
	}
}

func GetDirectLink(name string) (dl string) {
	return fmt.Sprintf("http://rghost.ru%s/image.png", name)
}

/*func createThumb(fileName string, dirname string, ext string, nowTime string) (thumb string, err error) {

	thumbDir := fmt.Sprintf("%s/thumb", dirname)
	thumbDirCreated, err := dirExists(thumbDir)
	if err != nil {
		fmt.Println("Error while looking for", thumbDirCreated, "-", err)
	}
	if !thumbDirCreated {
		err := os.Mkdir(thumbDir, 0777)
		if err != nil {
			fmt.Println("Error while creating", thumbDir, "-", err)
			return thumbDir, err
		}
	}
	filePath := fmt.Sprintf("%s/%s", dirname, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error while opening", filePath, "-", err)
		return fileName, err
	}
	var img image.Image
	if ext == "jpeg" || ext == "jpg" {
		img, err = jpeg.Decode(file)
	} else if ext == "png" {
		img, err = png.Decode(file)
	} else if ext == "gif" {
		img, err = gif.Decode(file)
	} else {
		err = errors.New("type is not supported")
	}
	if err != nil {
		fmt.Println("Error while decoding", fileName, "-", err)
		return fmt.Sprintf("images/%s/thumb/%s", nowTime, fileName), nil
	}
	file.Close()

	m := resize.Thumbnail(200, 150, img, resize.NearestNeighbor)
	thumbName := fmt.Sprintf("%s/%s", thumbDir, fileName)
	out, err := os.Create(thumbName)
	if err != nil {
		return fmt.Sprintf("images/%s/thumb/%s", nowTime, fileName), err
	}
	defer out.Close()

	if ext == "jpeg" || ext == "jpeg" {
		jpeg.Encode(out, m, nil)
	} else if ext == "png" {
		png.Encode(out, m)
	} else if ext == "gif" {
		gif.Encode(out, m, nil)
	}
	return fmt.Sprintf("images/%s/thumb/%s", nowTime, fileName), err
}
*/

func DownloadImage(url string, name string, ext string) (downloaded bool, fileName string, thumbName string) {
	tokens := strings.Split(url, "/")
	nowTime := time.Now().Format("20060102")
	dirname := fmt.Sprintf("%s/%s", img_dir, nowTime)
	dirCreated, err := dirExists(dirname)
	if err != nil {
		fmt.Println("Error while looking for", dirname, "-", err)
	}
	if !dirCreated {
		err := os.Mkdir(dirname, 0777)
		if err != nil {
			fmt.Println("Error while creating", dirname, "-", err)
			return false, "", ""
		}
	}

	fileName = fmt.Sprintf("%s/%s-%s", dirname, tokens[len(tokens)-2], name)

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return false, "", ""
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return false, "", ""
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return false, "", ""
	}

	//filePath := fmt.Sprintf("%s-%s", tokens[len(tokens)-2], name)
	//thumb, err := createThumb(filePath, dirname, ext, nowTime)
	//if err != nil {
	//	thumb = fmt.Sprintf("images/%s/%s-%s", nowTime, tokens[len(tokens)-2], name)
	//}
	thumb := fmt.Sprintf("images/%s/%s-%s", nowTime, tokens[len(tokens)-2], name)
	fileName = fmt.Sprintf("images/%s/%s-%s", nowTime, tokens[len(tokens)-2], name)
	return true, fileName, thumb

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
