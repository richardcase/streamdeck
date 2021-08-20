package main

import (
	"image/color"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	_ "github.com/magicmonkey/go-streamdeck/devices"
	"github.com/pterm/pterm"

	"github.com/richardcase/streamdeck/pkg/plugin/obs"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	pterm.ThemeDefault.SectionStyle = *pterm.NewStyle(pterm.FgCyan)
	letters := pterm.NewLettersFromString("streamdeck")
	pterm.DefaultBigText.WithLetters(letters).Render()
	pterm.Info.Printfln("version 0.1")

	sd, err := streamdeck.Open()
	if err != nil {
		pterm.Error.Printf("failed connecting to streamdeck: %s", err)
		os.Exit(1)
	}
	pterm.Info.Printf("Using device: %s\n", sd.GetName())
	sd.ClearButtons()
	sd.SetBrightness(50)

	obs := obs.New()
	obs.Configure()
	obsActions, err := obs.GetActions()
	if err != nil {
		pterm.Error.Printf("failed to get obs actions %s", err)
		os.Exit(1)
	}

	for i, act := range obsActions {
		pterm.Info.Printf("button: %s, %s\n", act.ID(), act.Description())

		sd.WriteTextToButton(i, act.ID(), color.White, color.Black)
	}

	sd.ButtonPress(func(i int, d *streamdeck.Device, e error) {
		if err != nil {
			pterm.Error.Println(err)
			return
		}

		act := obsActions[i]
		act.Pressed()
	})

	<-sigChan
	pterm.Info.Println("Signal received sending cancel")
}
