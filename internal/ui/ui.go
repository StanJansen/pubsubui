package ui

import (
	"github.com/rivo/tview"
	"github.com/stanjansen/pubsubui/internal/pubsub"
	"github.com/stanjansen/pubsubui/internal/ui/header"
	"github.com/stanjansen/pubsubui/internal/ui/keyactions"
	"github.com/stanjansen/pubsubui/internal/ui/screen"
	"github.com/stanjansen/pubsubui/internal/ui/theme"
)

type app interface {
	GetHost() string
	GetVersion() string
	GetProject() string
	SetProject(string) error
	Pubsub() *pubsub.Pubsub
}

type redrawablePrimitive interface {
	Primitive() tview.Primitive
	Redraw()
	Reset()
}

type ui struct {
	app        app
	view       *tview.Application
	pages      *tview.Pages
	keyActions *keyactions.KeyActions
	header     redrawablePrimitive
	screen     redrawablePrimitive
}

func Render(app app) error {
	theme.SetTheme()
	view := tview.NewApplication()
	pages := tview.NewPages()
	keyActions := keyactions.NewKeyActions(view)

	ui := ui{
		app:        app,
		view:       view,
		pages:      pages,
		keyActions: keyActions,
	}

	ui.setDefaultKeyActions()

	ui.screen = screen.NewScreen(&screen.Config{
		Pubsub:     app.Pubsub(),
		KeyActions: keyActions,
		Pages:      pages,
		RedrawApp:  ui.redraw,
		RefreshApp: ui.refresh,
	})

	ui.header = header.NewHeader(&header.Config{
		App:        app,
		KeyActions: keyActions,
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.header.Primitive(), 5, 1, false).
		AddItem(ui.screen.Primitive(), 0, 1, true)

	pages = pages.AddPage("main", flex, true, true)

	return view.SetRoot(pages, true).Run()
}

func (ui *ui) redraw() {
	ui.screen.Redraw()
	ui.header.Redraw()
}

func (ui *ui) reset() {
	ui.screen.Reset()
	ui.header.Reset()
}

func (ui *ui) refresh() {
	ui.view.Draw()
}
