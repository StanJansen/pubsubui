package pubsub

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

type Pubsub struct {
	client        *pubsub.Client
	emulatorHost  string
	subscriptions map[string]*pubsub.SubscriptionConfig
}

func New(project, emulatorHost string) (*Pubsub, error) {
	if emulatorHost != "" {
		os.Setenv("PUBSUB_EMULATOR_HOST", emulatorHost)
	}

	if project == "" {
		if emulatorHost != "" {
			project = "emulator-project"
		} else {
			ctx := context.Background()
			credentials, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
			if err != nil {
				return nil, err
			}

			project = credentials.ProjectID
		}
	}

	ps := &Pubsub{
		emulatorHost:  emulatorHost,
		subscriptions: make(map[string]*pubsub.SubscriptionConfig),
	}

	if err := ps.UpdateProject(project); err != nil {
		return nil, err
	}

	return ps, nil
}

func (p *Pubsub) UpdateProject(project string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return err
	}

	if p.client != nil {
		p.client.Close()
	}

	p.client = client

	return nil
}

func (p *Pubsub) Project() string {
	return p.client.Project()
}

func (p *Pubsub) Host() string {
	if p.emulatorHost != "" {
		return p.emulatorHost
	}

	return "pubsub.googleapis.com"
}

func (p *Pubsub) Close() error {
	return p.client.Close()
}
