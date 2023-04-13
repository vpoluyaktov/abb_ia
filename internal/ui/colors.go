package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
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
	// Text
	cview.Styles.PrimaryTextColor = gray           // Primary text.
	cview.Styles.SecondaryTextColor = yellow       // Secondary text (e.g. labels).
	cview.Styles.TertiaryTextColor = black         // Tertiary text (e.g. subtitles, notes).
	cview.Styles.InverseTextColor = black          // Text on primary-colored backgrounds.
	cview.Styles.ContrastPrimaryTextColor = black  // Primary text for contrasting elements.
	cview.Styles.ContrastSecondaryTextColor = gray // Secondary text on ContrastBackgroundColor-colored backgrounds.

	// Background
	cview.Styles.PrimitiveBackgroundColor = blue    // Main background color for primitives.
	cview.Styles.ContrastBackgroundColor = white    // Background color for contrasting elements.
	cview.Styles.MoreContrastBackgroundColor = gray // Background color for even more contrasting elements.
	cview.Styles.ScrollBarColor = white             // Scroll bar color
}
