package ui

import (
	"code.rocketnine.space/tslocum/cview"
)

type header struct {
	view   *cview.TextView
	colors *colors
}

func newHeader(colors *colors) *header {

	header := &header{}
	header.colors = colors
	header.view = cview.NewTextView()
	header.view.SetText("Use arrow keys to navigate, press ? for help ")
	header.view.SetBorder(false)
	header.view.SetTextColor(header.colors.textColor)
	header.view.SetBackgroundColor(colors.textBgColor)

	return header
}
