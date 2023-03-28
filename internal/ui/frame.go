package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

func Frame() {
  app := cview.NewApplication()
  defer app.HandlePanic()

  app.EnableMouse(true)

  box := cview.NewBox()
  box.SetBackgroundColor(tcell.ColorBlue.TrueColor())

  frame := cview.NewFrame(box)
  frame.SetBorders(2, 2, 2, 2, 4, 4)
  frame.AddText("Header left", true, cview.AlignLeft, tcell.ColorWhite.TrueColor())
  frame.AddText("Header middle", true, cview.AlignCenter, tcell.ColorWhite.TrueColor())
  frame.AddText("Header right", true, cview.AlignRight, tcell.ColorWhite.TrueColor())
  frame.AddText("Header second middle", true, cview.AlignCenter, tcell.ColorRed.TrueColor())
  frame.AddText("Footer middle", false, cview.AlignCenter, tcell.ColorGreen.TrueColor())
  frame.AddText("Footer second middle", false, cview.AlignCenter, tcell.ColorGreen.TrueColor())

  app.SetRoot(frame, true)
  if err := app.Run(); err != nil {
    panic(err)
  }
}