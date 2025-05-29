package webpconverter

import (
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/webp"
)

func main() {
	m := flag.String("m", "jpeg", "MIME type to convert to")
	q := flag.Int("q", 100, "Quality of image")
	flag.Parse()
	options := jpeg.Options{Quality: *q} // For jpeg encoding

	// Get directory user's system
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	// subdirectory will be the last argument when user runs program on command line
	subdir := os.Args[len(os.Args)-1]
	fmt.Printf("subdir path: %s\n", subdir)

	fullPath := filepath.Join(home, subdir)

	// Get all webp files in the directory
	files, err := filepath.Glob(fullPath + "/*.webp")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	for i, file := range files {

		// Get the current file's name, remove the .webp suffix, and replace with .png
		name := strings.TrimSuffix(filepath.Base(file), ".webp")
		img, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		// Decode the webp data, then create a new file
		webpImg, err := webp.Decode(img)

		var newFile *os.File
		if *m == "jpeg" {
			log.Printf("jpeg")
			newFile, err = os.Create(fullPath + "/" + name + strconv.Itoa(i) + ".jpg")
			err = jpeg.Encode(newFile, webpImg, &options)
			if err != nil {
				fmt.Println(err)
				log.Fatal(err)
			}
		} else {
			log.Printf("png")
			newFile, err = os.Create(fullPath + "/" + name + strconv.Itoa(i) + ".png")
			if err != nil {
				fmt.Println(err)
				log.Fatal(err)
			}

			// encode the data into png and write it to the new file
			err = png.Encode(newFile, webpImg)
			if err != nil {
				fmt.Println(err)

				log.Fatal(err)
			}

		}
		img.Close()
	}
}
