package ui

import (
	"image"
	"strings"

	. "github.com/gizak/termui/v3"
)

const (
	heightFromTopBorderToHeader   = 1
	heightTitle                   = 1
	heightRow                     = 1
	widthFromLeftBorderToContents = 2
)

type Row struct {
	Key   string
	Elems []string
}

type Table struct {
	*Block

	Headers     []string
	Widths      []int
	Rows        []Row
	Cursor      bool
	CursorColor Color

	SelectedRow    int
	DrawInitialRow int
}

func NewTable() *Table {
	return &Table{
		Block:  NewBlock(),
		Cursor: true,
	}
}

func (self *Table) Drawable() bool {
	return self.Inner.Dy() >= heightFromTopBorderToHeader+heightTitle+heightRow
}

func (self *Table) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	if self.Drawable() {
		// store start positions for each column
		var (
			colPos []int
			cur    int = widthFromLeftBorderToContents
		)
		for _, w := range self.Widths {
			colPos = append(colPos, cur)
			cur += w
		}

		// draw headers
		for i, h := range self.Headers {
			// replace to '…' if the field is over
			h := TrimString(h, self.Widths[i]-widthFromLeftBorderToContents)
			buf.SetString(
				h,
				NewStyle(Theme.Default.Fg, ColorClear, ModifierBold),
				image.Pt(
					self.Inner.Min.X+colPos[i],
					self.Inner.Min.Y+heightFromTopBorderToHeader),
			)
		}

		if self.SelectedRow < self.DrawInitialRow {
			self.DrawInitialRow = self.SelectedRow
		} else if self.SelectedRow >= self.DrawInitialRow+self.Inner.Dy()-(heightFromTopBorderToHeader+heightTitle) {
			self.DrawInitialRow += self.Inner.Dy() - (heightFromTopBorderToHeader + heightTitle)
		}

		// draw rows
		for idx := self.DrawInitialRow; idx >= 0 && idx < len(self.Rows) && idx < self.DrawInitialRow+self.Inner.Dy()-(heightFromTopBorderToHeader+heightTitle); idx++ {
			row := self.Rows[idx]
			// move y+1 for a header
			y := self.Inner.Min.Y + idx - self.DrawInitialRow + heightFromTopBorderToHeader + heightTitle
			style := NewStyle(Theme.Default.Fg)
			if self.Cursor {
				if idx == self.SelectedRow {
					style.Fg = self.CursorColor
					style.Modifier = ModifierReverse
					buf.SetString(
						strings.Repeat(" ", self.Inner.Dx()),
						style,
						image.Pt(self.Inner.Min.X, y),
					)
					self.SelectedRow = idx
				}
			}
			for i, width := range self.Widths {
				r := TrimString(row.Elems[i], width-widthFromLeftBorderToContents)
				buf.SetString(
					r,
					style,
					image.Pt(self.Inner.Min.X+colPos[i], y),
				)
			}
		}
	}
}

func (self *Table) scroll(i int) {
	self.SelectedRow += i
	maxRow := len(self.Rows) - 1
	if len(self.Rows) > 1 && self.SelectedRow < 0 {
		self.SelectedRow = maxRow
	} else if len(self.Rows) > 1 && self.SelectedRow > maxRow {
		self.SelectedRow = 0
	}
}

func (self *Table) ScrollUp() {
	self.scroll(-1)
}

func (self *Table) ScrollDown() {
	self.scroll(1)
}
