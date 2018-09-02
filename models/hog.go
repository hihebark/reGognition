package model

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hihebark/gore/core"
	"golang.org/x/image/draw"
)

//HogVect hog implementation.
func HogVect(imgsrc image.Image, i *core.ImageInfo) image.Image {
	bound := imgsrc.Bounds()
	hogimg := image.NewRGBA(bound)
	draw.Draw(hogimg, bound, &image.Uniform{color.Black}, image.ZP, draw.Src)
	cells := core.Divid(bound, i.Cellsize)
	midcell := image.Pt(int(i.Cellsize/2)+1, int(i.Cellsize/2)+1)
	c := color.RGBA{0xff, 0xff, 0xff, 0xee}
	fmt.Printf("+ There is %d cells\n", len(cells)-1)
	for k, cell := range cells {
		i.Wg.Add(1)
		fmt.Printf("- Processing with %d cell\r", k)
		imgcell := image.NewRGBA(cell)
		for y := cell.Min.Y; y < cell.Max.Y; y++ {
			for x := cell.Min.X; x < cell.Max.X; x++ {
				yd := math.Abs(float64(imgsrc.At(x, y-1).(color.Gray).Y - imgsrc.At(x, y+1).(color.Gray).Y))
				xd := math.Abs(float64(imgsrc.At(x-1, y).(color.Gray).Y - imgsrc.At(x+1, y).(color.Gray).Y))
				magnitude, orientation := core.GradientVector(xd, yd)
				imgcell = core.DrawLine(cell.Sub(midcell).Max, orientation, magnitude, imgcell, c)
			}

		}
		draw.Draw(hogimg, imgcell.Bounds(), imgcell, cell.Min, draw.Over)
		i.Wg.Done()
	}
	fmt.Print("\n")
	return hogimg
}
