package txtp

import (
	"image"
	"image/png"
	"log"
	"os"
	"testing"
)

func TestDrawRectangle(t *testing.T) {
	msk := New(256, 256)
	msk.FillRectangle(0, 0, 256, 256)
	msk.ClearRectangle(64, 64, 128, 128)
	msk.SavePNG("example-msk.png")

	blk := New(256, 256)
	blk.FillRectangle(128, 128, 128, 128)
	blk.SavePNG("example-blk.png")

	ctx := New(256, 256)
	ctx.SetLineWidth(2)
	ctx.SetMask(msk.AsMask())

	// black
	ctx.DrawSubImage(blk.Image(), image.Rect(0, 0, 128, 128), image.Pt(128, 128))

	// red
	ctx.SetColor(Red)
	ctx.StrokeRectangle(128, 128, 128, 128)

	ctx.SavePNG("example.png")

	imgFile, err := os.Open("rectangle_test.png")
	if err != nil {
		log.Fatalln(err)
	}

	img, err := png.Decode(imgFile)
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < img.Bounds().Dx(); i++ {
		for r := 0; r < img.Bounds().Dy(); r++ {
			if img.(*image.Paletted).At(i, r) != ctx.image.At(i, r) {
				t.Errorf("image not match at (%d, %d)\n", i, r)
				ctx.SavePNG("report-rectangle_test.png")
				return
			}
		}
	}
}
