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

	NodeLimit Straight
	PodLimit  Straight
	Usage     Curved
}

type Straight struct {
	Value float64
	Label string
	Style Style
}

type Curved struct {
	Values []float64
	Label  string
	Style  Style
}

func NewGraph() *Graph {
	return &Graph{
		Block: NewBlock(),
		Usage: Curved{
			Values: make([]float64, 0),
		},
	}
}

func (self *Graph) Reset() {
	self.NodeLimit.Label = ""
	self.NodeLimit.Value = 0
	self.PodLimit.Label = ""
	self.PodLimit.Value = 0
	self.Usage.Label = ""
	self.Usage.Values = make([]float64, 0)
}

func (self *Graph) getY(val float64) int {
	dy := self.Inner.Max.Y - (self.Inner.Min.Y + 4)
	return int(float64(self.Inner.Max.Y) - float64(dy)*(val/self.NodeLimit.Value))
}

func (self *Graph) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	startPosXForLabel := self.Inner.Min.X + widthFromLeftBorder
	trimWidth := self.Inner.Max.X - (widthFromLeftBorder + 1)
	buf.SetString(
		TrimString(self.NodeLimit.Label, trimWidth),
		self.NodeLimit.Style,
		image.Pt(startPosXForLabel, self.Inner.Min.Y+1),
	)
	buf.SetString(
		TrimString(self.PodLimit.Label, trimWidth),
		self.PodLimit.Style,
		image.Pt(startPosXForLabel, self.Inner.Min.Y+2),
	)
	buf.SetString(
		TrimString(self.Usage.Label, trimWidth),
		self.Usage.Style,
		image.Pt(startPosXForLabel, self.Inner.Min.Y+3),
	)

	canvas := NewCanvas()
	canvas.Rectangle = self.Inner

	canvas.SetLine(
		image.Pt(self.Inner.Min.X*2, self.getY(self.NodeLimit.Value)*4),
		image.Pt(self.Inner.Max.X*2, self.getY(self.NodeLimit.Value)*4),
		self.NodeLimit.Style.Fg,
	)
	canvas.SetLine(
		image.Pt(self.Inner.Min.X*2, self.getY(self.PodLimit.Value)*4),
		image.Pt(self.Inner.Max.X*2, self.getY(self.PodLimit.Value)*4),
		self.PodLimit.Style.Fg,
	)

	if len(self.Usage.Values) > 0 {
		dest := self.getY(self.Usage.Values[len(self.Usage.Values)-1])
		for di := len(self.Usage.Values) - 1; di >= 0; di-- {
			src := self.getY(self.Usage.Values[di])
			canvas.SetLine(
				image.Pt((self.Inner.Min.X+di)*2, dest*4),
				image.Pt((self.Inner.Min.X+di+1)*2, src*4),
				self.Usage.Style.Fg,
			)
			dest = src
		}
	}
	canvas.Draw(buf)
}
