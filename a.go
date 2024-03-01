package main

import (
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/images", imageHandler)
	http.ListenAndServe(":8090", nil)
}

func imageHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	imgUrl := req.URL.Query().Get("url")
	fmt.Printf("imgUrl: %s\n", imgUrl)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n\n", err)
	}

	fmt.Printf("%s: got / request. body:\n%s\n\n", ctx.Value("serverAddr"), body)

	fileName := "tmp-file.jpg"
	errDownloading := downloadFile(imgUrl, fileName)
	if errDownloading != nil {
		log.Fatal(errDownloading)
	}

	resizeImage(fileName)

	fmt.Fprintf(w, "imageHandler\n")
}

func resizeImage(imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(1000, 0, img, resize.Lanczos3)

	out, err := os.Create("test_resized.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}

func downloadFile(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
