package screen

import (
	"context"
	"sync"

	"github.com/rivo/tview"
	"github.com/stanjansen/pubsubui/internal/pubsub"
	"github.com/stanjansen/pubsubui/internal/ui/keyactions"
)

type Config struct {
	Pubsub     *pubsub.Pubsub
	KeyActions *keyactions.KeyActions
	Pages      *tview.Pages
	RedrawApp  func()
	RefreshApp func()
}

type screen struct {
	*Config
	table         *tview.Table
	subscriptions []pubsub.Subscription
	subscription  string
	sync.Mutex
	contextCancels []context.CancelFunc
	loading        bool
}

func NewScreen(c *Config) *screen {
	table := tview.NewTable().
		SetSelectable(true, false)

	table.SetBorderPadding(0, 0, 1, 1)

	screen := screen{
		Config:         c,
		table:          table,
		subscriptions:  []pubsub.Subscription{},
		contextCancels: []context.CancelFunc{},
	}

	screen.Redraw()

	return &screen
}

func (s *screen) Primitive() tview.Primitive {
	return s.table
}

func (s *screen) Redraw() {
	for _, cancel := range s.contextCancels {
		cancel()
	}
	s.contextCancels = []context.CancelFunc{}

	s.table.Clear()

	if s.subscription != "" {
		s.drawSubscription()
	} else {
		s.drawProject()
	}
}

func (s *screen) Reset() {
	s.subscription = ""
	s.subscriptions = []pubsub.Subscription{}
	s.Redraw()
}
