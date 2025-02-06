package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
)

var colorPalette = []color.Color{
	// Soft Pastel Shades
	color.Black,
	color.White,
	color.RGBA{255, 215, 255, 1}, // Light lavender
	color.RGBA{255, 182, 222, 1}, // Soft pink
	color.RGBA{200, 165, 255, 1}, // Light periwinkle
	color.RGBA{250, 180, 225, 1}, // Soft peachy pink
	color.RGBA{210, 170, 255, 1}, // Lilac

	// Bold & Vibrant Tones
	color.RGBA{255, 56, 56, 1},  // Strong red
	color.RGBA{25, 205, 25, 1},  // Fresh green
	color.RGBA{46, 204, 113, 1}, // Bright medium green
	color.RGBA{255, 160, 40, 1}, // Tangerine orange
	color.RGBA{128, 0, 255, 1},  // Deep purple

	// Rich & Deep Colors
	color.RGBA{40, 58, 80, 1},  // Slate blue-gray
	color.RGBA{24, 120, 50, 1}, // Forest green
	color.RGBA{18, 167, 52, 1}, // Emerald green
	color.RGBA{51, 120, 92, 1}, // Teal
	color.RGBA{27, 140, 70, 1}, // Dark sea green

	// Neon & Glowing
	color.RGBA{255, 0, 0, 1},   // Neon red
	color.RGBA{0, 255, 0, 1},   // Neon green
	color.RGBA{0, 151, 255, 1}, // Neon blue
	color.RGBA{255, 165, 0, 1}, // Neon orange
	color.RGBA{168, 0, 255, 1}, // Neon purple

	// Shimmery & Metallic Hues
	color.RGBA{232, 180, 120, 1}, // Champagne gold
	color.RGBA{211, 105, 55, 1},  // Coppery bronze
	color.RGBA{238, 173, 245, 1}, // Soft metallic lavender
	color.RGBA{180, 120, 100, 1}, // Antique gold
	color.RGBA{248, 220, 155, 1}, // Warm beige metallic

	// Earthy & Rustic Tones
	color.RGBA{116, 77, 53, 1}, // Chocolate brown
	color.RGBA{68, 38, 24, 1},  // Deep coffee brown
	color.RGBA{95, 47, 28, 1},  // Burnt sienna
	color.RGBA{107, 57, 34, 1}, // Deep terracotta
	color.RGBA{148, 88, 51, 1}, // Rustic amber

	// Moody & Dark Shades
	color.RGBA{15, 15, 15, 1},    // Almost black
	color.RGBA{30, 50, 70, 1},    // Dark charcoal
	color.RGBA{60, 85, 110, 1},   // Dark slate blue
	color.RGBA{90, 115, 130, 1},  // Steel blue
	color.RGBA{120, 140, 160, 1}, // Dark steel gray
}

var (
	background, _ = strconv.Atoi(os.Args[1:][0])
	foreground, _ = strconv.Atoi(os.Args[1:][1])
)

func main() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 255   // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 220   // image canvas covers [-size..+size]
		nframes = 128   // number of animation frames
		delay   = 8
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, colorPalette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			if cycles%2 > 1 {
				img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(background))
			} else {
				img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(foreground))
			}
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
