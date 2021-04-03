package ui

import (
	"image"

	. "github.com/gizak/termui/v3"
)

const (
	widthFromLeftBorder = 2
)

type Graph struct {
	*Block

	Header      string
	HeaderStyle Style

	Elements []Element
}

type Element struct {
	Data  []float64
	Label string
	Style Style
}

func NewGraph() *Graph {
	return &Graph{
		Block:       NewBlock(),
		HeaderStyle: NewStyle(Theme.Default.Fg, Theme.Default.Bg, ModifierBold),
	}
}

func (self *Graph) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	buf.SetString(
		TrimString(self.Header, self.Inner.Max.X-widthFromLeftBorder),
		self.HeaderStyle,
		image.Pt(
			self.Inner.Min.X+widthFromLeftBorder,
			self.Inner.Min.Y+1,
		),
	)

	if len(self.Elements) > 0 {

		for i, e := range self.Elements {
			buf.SetString(
				e.Label, e.Style,
				image.Pt(self.Inner.Min.X+2, self.Inner.Min.Y+2+i),
			)
		}

		canvas := NewCanvas()
		canvas.Rectangle = self.Inner
		for _, e := range self.Elements {
			// dest := self.height(ei, r - 1)
			for di := len(e.Data) - 1; di >= 0; di-- {
				// src := self.height(ei, di)
				canvas.SetLine(
					image.Pt(self.Inner.Min.X+di, 10),
					image.Pt(self.Inner.Min.X+di+1, 12),
					e.Style.Fg,
				)
				// dest = src
			}
		}
		canvas.Draw(buf)
	}
}
