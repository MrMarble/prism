package prism

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/fogleman/gg"
	"github.com/mrmarble/prism/themer"
	"github.com/mrmarble/prism/tokenizer"
	"golang.org/x/image/font"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Context struct {
	font   font.Face
	points float64
	theme  themer.Theme
	lang   tokenizer.Language
}

func NewContext() *Context {
	return &Context{}
}

func (ctx *Context) SetFontFace(path string, points float64) error {
	font, err := gg.LoadFontFace(path, points)
	if err != nil {
		return err
	}
	ctx.font = font
	ctx.points = points
	return nil
}

func (ctx *Context) SetTheme(theme themer.Theme) {
	ctx.theme = theme
}

func (ctx *Context) SetLanguage(lang tokenizer.Language) {
	ctx.lang = lang
}

func (ctx *Context) EncodePNG(code string, writer io.Writer) error {
	dc := ctx.parse(code)
	return dc.EncodePNG(writer)
}

func (ctx *Context) SavePNG(code string, output string) error {
	dc := ctx.parse(code)
	return dc.SavePNG(output)
}

func (ctx *Context) parse(code string) *gg.Context {
	const margin float64 = 10
	code = removeAccents(code)
	width, height := ctx.measureMultilineString(code)

	fmt.Printf("Width %.2f, Height %.2f\n", width, height)

	// gg
	dc := gg.NewContext(int(width+margin*2), int(height+margin*2))

	// background
	dc.SetHexColor(ctx.theme["body"])
	dc.Clear()
	// base color
	dc.SetHexColor(ctx.theme[tokenizer.COMMENT])

	// font
	dc.SetFontFace(ctx.font)

	// tokens
	tokens := tokenizer.Tokenize(code, ctx.lang)

	for _, token := range tokens {
		color := ctx.theme.GetColor(token)
		dc.SetHexColor(color)
		dc.DrawString(token.Content, margin+float64(token.Col)*ctx.measureString(" "), float64(margin+ctx.points+((1.3*ctx.points)*float64(token.Line))))
	}
	return dc
}

func (ctx *Context) measureMultilineString(s string) (width, height float64) {
	lineSpacing := 1.3
	lines := strings.Split(s, "\n")

	// sync h formula with DrawStringWrapped
	height = float64(len(lines)) * ctx.points * lineSpacing
	height -= (lineSpacing - 1) * ctx.points

	d := &font.Drawer{
		Face: ctx.font,
	}

	// max width from lines
	for _, line := range lines {
		adv := d.MeasureString(line)
		currentWidth := float64(adv >> 6) // from gg.Context.MeasureString
		if currentWidth > width {
			width = currentWidth
		}
	}

	return width, height
}

// measureString returns the rendered width and height of the specified text
// given the current font face.
func (ctx *Context) measureString(s string) float64 {
	d := &font.Drawer{
		Face: ctx.font,
	}
	a := d.MeasureString(s)
	return float64(a >> 6)
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}
