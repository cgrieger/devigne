package main

import (
	"os"
	"math"
	"log"
	"image"
	"image/color"
	r "github.com/cgrieger/devigne/resize"
	_ "image/jpeg"
	png "image/png"
)

func main() {
	grayImage := ReadGrayscale(os.Args[1])
	grayImage = r.Resize(grayImage, grayImage.Bounds(), 256, 256)
	grayImage = kangWeiss(grayImage)
	SaveGrayscale(grayImage, os.Args[2])
}

func kangWeiss(src image.Image) image.Image {
	// focal length
	f := float64(30)
	// geometric vignetting factor
	a1 := float64(1)

	angleX := float64(0);
	angleT := float64(0.2);

	result := eachPoint(src, func (x,y int, v uint8) uint8 {
		r := distanceToCenter(src, x ,y)
		A := 1 / sq(1 + sq(r / f))
		G := (1 - (a1 * r))
		vR := A * G

		T := math.Cos(a1) * 
			cube(1 + ((math.Tan(angleT) / f) * 
			(float64(x) * math.Sin(angleX)) - (float64(y) * math.Cos(angleX))))

		return uint8(vR * T)
	})
	
	return result
}

func eachPoint(src image.Image, fn func(x, y int, v uint8) uint8) image.Image {
	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := image.NewGray(bounds)

	for x := 1; x < w - 1; x++ {
		for y := 1; y < h - 1; y++ {
			v, _, _ , _ := src.At(x,y).RGBA()
			colorValue := uint8(v)
			color := &color.Gray{uint8(fn(x,y,colorValue))}
			result.Set(x,y, color)
		}
	}

	return result
}

func sq(v float64) float64{
	return v * v
}


func cube(v float64) float64{
	return v * v * v
}

func distanceToCenter(src image.Image, x, y int) float64 {
	rect := src.Bounds()
	w, h := rect.Dx(), rect.Dy()

	dist := float64(((x - w) * (x - w)) + ((y - h) * (y - h)))
	return math.Sqrt(dist)
}

func GradientImage(src image.Image) image.Image {
	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	grad := image.NewGray(bounds)


	for x := 1; x < w - 1; x++ {
		for y := 1; y < h - 1; y++ {

			vXPrev, _, _ , _ := src.At(x-1,y).RGBA()
			vXNext, _, _ , _ := src.At(x+1,y).RGBA()

			vYPrev, _, _ , _ := src.At(x,y-1).RGBA()
			vYNext, _, _ , _ := src.At(x,y+1).RGBA()

			gX := int((vXNext - vXPrev) / 2)
			gY := int((vYNext - vYPrev) / 2)

			abs := float64((gX * gX) + (gY * gY))
			v := &color.Gray{uint8(math.Sqrt(abs))}
			grad.Set(x,y,v)
		}
	}

	return grad
}

func SaveGrayscale(image image.Image, path string) {
	 // Encode the grayscale image to the output file
    outfile, err := os.Create(path)
    if err != nil {
        // replace this with real error handling
        panic(err)
    }
    defer outfile.Close()
    png.Encode(outfile, image)
}

func ReadGrayscale(path string) image.Image {
	log.Print("Reading image from ", path)
	infile, err := os.Open(path)

	if err != nil {
		panic(err)
	}
	defer infile.Close()

	src, _, err := image.Decode(infile)
    if err != nil {
        // replace this with real error handling
        panic(err)
    }

    // Create a new grayscale image
    bounds := src.Bounds()
    w, h := bounds.Dx(), bounds.Dy()

    gray := image.NewGray(bounds)
    for x := 0; x < w; x++ {
        for y := 0; y < h; y++ {
            oldColor := src.At(x, y)
            converter := image.NewUniform(oldColor)
            grayColor := converter.Convert(oldColor)
            gray.Set(x, y, grayColor)
        }
    }

    return gray
}