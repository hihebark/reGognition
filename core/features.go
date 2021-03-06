package core

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hihebark/gore/log"
)

//Blur image.
func Blur(imgsrc image.Image, radius float64) image.Image {
	maxY, maxX := imgsrc.Bounds().Max.Y, imgsrc.Bounds().Max.X
	imgdst := image.NewRGBA(imgsrc.Bounds())
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			var r, g, b, a uint32 = 0, 0, 0, 0
			var count uint32

			for ky := -radius; ky < radius; ky++ {
				for kx := -radius; kx <= radius; kx++ {
					kr, kg, kb, ka := imgsrc.At(x+int(kx), y+int(ky)).RGBA()
					r += kr
					g += kg
					b += kb
					a += ka
					count++
				}
			}
			c := color.RGBA{uint8(r/count) + 1, uint8(g/count) + 1, uint8(b/count) + 1, uint8(a / count)}
			imgdst.Set(x, y, c)
		}
	}
	return imgdst
}

//GaussianBlur blur image with gauss formula.
func GaussianBlur(imgsrc image.Image, kernel, radius int) image.Image {
	bounds := imgsrc.Bounds()
	maxX, maxY := bounds.Max.X, bounds.Max.Y
	imgdst := image.NewRGBA64(bounds)
	l := maxY * maxX
	kernels := gaussianMap(kernel, float64(radius))
	log.Inf("+ There is %d cells", l)
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			var r, g, b, a uint16
			k := -kernel
			if x == 0 || y == 0 {
				k = 0
			}
			for ky := -k; ky < kernel; ky++ {
				for kx := -k; kx < kernel; kx++ {
					kr, kg, kb, ka := imgsrc.At(x+kx, y+ky).RGBA()
					r += uint16(float64(kr) * kernels[kernel+kx][kernel+ky])
					g += uint16(float64(kg) * kernels[kernel+kx][kernel+ky])
					b += uint16(float64(kb) * kernels[kernel+kx][kernel+ky])
					a += uint16(float64(ka) * kernels[kernel+kx][kernel+ky])
				}
			}
			imgdst.SetRGBA64(x, y, color.RGBA64{r, g, b, a})
			fmt.Printf("- Processing with %5d cell\r", (maxX*maxY)-l)
			l--
		}
	}
	fmt.Printf("\n")
	return imgdst
}

func gaussianMap(ks int, sigma float64) [][]float64 {
	var sum float64
	l := ks*2 + 1
	kernel := make([][]float64, l)
	for i := 0; i < l; i++ {
		row := make([]float64, l)
		for j := 0; j < l; j++ {
			g := Gaussian(i, j, sigma)
			row[j] = g
			sum += g
		}
		kernel[i] = row
	}
	for i := 0; i < l; i++ {
		for j := 0; j < l; j++ {
			kernel[i][j] /= sum
		}
	}
	return kernel
}
func GaborFilter(imgsrc image.Image, bounds image.Rectangle) {
	//	maxX := bounds.Max.X
	//	maxY := bounds.Max.Y
}
func GaborFilterKernel(period, phase, angle float64, size int, imgsrc image.Image) [][]float64 {
	major := period / 3.0
	theta := angle + 90.0
	if size == -1 {
		size = int(math.Ceil(major * math.Sqrt(-2.0*math.Log(math.Exp(-5.0)))))
	} else {
		size /= 2
	}
	l := size*2 + 1
	kernels := make([][]float64, l)
	psi := PI / 180.0 * phase
	rtDeg := PI / 180.0 * theta
	omega := (2.0 * PI) / period
	co := math.Cos(rtDeg)
	si := math.Sin(rtDeg)
	j := 0
	majorsigq := 2.0 * major * major
	minorsigq := majorsigq
	for y := -size; y <= size; y++ {
		row := make([]float64, l)
		i := 0
		for x := -size; x <= size; x++ {
			maj, min := float64(x)*co+float64(y)*si, float64(x)*si-float64(y)*co
			row[i] = math.Cos(omega*maj+psi) * math.Exp(-(maj*maj)/majorsigq) * math.Exp(-(min*min)/minorsigq)
			i++
		}
		kernels[j] = row
		j++
	}
	return kernels
}
func IntensityFeatures(imgsrc image.Image) image.Image {
	log.Inf("Extracting intensity features from image")
	maxX, maxY := imgsrc.Bounds().Max.X, imgsrc.Bounds().Max.Y
	imgdst := image.NewGray(imgsrc.Bounds())
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			r, g, b, _ := imgsrc.At(x, y).RGBA()
			imgdst.SetGray(x, y, color.Gray{uint8((r + g + b) / 3)})
		}
	}
	return imgdst
}
