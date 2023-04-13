package ui

import (
	"code.rocketnine.space/tslocum/cview"
)

func newText(text string) *cview.TextView {
	tv := cview.NewTextView()
	tv.SetTextAlign(cview.AlignCenter)
	tv.SetText(text)
	return tv
}

func newButton(text string, f Fn) *cview.Button {
	bt := cview.NewButton(text)
	bt.SetRect(0, 0, 20, 1)
	bt.SetBorder(false)
	bt.SetSelectedFunc(f)
	return bt
}

func box(title string) *cview.Box {
  b := cview.NewBox()
  b.SetBorder(true)
  b.SetTitle(title)
  return b
}

