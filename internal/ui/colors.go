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
	red    = tcell.ColorRed

	footerFgColor = black
	footerBgColor = cyan

	headerFgColor = black
	headerBGColor = cyan

	busyIndicatorFgColor = yellow
	busyIndicatorBgColor = cyan
)

func setColorTheme() {
	// Text colors
	tview.Styles.PrimaryTextColor = gray            // Primary text.
	tview.Styles.SecondaryTextColor = yellow        // Secondary text (e.g. labels).
	tview.Styles.TertiaryTextColor = white          // Tertiary text (e.g. subtitles, notes).
	tview.Styles.InverseTextColor = white           // Text on primary-colored backgrounds.
	tview.Styles.ContrastSecondaryTextColor = white // Primary text for contrasting elements.
	tview.Styles.ContrastSecondaryTextColor = gray  // Secondary text on ContrastBackgroundColor-colored backgrounds.

	// Background colors
	tview.Styles.PrimitiveBackgroundColor = blue     // Main background color for primitives.
	tview.Styles.ContrastBackgroundColor = gray      // Background color for contrasting elements.
	tview.Styles.MoreContrastBackgroundColor = white // Background color for even more contrasting elements.

	// Elements colors
	tview.Styles.BorderColor = gray     // Box borders.
	tview.Styles.TitleColor = yellow    // Box titles.
	tview.Styles.GraphicsColor = yellow // Graphics.
}
