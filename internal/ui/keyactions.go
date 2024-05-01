package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stanjansen/pubsubui/internal/ui/modal"
	"github.com/stanjansen/pubsubui/internal/ui/theme"
)

func (ui *ui) setDefaultKeyActions() {
	ui.keyActions.
		Add("Change project", tcell.KeyRune, 'p', ui.openChangeProjectModal).
		Add("Back", tcell.KeyEsc, ' ', nil).
		Add("Quit", tcell.KeyCtrlC, ' ', nil)
}

func (ui *ui) openChangeProjectModal() bool {
	if _, ok := ui.view.GetFocus().(tview.FormItem); ok {
		return false
	}

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
