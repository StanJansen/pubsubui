package header

import (
	"github.com/rivo/tview"
	"github.com/stanjansen/pubsubui/internal/ui/theme"
)

func (h *header) infoView() tview.Primitive {
	info := [3][2]string{
		{"Host", h.App.GetHost()},
		{"Project", h.App.GetProject()},
		{"Version", h.App.GetVersion()},
	}

	table := tview.NewTable().
		SetBorders(false)

	for i, row := range info {
		for j, text := range row {
			if j == 0 {
				table.SetCell(i, j, tview.NewTableCell(text+":").
					SetTextColor(theme.MainColor))
			} else {
				table.SetCell(i, j, tview.NewTableCell("[::b]"+text))
			}
		}
	}

	return table
}
