package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type result struct {
	srcImage       string
	thumbnailImage *image.NRGBA
	err            error
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Need to send directory path of images")
	}

	start := time.Now()

	err := setupPipeline(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Time taken is : %s\n", time.Since(start))
}

func setupPipeline(root string) error {
	done := make(chan struct{})
	defer close(done)

	// First stage of pipeline
	paths, errc := walkFiles(done, root)

	// second stage
	results := processImage(done, paths)

	//Third stage
	for r := range results{
		if r.err != nil{
			return r.err
		}
		saveThumbnail(r.srcImage, r.thumbnailImage)
	}
	if err := <-errc; err != nil{
		return err
	}
	return nil
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {

	paths := make(chan string)
	errch := make(chan error, 1)

	go func() {
		defer close(paths)
		defer close(errch)
		errch <- filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}

			contentType, _ := getFileContentType(path)
			if contentType != "image/jpeg" {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return fmt.Errorf("walk was cancelled")
			}
			paths <- path
			return nil
		})
	}()

	return paths, errch
}

func getFileContentType(path string) (string, error) {
	out, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	//Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func processImage(done <-chan struct{}, paths <-chan string) <-chan *result {
	results := make(chan *result)

	thumbnailer := func() {
		for path := range paths {
			srcImage, err := imaging.Open(path)
			if err != nil {
				select {
				case results <- &result{path, nil, err}:
				case <-done:
					fmt.Errorf("")

				}
			}

			thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)

			select {
			case results <- &result{path, thumbnailImage, err}:
			case <-done:
				return

			}

		}
	}

	const numThumbNailer = 5

	var wg sync.WaitGroup

	wg.Add(numThumbNailer)
	for i := 0; i < numThumbNailer; i++ {
		go func() {
			thumbnailer()
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

func saveThumbnail(srcImagePath string, thumbnailImage *image.NRGBA) error {
	filename := filepath.Base(srcImagePath)
	dstImagePath := "thumbnail/" + filename

	err := imaging.Save(thumbnailImage, dstImagePath)
	return err
}
