package main

import (
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

// provided by ldflags
var version string
var build string

// spacing in mm
var horizontalSpacing = .4
var verticalSpacing = .4
var leftMargin = 10
var topMargin = 10

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please provide a folder with images as argument.")
	}
	dir := os.Args[1]
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var images = make([]string, 0, len(files))
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext == ".jpg" || ext == ".png" {
			images = append(images, f.Name())
		}
	}
	length := len(images)
	pages := int(math.Ceil(float64(length) / 9))
	log.Println("Creating A4 PDF from", length, "images. with number of pages", pages)

	pdf := gofpdf.New("P", "mm", "A4", "")
	var opt gofpdf.ImageOptions
	opt.ReadDpi = true
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(0, 10, "Generated with https://github.com/yene/ANRProxyGenerator", "", 0, "C", false, 0, "")
	})

	for i := 0; i < pages; i++ {
		pdf.AddPage()
		for j, card := range chunkSlice(images, i*9, (i+1)*9) {
			p := filepath.Join(dir, card)
			// Netrunner cards are exactly 88x61mm
			poscol := j % 3
			x := float64(leftMargin + (poscol * 61))
			x = x + (float64(poscol) * horizontalSpacing)

			posrow := j / 3
			y := float64(topMargin + (posrow * 88))
			y = y + (float64(posrow) * verticalSpacing)
			pdf.ImageOptions(p, float64(x), float64(y), 61, 88, false, opt, 0, "")
		}
	}

	outname := filepath.Base(dir)
	err = pdf.OutputFileAndClose(outname + ".pdf")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Written PDF to:", outname+".pdf")
}

func chunkSlice(slice []string, start int, end int) []string {
	if end > len(slice) {
		return slice[start:]
	}
	return slice[start:end]
}
