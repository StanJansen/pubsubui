package keyactions

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type KeyActions struct {
	app          *tview.Application
	keyActions   []KeyAction
	onUpdateFunc func()
}

type KeyAction struct {
	Description string
	key         tcell.Key
	char        rune
	action      func() bool
}

func NewKeyActions(app *tview.Application) *KeyActions {
	k := KeyActions{
		app:        app,
		keyActions: make([]KeyAction, 0, 5),
	}

	k.setInputCapture()

	return &k
}

func (k *KeyActions) Add(description string, key tcell.Key, char rune, action func() bool) *KeyActions {
	k.Remove(key, char)
	k.keyActions = append(k.keyActions, KeyAction{
		Description: description,
		key:         key,
		char:        char,
		action:      action,
	})

	k.updated()

	return k
}

func (k *KeyActions) Remove(key tcell.Key, char rune) *KeyActions {
	for i, s := range k.keyActions {
		if s.key != key || s.char != char {
			continue
		}

		k.keyActions = append(k.keyActions[:i], k.keyActions[i+1:]...)
	}

	k.updated()

	return k
}

func (k *KeyActions) Replace(key tcell.Key, char rune, action func() bool) *KeyActions {
	for i, s := range k.keyActions {
		if s.key != key || s.char != char {
			continue
		}

		k.keyActions[i].action = action
	}

	k.updated()

	return k
}

func (k *KeyActions) GetAll() []KeyAction {
	return k.keyActions
}

func (k *KeyActions) GetAction(key tcell.Key, char rune) func() bool {
	for _, s := range k.keyActions {
		if s.key == key && s.char == char {
			return s.action
		}
	}

	return nil
}

func (k *KeyActions) setInputCapture() {
	k.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if _, ok := k.app.GetFocus().(tview.FormItem); ok {
			return event
		}

		for _, s := range k.keyActions {
			if s.key != event.Key() {
				continue
			}
			if s.key == tcell.KeyRune && s.char != event.Rune() {
				continue
			}

			if s.action == nil {
				return event
			}

			if ok := s.action(); !ok {
				return event
			}

			return nil
		}

		return event
	})
}

func (k KeyAction) Key() string {
	if k.key == tcell.KeyRune {
		return string(k.char)
	}

	return strings.ToLower(tcell.KeyNames[k.key])
}

func (k *KeyActions) OnUpdate(f func()) {
	k.onUpdateFunc = f
}

func (k *KeyActions) updated() {
	if k.onUpdateFunc != nil {
		k.onUpdateFunc()
	}
}
