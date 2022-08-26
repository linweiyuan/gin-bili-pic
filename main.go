package main

import (
	"embed"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"html/template"

	//"image/jpeg"
	"io"
	"net/http"
	"os"
	//"strconv"
	"strings"
)

const fileExt = ".jpg"

var downloadImageCh = make(chan bool)
var videoNameCh = make(chan string)
var fileDeleteCh = make(chan bool)

//go:embed templates
var fs embed.FS

// func downloadImage(bv, width, height string) {
func downloadImage(bv string) {
	resp, _ := http.Get("https://www.bilibili.com/video/" + bv)
	defer resp.Body.Close()

	document, _ := goquery.NewDocumentFromReader(resp.Body)
	if title, exists := document.Find(".video-title").First().Attr("title"); exists {
		videoNameCh <- title + fileExt
	} else {
		videoNameCh <- "video not found"
	}
	selection := document.Find(`meta[itemprop="image"]`)
	if url, exists := selection.First().Attr("content"); exists {
		imageResp, _ := http.Get(url)
		defer imageResp.Body.Close()

		//image, _ := jpeg.Decode(imageResp.Body)
		//originalWidth := image.Bounds().Size().X
		//originalHeight := image.Bounds().Size().Y
		//if width != "" {
		//	w, _ := strconv.Atoi(width)
		//	if w < originalWidth {
		//		originalWidth = w
		//	}
		//}
		//if height != "" {
		//	h, _ := strconv.Atoi(height)
		//	if h < originalHeight {
		//		originalHeight = h
		//	}
		//}
		//newUrl := fmt.Sprintf("%s@%dw_%dh.png", url, originalWidth, originalHeight)
		//newImageResp, _ := http.Get(newUrl)
		//defer newImageResp.Body.Close()

		file, _ := os.Create(bv + fileExt)
		defer file.Close()
		//io.Copy(file, newImageResp.Body)
		io.Copy(file, imageResp.Body)

		downloadImageCh <- true

		if ok := <-fileDeleteCh; ok {
			os.Remove(file.Name())
		}
	} else {
		downloadImageCh <- false
	}
}
func main() {
	router := gin.Default()
	//router.LoadHTMLGlob("templates/*")
	tmpl := template.Must(template.New("").ParseFS(fs, "templates/*.html"))
	router.SetHTMLTemplate(tmpl)
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.POST("/", func(context *gin.Context) {
		bv := context.PostForm("bv")
		//width := context.PostForm("width")
		//height := context.PostForm("height")
		if !strings.HasPrefix(bv, "BV") {
			return
		}

		//go downloadImage(bv, width, height)
		go downloadImage(bv)

		fileName := <-videoNameCh
		result := <-downloadImageCh
		if result {
			context.Header("Content-Type", "application/octet-stream")
			context.Header("Content-Disposition", "attachment; filename="+fileName)
			context.File(bv + fileExt)
		}
		fileDeleteCh <- true
	})
	router.Run(":8080")
}
