package txtp

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"regexp"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Context struct {
	width, height int
	face          font.Face
	color         color.Color
	image         *image.NRGBA
	mask          *image.Alpha
	lineHeight    float64
	strokeWidth   int
}

func (ctx *Context) SetFontFace(face font.Face) {
	ctx.face = face
}

func (ctx *Context) SetColor(clr color.Color) {
	ctx.color = clr
}

func (ctx *Context) SetLineWidth(d int) {
	ctx.strokeWidth = d
}

func (ctx *Context) SetMask(m *image.Alpha) {
	ctx.mask = m
}

func (ctx *Context) SetHexColor(s string) {
	r, g, b, a := parseHexColor(s)
	ctx.color = color.NRGBA{r, g, b, a}
}

// ---------------- Strings ----------------

func (ctx *Context) DrawString(s string, x, y float64) {
	_, h := ctx.MeasureString(s)
	tmp := image.NewRGBA(ctx.image.Bounds())
	dr := &font.Drawer{
		Src:  image.NewUniform(ctx.color),
		Dst:  tmp,
		Face: ctx.face,
		Dot:  fixed.P(int(x), int(y+h)),
	}

	dr.DrawString(s)
	ctx.DrawImage(tmp, 0, 0)
}

func (ctx *Context) DrawStringAnchored(s string, x, y float64, ax, ay float64) {
	w, h := ctx.MeasureString(s)
	tmp := image.NewRGBA(ctx.image.Bounds())
	dr := &font.Drawer{
		Src:  image.NewUniform(ctx.color),
		Dst:  tmp,
		Face: ctx.face,
		Dot: fixed.Point26_6{
			X: fixed.Int26_6((x - ax*w) * 64),
			Y: fixed.Int26_6((y + h - ay*h) * 64),
		},
	}
	dr.DrawString(s)
	ctx.DrawImage(tmp, 0, 0)
}

func (ctx *Context) DrawStringWrapped(s string, x, y, maxWidth float64) {
	sp := regexp.MustCompile(`\n+`)
	for _, paragraph := range sp.Split(s, -1) {
		lines := ctx.textWrap(paragraph, maxWidth)
		for _, s := range lines {
			_, h := ctx.MeasureString(s)
			ctx.DrawString(s, x, y)
			y += h * ctx.lineHeight
		}

		y += float64(ctx.face.Metrics().Height.Floor())
	}
}

func (ctx *Context) MeasureString(s string) (w, h float64) {
	b, _ := font.BoundString(ctx.face, s)
	return float64((b.Max.X - b.Min.X) >> 6), float64((b.Max.Y - b.Min.Y) >> 6)
}

func (ctx *Context) MeasureMultilineString(s string, maxWidth float64) (w, h float64) {
	sp := regexp.MustCompile(`\n+`)
	for _, paragraph := range sp.Split(s, -1) {
		lines := ctx.textWrap(paragraph, maxWidth)
		for _, s := range lines {
			lineWidth, lineHeight := ctx.MeasureString(s)
			h += lineHeight
			if lineWidth > w {
				w = lineWidth
			}
		}

		h += float64(ctx.face.Metrics().Height.Floor())
	}

	h -= float64(ctx.face.Metrics().Height.Floor())
	return
}

// ---------------- IMAGE ----------------

func (ctx *Context) DrawImage(img image.Image, x, y float64) {
	dr := img.Bounds().Add(image.Pt(int(x), int(y)))
	if ctx.mask != nil {
		draw.DrawMask(ctx.image, dr, img, image.Point{}, ctx.mask, image.Point{}, draw.Over)
		return
	}

	draw.Draw(ctx.image, dr, img, image.Point{}, draw.Over)
}

func (ctx *Context) DrawSubImage(img image.Image, dr image.Rectangle, p image.Point) {
	if ctx.mask != nil {
		draw.DrawMask(ctx.image, dr, img, p, ctx.mask, image.Point{}, draw.Over)
		return
	}

	draw.Draw(ctx.image, dr, img, p, draw.Over)
}

// ---------------- MASK ---------------

func (ctx *Context) ClearMask() {
	ctx.mask = nil
}

func (ctx *Context) InvertMask() {
	if ctx.mask == nil {
		ctx.mask = image.NewAlpha(ctx.image.Bounds())
		return
	}

	for i, a := range ctx.mask.Pix {
		ctx.mask.Pix[i] = 255 - a
	}
}

func (ctx *Context) AsMask() *image.Alpha {
	a := image.NewAlpha(ctx.image.Bounds())
	draw.Draw(a, ctx.image.Bounds(), ctx.image, image.Point{}, draw.Src)
	return a
}

// ---------------- IO ----------------

func (ctx *Context) SavePNG(filename string) error {
	imgFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	return png.Encode(imgFile, ctx.image)
}

func (ctx *Context) Image() image.Image {
	return ctx.image
}

func New(w, h int) *Context {
	return FromImage(image.NewRGBA(image.Rect(0, 0, w, h)))
}

func FromImage(img image.Image) *Context {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	tmp := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.Draw(tmp, tmp.Rect, img, image.Pt(0, 0), draw.Src)

	return &Context{
		image:       tmp,
		face:        nil,
		width:       w,
		height:      h,
		color:       color.Black,
		lineHeight:  1.5,
		strokeWidth: 1,
	}
}

// ---------------- Rectangle ----------------
func (ctx *Context) FillRectangle(x, y, w, h int) {
	draw.Draw(ctx.image, image.Rect(x, y, x+w, y+h), image.NewUniform(ctx.color), image.Pt(0, 0), draw.Over)
	if ctx.mask != nil {
		draw.DrawMask(ctx.image, image.Rect(x, y, x+w, y+h), image.NewUniform(ctx.color), image.Pt(x, y), ctx.mask, image.Pt(0, 0), draw.Over)
		return
	}
}

func (ctx *Context) ClearRectangle(x, y, w, h int) {
	draw.Draw(ctx.image, image.Rect(x, y, x+w, y+h), image.NewUniform(color.NRGBA{}), image.Pt(0, 0), draw.Src)
}

func (ctx *Context) StrokeRectangle(x, y, w, h int) {
	tmp := New(ctx.width, ctx.height)
	tmp.SetColor(ctx.color)
	tmp.FillRectangle(x, y, w, h)
	tmp.ClearRectangle(x+ctx.strokeWidth, y+ctx.strokeWidth, w-ctx.strokeWidth*2, h-ctx.strokeWidth*2)

	ctx.DrawImage(tmp.Image(), 0, 0)
}
