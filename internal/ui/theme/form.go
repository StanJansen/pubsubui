package theme

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewForm() *tview.Form {
	form := tview.NewForm().
		SetLabelColor(TextColor).
		SetButtonBackgroundColor(MutedColor).
		SetButtonActivatedStyle(invertedStyle()).
		SetButtonsAlign(tview.AlignCenter).
		SetFieldBackgroundColor(MutedColor)

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
		case tcell.KeyUp:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
		}
		return event
	})

	return form
}

func invertedStyle() tcell.Style {
	return tcell.StyleDefault.
		Background(MainColor)
}
