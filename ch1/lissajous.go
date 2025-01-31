package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)

//var palette = []color.Color{color.White, color.Black}
var colorPalette = []color.Color{
	color.RGBA{255, 145, 199, 1}, 
	color.RGBA{0, 145, 36, 2}, 
	color.RGBA{100, 2, 3, 4}, 
	color.RGBA{122, 23, 12, 199},
}

const (
	background = 0 // or first index
	foreground = 1 // second index
	color3 = 2
	color4 = 3
)

func main() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles = 155 // number of complete x oscillator revolutions
		res = 0.001 // angular resolution
		size = 100 // image canvas covers [-size..+size]
		nframes = 64 // number of animation frames
		delay = 8
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1,2*size+1)
		img := image.NewPaletted(rect, colorPalette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			if cycles % 2 > 1 {
				img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), color4)
			} else {
				img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), color3)
			}
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
