// Command ogbase generates the 1200x630 background used for social cards.
// Colours mirror the dark theme in assets/css/tokens.css.
package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

const (
	width  = 1200
	height = 630
	inset  = 40
	rule   = 3
)

var (
	bg     = color.RGBA{0x0d, 0x11, 0x17, 0xff} // --bg
	border = color.RGBA{0x30, 0x36, 0x3d, 0xff} // --border
	accent = color.RGBA{0xe0, 0x85, 0x3f, 0xff} // --accent
)

func main() {
	out := flag.String("out", "assets/og/base.png", "output png path")
	flag.Parse()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// hairline frame
	frame := image.Rect(inset, inset, width-inset, height-inset)
	drawRect(img, frame, border, 1)

	// accent rule under where the title will sit
	draw.Draw(img,
		image.Rect(inset+40, height-inset-90, inset+40+120, height-inset-90+rule),
		&image.Uniform{accent}, image.Point{}, draw.Src)

	f, err := os.Create(*out)
	if err != nil {
		log.Fatalf("create %s: %v", *out, err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatalf("encode: %v", err)
	}
	log.Printf("wrote %s (%dx%d)", *out, width, height)
}

// drawRect strokes the outline of r with the given colour and thickness.
func drawRect(dst draw.Image, r image.Rectangle, c color.Color, t int) {
	u := &image.Uniform{c}
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+t), u, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y-t, r.Max.X, r.Max.Y), u, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Min.X+t, r.Max.Y), u, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(r.Max.X-t, r.Min.Y, r.Max.X, r.Max.Y), u, image.Point{}, draw.Src)
}
