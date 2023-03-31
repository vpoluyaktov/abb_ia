package ui

import (
	"code.rocketnine.space/tslocum/cview"
)

type footer struct {
	view   *cview.TextView
	colors *colors
}

func newFooter(colors *colors) *footer {

	footer := &footer{}
	footer.colors = colors
	footer.view = cview.NewTextView()
	footer.view.SetText("Footer")
	footer.view.SetBorder(false)
	footer.view.SetTextColor(colors.textColor)
	footer.view.SetBackgroundColor(colors.textBgColor)

	return footer
}
