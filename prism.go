package prism

import (
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strings"
	"unicode"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/mrmarble/prism/themer"
	"github.com/mrmarble/prism/tokenizer"
	"golang.org/x/image/font"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Context struct {
	points      float64
	lineSpacing float64
	margin      int
	lines       int
	font        font.Face
	theme       themer.Theme
	lang        tokenizer.Language
}

type Options struct {
	LineNumbers bool
}

//go:embed fonts/FiraCode-Regular.ttf
var firaCode []byte

func NewContext() *Context {
	font, _ := parseFontFace(firaCode, 46)
	return &Context{points: 46, margin: 50, lineSpacing: 1.25, font: font, theme: themer.Dark()}
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

func (ctx *Context) calculate(code string, options Options) (width, height int, horizontalMargin, verticalMargin float64) {
	codeWidth, codeHeight := ctx.measureMultilineString(code)
	horizontalMargin = float64(ctx.margin)
	verticalMargin = horizontalMargin

	if options.LineNumbers {
		horizontalMargin += ctx.measureString(fmt.Sprint(ctx.lines)) + (horizontalMargin / 2)
	}

	return int(codeWidth) + ctx.margin*2, int(codeHeight + verticalMargin*2), horizontalMargin, verticalMargin
}

func (ctx *Context) parse(code string, options Options) *gg.Context {
	code = removeAccents(code)
	ctx.lines = strings.Count(code, "\n")

	width, height, hMargin, _ := ctx.calculate(code, options)

	// gg
	dc := gg.NewContext(width, height)

	// background
	dc.SetHexColor(ctx.theme["body"])
	dc.Clear()
	// base color
	dc.SetHexColor(ctx.theme[tokenizer.COMMENT])

	// font
	dc.SetFontFace(ctx.font)

	if options.LineNumbers {
		pad := int(math.Log10(float64(ctx.lines)) + 1)
		for i := 0; i < ctx.lines; i++ {
			dc.SetHexColor(ctx.theme[tokenizer.COMMENT])
			dc.DrawString(fmt.Sprintf("%*d", pad, i+1), float64(ctx.margin), float64(ctx.margin)+ctx.points+((ctx.points*ctx.lineSpacing)*float64(i)))
		}
	}

	// tokens
	tokens := tokenizer.Tokenize(code, ctx.lang)

	runeWidth := ctx.measureString(" ") // rethink this for variable-width fonts
	for _, token := range tokens {
		color := ctx.theme.GetColor(token)
		dc.SetHexColor(color)
		x := float64(hMargin) + float64(token.Col)*runeWidth
		y := float64(ctx.margin) + ctx.points + ((ctx.points * ctx.lineSpacing) * float64(token.Line))
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

func loadFontFace(path string, points float64) (font.Face, error) {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseFontFace(fontBytes, points)
}

func parseFontFace(fontBytes []byte, points float64) (font.Face, error) {
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face, nil
}
