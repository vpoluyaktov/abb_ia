package ui

import (
	"github.com/gdamore/tcell/v2"
)

type colors struct {
	textColor   tcell.Color
	textBgColor tcell.Color
}

func newColors() *colors {
	colors := colors{}
	colors.textColor = tcell.NewRGBColor(255, 255, 255)
	colors.textBgColor = tcell.NewRGBColor(0, 0, 255)
	return &colors
}
