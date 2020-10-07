package ui

import (
	"image"
	"strings"

	. "github.com/gizak/termui/v3"
)

const (
	spaceSizeFromTopBorder  = 1
	spaceSizeFromLeftBorder = 2
	capSizeTitleAndFirstRow = 2 + spaceSizeFromTopBorder
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
	topRow      int

	SelectedRow int
}

func NewTable() *Table {
	return &Table{
		Block:       NewBlock(),
		Cursor:      true,
		topRow:      0,
		SelectedRow: 0,
	}
}

func (self *Table) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	if self.drawable() {
		// store start positions for each column
		var (
			colPos []int
			cur int = spaceSizeFromLeftBorder
		)
		for _, w := range self.Widths {
			colPos = append(colPos, cur)
			cur += w
		}

		// draw headers
		for i, h := range self.Headers {
			// replace to '…' if the field is over
			h := TrimString(h, self.Widths[i] - spaceSizeFromLeftBorder)
			buf.SetString(
				h,
				NewStyle(Theme.Default.Fg, ColorClear, ModifierBold),
				image.Pt(
					self.Inner.Min.X+colPos[i],
					self.Inner.Min.Y+spaceSizeFromTopBorder),
			)
		}

		if self.SelectedRow < self.topRow {
			self.topRow = self.SelectedRow
		} else if self.SelectedRow > self.cursorBottom() {
			self.topRow = self.cursorBottom() + spaceSizeFromTopBorder
		}

		// describe rows
		for idx := self.topRow; idx >= 0 && idx < len(self.Rows) && idx <= self.bottom(); idx++ {
			row := self.Rows[idx]
			// move y+1 for a header
			y := self.Inner.Min.Y + 1 + idx - self.topRow + spaceSizeFromTopBorder
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
				r := TrimString(row.Elems[i], width-spaceSizeFromLeftBorder)
				buf.SetString(
					r,
					style,
					image.Pt(self.Inner.Min.X+colPos[i], y),
				)
			}
		}
	}
}

func (self *Table) drawable() bool {
	return self.Inner.Dy() >= capSizeTitleAndFirstRow
}

func (self *Table) cursorBottom() int {
	return self.topRow + self.Inner.Dy() - capSizeTitleAndFirstRow
}

func (self *Table) bottom() int {
	return self.topRow + self.Inner.Dy() - capSizeTitleAndFirstRow
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
