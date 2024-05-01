package app

import (
	"github.com/stanjansen/pubsubui/internal/pubsub"
	"github.com/stanjansen/pubsubui/internal/ui"
)

type app struct {
	Config
	pubsub *pubsub.Pubsub
}

type Config struct {
	PubsubEmulatorHost string
	Project            string
}

func New(c Config) app {
	return app{
		Config: c,
	}
}

func (a *app) Run() error {
	pubsub, err := pubsub.New(a.Project, a.PubsubEmulatorHost)
	if err != nil {
		return err
	}
	a.pubsub = pubsub

	err = ui.Render(a)
	if err != nil {
		return err
	}

	return pubsub.Close()
}

func (a *app) GetHost() string {
	return a.pubsub.Host()
}

func (a *app) GetVersion() string {
	return "v0.0.1"
}

func (a *app) GetProject() string {
	return a.pubsub.Project()
}

func (a *app) SetProject(p string) error {
	err := a.pubsub.UpdateProject(p)
	if err != nil {
		return err
	}

	return nil
}

func (a *app) Pubsub() *pubsub.Pubsub {
	return a.pubsub
}
