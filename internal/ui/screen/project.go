package screen

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/rivo/tview"
	"github.com/stanjansen/pubsubui/internal/pubsub"
	"github.com/stanjansen/pubsubui/internal/ui/theme"
)

func (s *screen) drawProject() {
	s.loading = true

	ctx, cancel := context.WithCancel(context.Background())
	s.contextCancels = append(s.contextCancels, cancel)
	go s.loadSubscriptions(ctx, cancel)

	s.drawProjectTable(cancel)
}

func (s *screen) drawProjectTable(cancel context.CancelFunc) {
	s.table.Clear()
	s.table.SetFixed(1, 1)
	s.table.SetCell(0, 0, tview.NewTableCell("NAME").SetAlign(tview.AlignLeft).SetExpansion(10).SetSelectable(false))
	s.table.SetCell(0, 1, tview.NewTableCell("TOPIC").SetAlign(tview.AlignLeft).SetExpansion(10).SetSelectable(false))
	s.table.SetCell(0, 2, tview.NewTableCell("DL TOPIC").SetAlign(tview.AlignLeft).SetExpansion(10).SetSelectable(false))

	for i, sub := range s.sortedSubscriptions() {
		s.table.SetCell(i+1, 0, tview.NewTableCell(sub.Name).SetAlign(tview.AlignLeft).SetExpansion(10).SetTextColor(theme.ItemColor))
		s.table.SetCell(i+1, 1, tview.NewTableCell(sub.Topic).SetAlign(tview.AlignLeft).SetExpansion(10).SetTextColor(theme.ItemColor))
		s.table.SetCell(i+1, 2, tview.NewTableCell(sub.DeadLetterTopic).SetAlign(tview.AlignLeft).SetExpansion(10).SetTextColor(theme.ItemColor))
	}
	s.table.SetSelectedFunc(func(row int, column int) {
		cancel()
		s.subscription = s.table.GetCell(row, 0).Text
		s.Redraw()
	})

	s.drawProjectTitle()
}

func (s *screen) drawProjectTitle() {
	statusColor := "lightblue"
	if s.loading {
		statusColor = "red"
	}

	s.table.SetTitle(fmt.Sprintf(" Subscriptions [%s::b]<%s> ", statusColor, s.Pubsub.Project())).SetBorder(true)
}

func (s *screen) loadSubscriptions(ctx context.Context, cancel context.CancelFunc) {
	load := func() {
		if s.subscription != "" {
			return
		}

		s.Lock()
		defer s.Unlock()

		s.loading = true
		s.drawProjectTitle()
		s.RefreshApp()

		// TODO: Handle error
		sub, _ := s.Pubsub.Subscriptions()
		s.subscriptions = sub
		s.loading = false
		if s.subscription != "" {
			// Doublecheck if not navigated in the meanwhile
			return
		}

		s.drawProjectTable(cancel)
		s.RefreshApp()
	}
	load()

	timer := time.NewTimer(10 * time.Second)
	select {
	case <-timer.C:
		load()
	case <-ctx.Done():
		return
	}
}

func (s *screen) sortedSubscriptions() []pubsub.Subscription {
	subs := make([]pubsub.Subscription, len(s.subscriptions))
	copy(subs, s.subscriptions)

	sort.Slice(subs, func(i, j int) bool {
		return subs[i].Name < subs[j].Name
	})

	return subs
}
