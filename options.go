package prism

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
	"github.com/mrmarble/prism/tokenizer"
)

type Options struct {
	LineNumbers bool
	Header      bool
	Relative    bool
	Range       Range
}

func (o *Options) Apply(dc *gg.Context, ctx *Context) {
	if o.Header {
		radius := 15.0

		colors := [3]string{"#ff5f58", "#ffbd2e", "#18c132"}
		for i, color := range colors {
			dc.DrawCircle(float64(ctx.margin*(i+1)), float64(ctx.margin), radius)
			dc.SetHexColor(color)
			dc.Fill()
		}
	}

	if o.LineNumbers {
		pad := int(math.Log10(float64(ctx.lines)) + 1)

		for i := 0; i < ctx.lines; i++ {
			dc.SetHexColor(ctx.theme[tokenizer.COMMENT])

			x := float64(ctx.margin)
			y := ctx.vMargin + ctx.points + ((ctx.points * ctx.lineSpacing) * float64(i))
			lineNumber := i + 1

			if o.Relative {
				lineNumber = i + o.Range.Start
			}

			dc.DrawString(fmt.Sprintf("%*d", pad, lineNumber), x, y)
		}
	}
}
