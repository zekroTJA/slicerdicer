package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"path"
	"strings"

	"github.com/nfnt/resize"
	"github.com/zekroTJA/slicerdicer/pkg/slicerdicer"
)

var (
	flagImageFile     = flag.String("i", "", "imput image file")
	flagSlicesPerSide = flag.Int("s", 2, "ammount of slices per side")
	flagOutLocation   = flag.String("o", "results", "output files location")
	flagOutTyp        = flag.String("otyp", "png", "output image file type")
	flagOutName       = flag.String("oname", "slice", "output name prefix")
	flagScale         = flag.Float64("scale", 1, "The scale of the result images as multiplier.")
)

func checkErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func getFileExt(fname string) (string, error) {
	iExt := strings.LastIndex(*flagImageFile, ".")
	if iExt < 0 || iExt == len(fname)-1 {
		return "", errors.New("file name has no extension")
	}

	return fname[iExt+1:], nil
}

func getDecoderByExt(ext string) (decoder func(io.Reader) (image.Image, error), err error) {
	switch ext {

	case "png":
		decoder = png.Decode
	case "jpg", "jpeg":
		decoder = jpeg.Decode

	default:
		err = errors.New("unsupported file type")
	}

	return
}

func getEncoderByExt(ext string) (encoder func(io.Writer, image.Image) error, err error) {
	switch ext {

	case "png":
		encoder = png.Encode
	case "jpg", "jpeg":
		encoder = func(w io.Writer, i image.Image) error {
			return jpeg.Encode(w, i, &jpeg.Options{Quality: 85})
		}

	default:
		err = errors.New("unsupported file type")
	}

	return
}

func writeImage(encoder func(io.Writer, image.Image) error, i image.Image, name string) error {
	loc := path.Join(*flagOutLocation, name)

	f, err := os.Create(loc)
	if err != nil {
		return err
	}
	defer f.Close()

	return encoder(f, i)
}

func resizeImage(img image.Image, m float64) image.Image {
	if m == 1 {
		return img
	}

	bounds := img.Bounds()
	w := math.Floor(float64(bounds.Dx()) * m)
	h := math.Floor(float64(bounds.Dy()) * m)

	return resize.Resize(uint(w), uint(h), img, resize.Bicubic)
}

func main() {
	flag.Parse()

	if *flagImageFile == "" {
		log.Fatal("Input file must be specified")
	}

	if *flagScale <= 0 {
		log.Fatal("Scale must be larger than 0")
	}

	f, err := os.Open(*flagImageFile)
	checkErr(err)

	ext, err := getFileExt(*flagImageFile)
	checkErr(err)

	decoder, err := getDecoderByExt(ext)
	checkErr(err)

	img, err := decoder(f)
	checkErr(err)

	img = resizeImage(img, *flagScale)

	res, err := slicerdicer.Slice(img, *flagSlicesPerSide)
	checkErr(err)

	encoder, err := getEncoderByExt(*flagOutTyp)
	checkErr(err)

	dInfo, err := os.Stat(*flagOutLocation)
	if os.IsNotExist(err) {
		checkErr(
			os.MkdirAll(*flagOutLocation, os.ModeDir))
	} else if err == nil {
		if !dInfo.IsDir() {
			log.Fatal("given output location is a file")
		}
	} else {
		checkErr(err)
	}

	for y, row := range res {
		for x, slice := range row {
			name := fmt.Sprintf("%s_%d_%d.%s", *flagOutName, x, y, *flagOutTyp)
			log.Printf("writing image %s", name)
			checkErr(
				writeImage(encoder, slice, name))
		}
	}

	log.Println("finished")
}
