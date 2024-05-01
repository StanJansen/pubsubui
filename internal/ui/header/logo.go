package header

import (
	"strings"

	"github.com/rivo/tview"
	"github.com/stanjansen/pubsubui/internal/ui/theme"
)

var asciiLogo = []string{
	`__________    __  _________`,
	`\______   \  / / /   _____/`,
	` |     ___/ / /  \_____  \ `,
	` |____|    /_/  /________/`,
}

func (h *header) logoView() tview.Primitive {
	return tview.NewTextView().
		SetTextColor(theme.TitleColor).
		SetDynamicColors(true).
		SetText("[::b]" + strings.Join(asciiLogo, "\n"))
}
