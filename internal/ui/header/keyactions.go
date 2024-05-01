package header

import (
	"github.com/rivo/tview"
	"github.com/stanjansen/pubsubui/internal/ui/keyactions"
	"github.com/stanjansen/pubsubui/internal/ui/theme"
)

func (h *header) keyActionsView() tview.Primitive {
	drawTable := func(actions []keyactions.KeyAction) tview.Primitive {
		table := tview.NewTable().
			SetBorders(false)

		i := 0
		for _, k := range actions {
			table.SetCell(i, 0, tview.NewTableCell("[::b]<"+k.Key()+">").
				SetTextColor(theme.SecondaryColor))

			table.SetCell(i, 1, tview.NewTableCell(k.Description))
			i++
		}

		return table
	}

	actions := h.KeyActions.Get()
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	if len(actions) < 4 {
		flex.AddItem(drawTable(actions), 0, 1, false)
	} else {
		flex.AddItem(drawTable(actions[:3]), 0, 1, false).
			AddItem(drawTable(actions[3:]), 0, 1, false)
	}

	return flex
}
