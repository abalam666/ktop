package ui

import (
	"image"
	"strings"

	. "github.com/gizak/termui/v3"
)

const (
	widthFromLeftBorder = 2
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
	// consider: space from border + header + initial row = 3
	return self.Inner.Dy() >= 3
}

func (self *Table) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	if self.Drawable() {
		// store start positions for each column
		var (
			colPos []int
			cur    int = widthFromLeftBorder
		)
		for _, w := range self.Widths {
			colPos = append(colPos, cur)
			cur += w
		}

		// draw headers
		for i, h := range self.Headers {
			// replace to '…' if the field is over
			h := TrimString(h, self.Widths[i]-widthFromLeftBorder)
			buf.SetString(
				h,
				NewStyle(Theme.Default.Fg, ColorClear, ModifierBold),
				image.Pt(
					self.Inner.Min.X+colPos[i],
					// consider: space from border = 1
					self.Inner.Min.Y+1),
			)
		}

		if self.SelectedRow < self.DrawInitialRow {
			self.DrawInitialRow = self.SelectedRow
		} else if self.SelectedRow >= self.DrawInitialRow+self.Inner.Dy()-2 {
			// consider: space from border + header = 2
			self.DrawInitialRow += self.Inner.Dy() - 2
		}

		// draw rows
		for idx := self.DrawInitialRow; idx >= 0 && idx < len(self.Rows) && idx < self.DrawInitialRow+self.Inner.Dy()-2; idx++ {
			row := self.Rows[idx]
			// consider: space from border + header = 2
			y := self.Inner.Min.Y + idx - self.DrawInitialRow + 2
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
					self.setselected(idx)
				}
			}
			for i, width := range self.Widths {
				r := TrimString(row.Elems[i], width-widthFromLeftBorder)
				buf.SetString(
					r,
					style,
					image.Pt(self.Inner.Min.X+colPos[i], y),
				)
			}
		}
	}
}

func (self *Table) setselected(idx int) {
	self.SelectedRow = idx
	max := len(self.Rows) - 1
	if max >= 0 && self.SelectedRow < 0 {
		self.SelectedRow = max
	} else if max >= 0 && self.SelectedRow > max {
		self.SelectedRow = 0
	}
}

func (self *Table) ScrollUp() {
	self.setselected(self.SelectedRow - 1)
}

func (self *Table) ScrollDown() {
	self.setselected(self.SelectedRow + 1)
}
