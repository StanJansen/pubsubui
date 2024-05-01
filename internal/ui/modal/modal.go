package modal

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stanjansen/pubsubui/internal/ui/keyactions"
)

type modal struct {
	view       tview.Primitive
	keyActions *keyactions.KeyActions
	pages      *tview.Pages
}

func NewModal(k *keyactions.KeyActions, pages *tview.Pages, p tview.Primitive, title string, width, height int) modal {
	view := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)

	return modal{
		view:       view,
		keyActions: k,
		pages:      pages,
	}
}

func (m modal) Open() {
	m.pages.AddPage("modal", m.view, true, true)
	m.keyActions.Replace(tcell.KeyEsc, ' ', m.Close)
	m.pages.ShowPage("modal")
}

func (m modal) Close() bool {
	m.pages.RemovePage("modal")

	return true
}
