package screen

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stanjansen/pubsubui/internal/pubsub"
	"github.com/stanjansen/pubsubui/internal/ui/theme"
	"golang.design/x/clipboard"
)

type messageList struct {
	screen        *Screen
	msgChan       <-chan pubsub.Message
	loading       bool
	ackedMessages map[string]bool
	messages      map[string]pubsub.Message
}

func (s *Screen) drawSubscription() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	msgs := messageList{
		screen:        s,
		loading:       true,
		messages:      make(map[string]pubsub.Message),
		ackedMessages: make(map[string]bool),
	}

	go func() {
		<-ctx.Done()
		msgs.loading = false
		if s.subscription != "" {
			msgs.drawSubscriptionTitle()
			s.RefreshApp()
		}
	}()

	s.contextCancels = append(s.contextCancels, cancel)
	go msgs.load(ctx)

	s.KeyActions.Replace(tcell.KeyEsc, ' ', func() bool {
		s.subscription = ""
		s.KeyActions.Remove(tcell.KeyRune, 'r')
		s.KeyActions.Remove(tcell.KeyRune, 'a')
		s.KeyActions.Remove(tcell.KeyRune, 'c')
		s.Redraw()
		return true
	})
	s.KeyActions.Add("Reload", tcell.KeyRune, 'r', func() bool {
		cancel()
		s.RedrawApp()
		return true
	})
	s.KeyActions.Add("Ack", tcell.KeyRune, 'a', func() bool {
		msgId := s.table.GetCell(s.table.GetSelection()).Text
		if _, ok := msgs.ackedMessages[msgId]; !ok && msgId != "" {
			msgs.ackedMessages[msgId] = true
			msgs.messages[msgId].Ack()
			row, _ := s.table.GetSelection()
			for cell := 0; cell <= 2; cell++ {
				s.table.GetCell(row, cell).SetTextColor(theme.MutedColor)
			}
		}
		return true
	})
	s.KeyActions.Add("Copy data", tcell.KeyRune, 'c', func() bool {
		row, _ := s.table.GetSelection()
		data := s.table.GetCell(row, 2).Text
		if data != "" {
			clipboard.Write(clipboard.FmtText, []byte(data))
		}
		return true
	})

	s.drawSubscriptionTable(&msgs)
}

func (s *Screen) drawSubscriptionTable(msgs *messageList) {
	s.table.Clear()
	s.table.SetFixed(1, 1)
	s.table.SetCell(0, 0, tview.NewTableCell("ID").SetAlign(tview.AlignLeft).SetExpansion(1).SetSelectable(false))
	s.table.SetCell(0, 1, tview.NewTableCell("Timestamp").SetAlign(tview.AlignLeft).SetExpansion(1).SetSelectable(false))
	s.table.SetCell(0, 2, tview.NewTableCell("Data").SetAlign(tview.AlignLeft).SetExpansion(10).SetSelectable(false))

	go func() {
		i := 0
		for msg := range msgs.msgChan {
			msgs.messages[msg.ID] = msg
			i++
			s.table.SetCell(i, 0, tview.NewTableCell(msg.ID).SetAlign(tview.AlignLeft).SetExpansion(1).SetTextColor(theme.ItemColor))
			s.table.SetCell(i, 1, tview.NewTableCell(msg.Timestamp.Format(time.RFC3339)).SetAlign(tview.AlignLeft).SetExpansion(1).SetTextColor(theme.ItemColor))
			s.table.SetCell(i, 2, tview.NewTableCell(string(msg.Data)).SetAlign(tview.AlignLeft).SetExpansion(10).SetTextColor(theme.ItemColor))
			s.RefreshApp()
		}
	}()

	msgs.drawSubscriptionTitle()
}

func (msgs *messageList) drawSubscriptionTitle() {
	statusColor := "lightblue"
	if msgs.loading {
		statusColor = "red"
	}

	msgs.screen.table.SetTitle(fmt.Sprintf(" Subscription [%s::b]<%s> ", statusColor, msgs.screen.subscription)).SetBorder(true)
}

func (msgs *messageList) load(ctx context.Context) {
	s := msgs.screen
	s.RefreshApp()

	// TODO: Handle error
	msgs.msgChan = s.Pubsub.Messages(ctx, s.subscription)
	s.drawSubscriptionTable(msgs)
	s.RefreshApp()
}
