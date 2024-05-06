package header

import (
	"github.com/rivo/tview"
	"github.com/stanjansen/pubsubui/internal/ui/keyactions"
)

type app interface {
	GetHost() string
	GetVersion() string
	GetProject() string
}

type Config struct {
	App        app
	KeyActions *keyactions.KeyActions
}

type header struct {
	*Config
	flex *tview.Flex
}

func NewHeader(c *Config) header {
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	header := header{
		Config: c,
		flex:   flex,
	}

	header.Redraw()

	c.KeyActions.OnUpdate(header.Redraw)

	return header
}

func (h header) Primitive() tview.Primitive {
	return h.flex
}

func (h header) Redraw() {
	h.flex.Clear().
		AddItem(tview.NewBox(), 1, 1, false). // Spacing
		AddItem(h.infoView(), 0, 1, false).
		AddItem(tview.NewBox(), 3, 1, false). // Spacing
		AddItem(h.keyActionsView(), 0, 1, false).
		AddItem(tview.NewBox(), 3, 1, false). // Spacing
		AddItem(h.logoView(), 28, 1, false)
}

func (h header) Reset() {
	h.Redraw()
}
