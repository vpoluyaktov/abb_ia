package ui

import (
	"code.rocketnine.space/tslocum/cview"
)

type footer struct {
	view   *cview.TextView
	colors *colors
}

func newFooter(colors *colors) *footer {
	f := &footer{}
	f.colors = colors
	f.view = cview.NewTextView()
	f.view.SetText("Version 0.0.1")
	f.view.SetBorder(false)
	f.view.SetTextColor(colors.textColor)
	f.view.SetBackgroundColor(colors.textBgColor)
	return f
}
