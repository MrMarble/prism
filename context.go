package prism

import (
	_ "embed"
	"fmt"
	"io"
	"strings"

	"github.com/fogleman/gg"
	"github.com/mrmarble/prism/themer"
	"github.com/mrmarble/prism/tokenizer"
	"golang.org/x/image/font"
)

type Context struct {
	points      float64
	lineSpacing float64
	margin      int
	vMargin     float64
	hMargin     float64
	width       int
	height      int
	lines       int
	font        font.Face
	theme       themer.Theme
	lang        tokenizer.Language
}

//go:embed fonts/FiraCode-Regular.ttf
var firaCode []byte

func NewContext() *Context {
	points := 46.0
	margin := 50
	lineSpacing := 1.25
	font, _ := parseFontFace(firaCode, points) //nolint

	return &Context{points: points, margin: margin, lineSpacing: lineSpacing, font: font, theme: themer.Dark()}
}

func (ctx *Context) SetFontFace(path string, points float64) error {
	font, err := loadFontFace(path, points)
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

func (ctx *Context) EncodePNG(code string, writer io.Writer, options Options) error {
	dc := ctx.parse(code, options)
	return dc.EncodePNG(writer)
}

func (ctx *Context) SavePNG(code string, output string, options Options) error {
	dc := ctx.parse(code, options)
	return dc.SavePNG(output)
}

func (ctx *Context) calculateSize(code string, options Options) {
	codeWidth, codeHeight := ctx.measureMultilineString(code)
	ctx.hMargin = float64(ctx.margin)
	ctx.vMargin = ctx.hMargin

	if options.LineNumbers {
		ctx.hMargin += ctx.measureString(fmt.Sprint(ctx.lines)) + (ctx.hMargin / 2) //nolint
	}

	if options.Header {
		ctx.vMargin += 70
	}

	ctx.width = int(codeWidth) + ctx.margin + int(ctx.hMargin)
	ctx.height = int(codeHeight) + ctx.margin + int(ctx.vMargin)
}

func (ctx *Context) calculateLines(code string, options Options) string {
	ctx.lines = strings.Count(code, "\n") + 1

	if (Range{}) != options.Range {
		var err error

		code, err = substr(code, options.Range)
		if err != nil {
			panic(err)
		}

		if options.Relative {
			ctx.lines = 1 + options.Range.End + -options.Range.Start
		} else {
			ctx.lines = strings.Count(code, "\n") + 1
		}
	}

	return code
}

func (ctx *Context) parse(code string, options Options) *gg.Context {
	code = removeAccents(code)

	code = ctx.calculateLines(code, options)
	ctx.calculateSize(code, options)

	// gg
	dc := gg.NewContext(ctx.width, ctx.height)

	// background
	dc.SetHexColor(ctx.theme["body"])
	dc.Clear()
	// base color
	dc.SetHexColor(ctx.theme[tokenizer.COMMENT])

	// font
	dc.SetFontFace(ctx.font)

	options.Apply(dc, ctx)

	// tokens
	tokens := tokenizer.Tokenize(code, ctx.lang)

	runeWidth := ctx.measureString(" ") // rethink this for variable-width fonts

	for _, token := range tokens {
		color := ctx.theme.GetColor(token)
		dc.SetHexColor(color)

		x := ctx.hMargin + float64(token.Col)*runeWidth
		y := ctx.vMargin + ctx.points + ((ctx.points * ctx.lineSpacing) * float64(token.Line))

		dc.DrawString(token.Content, x, y)
	}

	return dc
}

func (ctx *Context) measureMultilineString(s string) (width, height float64) {
	lines := strings.Split(s, "\n")

	// sync h formula with DrawStringWrapped
	height = float64(len(lines)) * ctx.points * ctx.lineSpacing
	height -= (ctx.lineSpacing - 1) * ctx.points

	d := &font.Drawer{
		Face: ctx.font,
	}

	// max width from lines
	for _, line := range lines {
		adv := d.MeasureString(line)
		currentWidth := float64(adv >> 6) //nolint

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

	return float64(a >> 6) //nolint
}
