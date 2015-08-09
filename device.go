package main

import (
	"github.com/Grayda/go-dell"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/channels"
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

	// player.ApplyIsOn = func() (bool, error) {
	// 	return true, nil
	// }

	// player.ApplyGetPower = func() (bool, error) {
	// 	return true, nil
	// }

	// Volume Channel
	player.ApplyVolumeUp = func() error {
		_, err := dell.SendCommand(projector, dell.Commands.Volume.Up)
		projector.Volume += 1
		if err != nil {
			return err
		}

		return err
	}

	player.ApplyVolumeDown = func() error {
		_, err := dell.SendCommand(projector, dell.Commands.Volume.Down)
		projector.Volume -= 1
		if err != nil {
			return err
		}

		return err
	}

	// player.ApplyVolume = func(state *channels.VolumeState) error {
	//
	// 	return nil
	// }

	player.ApplyToggleMuted = func() error {
		if projector.VolumeMuted == true {
			dell.SendCommand(projector, dell.Commands.Volume.Unmute)
			projector.VolumeMuted = false
		} else {
			dell.SendCommand(projector, dell.Commands.Volume.Mute)
			projector.VolumeMuted = true
		}

		player.UpdateVolumeState(&channels.VolumeState{Muted: &projector.VolumeMuted})
		return err
	}

	// enable the volume channel, supporting mute (parameter is true)
	if err := player.EnableVolumeChannel(true); err != nil {
		player.Log().Errorf("Failed to enable volume channel: %s", err)
	}

	// on-off channel methods
	player.ApplyOff = func() error {
		player.UpdateOnOffState(false)
		dell.SendCommand(projector, dell.Commands.Power.Off)
		projector.PowerState = false
		return nil
	}

	player.ApplyOn = func() error {
		player.UpdateOnOffState(true)
		dell.SendCommand(projector, dell.Commands.Power.On)
		projector.PowerState = true
		return nil
	}

	player.ApplyToggleOnOff = func() error {
		if projector.PowerState == true {
			dell.SendCommand(projector, dell.Commands.Power.Off)
		} else {
			dell.SendCommand(projector, dell.Commands.Power.On)
		}
		return nil
	}

	// I can't find anywhere that the on/off states ever get set - on the sphereamid or in the app
	if err := player.EnableOnOffChannel("state"); err != nil {
		player.Log().Errorf("Failed to enable on-off channel: %s", err)
	}

	// NOTE: this is a workaround to get on/off when dragging to on/play or off/pause. Find a better way if possible
	// https://discuss.ninjablocks.com/t/mediaplayer-device-drivers/3776/2 (question asked)
	player.ApplyPlayPause = func(isPlay bool) error {
		if isPlay {
			player.UpdateControlState(channels.MediaControlEventPlaying)
			return player.ApplyOn()
		} else {
			player.UpdateControlState(channels.MediaControlEventPaused)
			return player.ApplyOff()
		}
	}

	if err := player.EnableControlChannel([]string{}); err != nil {
		player.Log().Errorf("Failed to enable control channel: %s", err)
	}

	return &Device{*player, &projector}, nil
}
