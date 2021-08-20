package obs

import (
	"fmt"
	"strconv"

	"github.com/pterm/pterm"
	"github.com/richardcase/streamdeck/pkg/plugin"
	"github.com/spf13/viper"

	obsws "github.com/christopher-dG/go-obs-websocket"
)

func New() plugin.Plugin {
	return &Obs{}
}

type Obs struct {
	obs_client obsws.Client
}

func (o *Obs) Configure() error {
	viper.SetDefault("obs.host", "localhost")
	viper.SetDefault("obs.port", "4444")

	return nil
}

func (o *Obs) GetActions() ([]plugin.Action, error) {
	client, err := o.getClientAndConnect()
	if err != nil {
		return nil, err
	}

	scene_req := obsws.NewGetSceneListRequest()
	scenes, err := scene_req.SendReceive(*client)
	if err != nil {
		return nil, fmt.Errorf("getting obs scenes: %w", err)
	}

	actions := []plugin.Action{}
	for _, scene := range scenes.Scenes {
		action := &obsAction{
			id:          scene.Name,
			description: fmt.Sprintf("OBS Scene %s", scene.Name),
			client:      client,
		}

		actions = append(actions, action)
	}

	return actions, nil
}

func (o *Obs) getClientAndConnect() (*obsws.Client, error) {
	port, err := strconv.Atoi(viper.GetString("obs.port"))
	if err != nil {
		return nil, fmt.Errorf("obs port must be a number")
	}

	client := obsws.Client{
		Host: viper.GetString("obs.host"),
		Port: port,
	}

	err = client.Connect()
	if err != nil {
		return nil, fmt.Errorf("connecting to obs: %s", err)
	}

	return &client, nil
}

type obsAction struct {
	id          string
	description string
	client      *obsws.Client
}

func (a *obsAction) ID() string {
	return a.id
}

func (a *obsAction) Description() string {
	return a.description
}

func (a *obsAction) Pressed() {
	req := obsws.NewSetCurrentSceneRequest(a.id)
	_, err := req.SendReceive(*a.client)
	if err != nil {
		pterm.Warning.Print("cannot change to secene %s: %s", a.id, err)
	}
}
