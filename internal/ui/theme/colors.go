package theme

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	BackgroundColor = tcell.ColorBlack
	MainColor       = tcell.ColorBlue
	TextColor       = tcell.ColorWhite
	SecondaryColor  = tcell.ColorDeepPink
	TitleColor      = tcell.ColorOrange
	MutedColor      = tcell.ColorDarkSlateGray
	ItemColor       = tcell.ColorLightBlue
	ErrorColor      = tcell.ColorRed
)

func SetTheme() {
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    BackgroundColor,
		ContrastBackgroundColor:     MainColor,
		MoreContrastBackgroundColor: SecondaryColor,
		BorderColor:                 MainColor,
		TitleColor:                  TitleColor,
		GraphicsColor:               TitleColor,
		PrimaryTextColor:            TextColor,
		SecondaryTextColor:          MainColor,
		TertiaryTextColor:           SecondaryColor,
		InverseTextColor:            MainColor,
		ContrastSecondaryTextColor:  MainColor,
	}
}
