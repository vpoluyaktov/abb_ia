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
	h.view.SetText("Use arrow keys to navigate, press ? for help ")
	h.view.SetBorder(false)
	h.view.SetTextColor(h.colors.textColor)
	h.view.SetBackgroundColor(colors.textBgColor)
	return h
}
