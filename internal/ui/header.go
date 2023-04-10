package ui

import (
	"code.rocketnine.space/tslocum/cview"
)

type header struct {
	view   *cview.TextView
	colors *colors
}

func newHeader(colors *colors) *header {
	h := &header{}
	h.colors = colors
	h.view = cview.NewTextView()
	h.view.SetText("Audiobook Creator (Internet Archive version)")
	h.view.SetBorder(false)
	h.view.SetTextColor(h.colors.textColor)
	h.view.SetBackgroundColor(colors.textBgColor)
	return h
}
