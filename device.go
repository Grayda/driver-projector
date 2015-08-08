package main

import (
	"github.com/Grayda/go-dell"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/devices"
	"github.com/ninjasphere/go-ninja/model"
)

type Device struct {
	devices.MediaPlayerDevice
	projector *dell.Projector
}

func newDevice(driver ninja.Driver, conn *ninja.Connection, projector dell.Projector) (*Device, error) {

	player, err := devices.CreateMediaPlayerDevice(driver, &model.Device{
		NaturalID:     projector.UUID,
		NaturalIDType: "dell-projector",
		Name:          &projector.UUID,
		Signatures: &map[string]string{
			"ninja:manufacturer": "Dell",
			"ninja:productName":  "Projector",
			"ninja:thingType":    "mediaplayer",
			"ip:serial":          projector.UUID,
		},
	}, conn)

	if err != nil {
		return nil, err
	}

	// Volume Channel
	player.ApplyVolumeUp = func() error {
		_, err := dell.SendCommand(projector, dell.Commands.Volume.Up)
		return err

	}

	player.ApplyVolumeDown = func() error {
		_, err := dell.SendCommand(projector, dell.Commands.Volume.Down)
		return err

	}

	player.ApplyToggleMuted = func() error {
		if projector.Muted == false {
			_, err := dell.SendCommand(projector, dell.Commands.Volume.Mute)
			return err
		} else {
			_, err := dell.SendCommand(projector, dell.Commands.Volume.Unmute)
			return err
		}

		return err

	}

	// On-off Channel
	player.ApplyOff = func() error {
		player.UpdateOnOffState(false)
		_, err := dell.SendCommand(projector, dell.Commands.Power.Off)
		return err
	}

	player.ApplyOn = func() error {
		player.UpdateOnOffState(true)
		_, err := dell.SendCommand(projector, dell.Commands.Power.On)
		return err
	}

	if err := player.EnableOnOffChannel("state"); err != nil {
		player.Log().Fatalf("Failed to enable control channel: %s", err)
	}

	if err := player.EnableVolumeChannel(false); err != nil {
		player.Log().Fatalf("Failed to enable volume channel: %s", err)
	}

	if err := player.EnableControlChannel([]string{}); err != nil {
		player.Log().Fatalf("Failed to enable control channel: %s", err)
	}

	return &Device{*player, &projector}, nil
}
