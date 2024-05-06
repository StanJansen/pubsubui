package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stanjansen/pubsubui/internal/ui/modal"
	"github.com/stanjansen/pubsubui/internal/ui/theme"
)

func (ui *ui) setDefaultKeyActions() {
	ui.keyActions.
		Add("Change project", tcell.KeyCtrlP, ' ', ui.openChangeProjectModal).
		Add("Back", tcell.KeyEsc, ' ', nil).
		Add("Quit", tcell.KeyCtrlC, ' ', nil).
		Add("Publish", tcell.KeyRune, 'p', ui.openPublishModal)
}

func (ui *ui) openChangeProjectModal() bool {
	form := theme.NewForm()

	form.
		SetTitle(" Change project ").
		SetBorder(true)

	modal := modal.NewModal(ui.keyActions, ui.pages, form, "Change project", 45, 7)

	inputField := tview.NewInputField().
		SetLabel("Project:").
		SetText(ui.app.GetProject())

	form.
		AddFormItem(inputField).
		AddButton("OK", func() {
			val := inputField.GetText()
			if val == "" || val == ui.app.GetProject() {
				modal.Close()
				return
			}

			err := ui.app.SetProject(val)
			if err != nil {
				form.
					SetTitle(" Error: " + err.Error() + " ").
					SetTitleColor(theme.ErrorColor).
					SetBorderColor(theme.ErrorColor)
				return
			}

			modal.Close()
			ui.reset()
		})

	modal.Open()

	return true
}

func (ui *ui) openPublishModal() bool {
	subscription := ui.screen.SelectedSubscription()
	if subscription == "" {
		return false
	}

	form := theme.NewForm()

	form.
		SetTitle(fmt.Sprintf(" Publish message [lightblue::b]<%s> ", subscription)).
		SetBorder(true)

	modal := modal.NewModal(ui.keyActions, ui.pages, form, "Publish message", 60, 11)

	inputField := tview.NewTextArea()

	form.
		AddFormItem(inputField).
		AddButton("OK", func() {
			content := inputField.GetText()
			ui.app.Pubsub().Publish(subscription, content)

			modal.Close()
		})

	modal.Open()

	return true
}
