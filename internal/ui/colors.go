package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	white  = tcell.NewRGBColor(255, 255, 255)
	gray   = tcell.NewRGBColor(202, 202, 202)
	blue   = tcell.NewRGBColor(12, 34, 184)
	yellow = tcell.ColorYellow
	black  = tcell.NewRGBColor(0, 0, 0)
	cyan   = tcell.NewRGBColor(80, 176, 189)

	footerFgColor = black
	footerBgColor = cyan

	headerFgColor = black
	headerBGColor = cyan
)

func setColorTheme() {
	// // Text
	tview.Styles.PrimaryTextColor = gray           // Primary text.
	tview.Styles.SecondaryTextColor = yellow       // Secondary text (e.g. labels).
	tview.Styles.TertiaryTextColor = black         // Tertiary text (e.g. subtitles, notes).
	tview.Styles.InverseTextColor = black          // Text on primary-colored backgrounds.
	tview.Styles.ContrastSecondaryTextColor = black  // Primary text for contrasting elements.
	tview.Styles.ContrastSecondaryTextColor = gray // Secondary text on ContrastBackgroundColor-colored backgrounds.

	// Background
	tview.Styles.PrimitiveBackgroundColor = blue    // Main background color for primitives.
	tview.Styles.ContrastBackgroundColor = white    // Background color for contrasting elements.
	tview.Styles.MoreContrastBackgroundColor = gray // Background color for even more contrasting elements.
	// tview.Styles.Scroll = white             // Scroll bar color
}
