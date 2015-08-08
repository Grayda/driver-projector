package main

import (
	"fmt"

	"github.com/Grayda/go-dell"
	"github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/suit"
)

// This file contains most of the code for the UI (i.e. what appears in the Labs)

type configService struct {
	driver *ProjectorDriver
}

// This function is common across all UIs, and is called by the Sphere. Shows our menu option on the main Labs screen
// The "c" bit at the start means this func is an extension of the configService struct (like prototyping, I think?)
func (c *configService) GetActions(request *model.ConfigurationRequest) (*[]suit.ReplyAction, error) {
	// What we're going to show
	var screen []suit.ReplyAction
	// Loop through all Orvibo devices. We do this so we can find an AllOne

	screen = append(screen, suit.ReplyAction{
		Name:        "",
		Label:       "Configure Projectors",
		DisplayIcon: "video-camera",
	},
	)
	// We've found at least one AllOne. That's enough to show the UI, so we stop looping

	// Return our screen to the sphere-ui for rendering
	return &screen, nil
}

// When you click on a ReplyAction button (e.g. the "Configure AllOne" button defined above), Configure is called. requests.Action == the "Name" portion of the ReplyAction
func (c *configService) Configure(request *model.ConfigurationRequest) (*suit.ConfigurationScreen, error) {
	fmt.Sprintf("Incoming configuration request. Action:%s Data:%s", request.Action, string(request.Data))

	switch request.Action {
	case "list": // Listing the IR codes
		fmt.Println("Showing list of IR codes..")
		return c.list()

	case "": // Coming in from the main menu
		return c.list()

	default: // Everything else

		// return c.list()
		return c.error(fmt.Sprintf("Unknown action: %s", request.Action))
	}

	// If this code runs, then we done fucked up, because default: didn't catch. When this code runs, the universe melts into a gigantic heap. But
	// removing this violates Apple guidelines and ensures the downfall of humanity (probably) so I don't want to risk it.
	// Then again, I could be making all this up. Do you want to remove it and try? ( ͡° ͜ʖ ͡°)
	return nil, nil
}

// So this function (which is an extension of the configService struct that suit (or Sphere-UI) requires) creates a box with a single "Okay" button and puts in a title and text

// The meat of our UI. Shows a list of IR codes to be blasted. This could show anything you like, really.
func (c *configService) list() (*suit.ConfigurationScreen, error) {

	// ActionListOption are buttons within a section that are sent to c.Configuration

	// Sections, for logical grouping
	var sections []suit.Section
	// Loop through all the CodeGroups in our driver
	for _, p := range dell.Projectors {

		// Now that we've looped through the codes for this group, create our UI for that section
		sections = append(sections, suit.Section{ // Append a new suit.Section into our sections variable

			Contents: []suit.Typed{ // Again, dunno what this means
				suit.StaticText{ // Create static text (a heading, basically)
					Title: p.Model + " " + p.Make, // With the name of the IR code group
					Value: p.UUID,                 // And it's description
				},
			},
		})
	}
	// Now that we've looped and got our sections, it's time to build the actual screen
	screen := suit.ConfigurationScreen{
		Title:    "Saved IR Codes",
		Sections: sections, // Our sections. Contains all the buttons and everything!
		Actions: []suit.Typed{ // Actiosn you can take on this page
			suit.CloseAction{ // Here we go! This takes a label and probably a DisplayIcon and DisplayClass and just takes you back to the main screen. Not YOUR main screen though, so use a ReplyAction with a "" name to go back to YOUR menu
				Label: "Close",
			},
			suit.ReplyAction{ // Reply action. Same as the rest
				Label:        "Manually Add Projector",
				Name:         "new", // Back in c.Configuration, show the new code UI
				DisplayClass: "success",
				DisplayIcon:  "asterisk",
			},
		},
	}

	return &screen, nil
}

func (c *configService) error(message string) (*suit.ConfigurationScreen, error) {

	return &suit.ConfigurationScreen{
		Sections: []suit.Section{
			suit.Section{
				Contents: []suit.Typed{
					suit.Alert{
						Title:        "Error",
						Subtitle:     message,
						DisplayClass: "danger",
					},
				},
			},
		},
		Actions: []suit.Typed{
			suit.ReplyAction{ // Shows a button we can click on. Takes us back to c.Configuration (reply.Action will be "list")
				Label:        "Cancel",
				Name:         "list",
				DisplayClass: "success",
				DisplayIcon:  "ok",
			},
		},
	}, nil
}

// // Shows the UI to learn a new IR code
// func (c *configService) new(config *ProjectorDriverConfig) (*suit.ConfigurationScreen, error) {
//
// 	// What radio options we're going to show. Are you seeing a pattern here now? The UI is rather easy once you do it for a while
// 	// If you want to know what options the UI supports, and what values you can use with them, check out https://github.com/ninjasphere/go-ninja/blob/master/suit/screen.go
// 	var allones []suit.RadioGroupOption
// 	var groups []suit.RadioGroupOption
//
// 	// Add a new RadioGroupOption to our list. This one blasts from All AllOnes connected ("ALL" is a special MAC Address in go-orvibo)
// 	allones = append(allones, suit.RadioGroupOption{
// 		Title:       "All Connected AllOnes",
// 		Value:       "ALL",
// 		DisplayIcon: "globe",
// 	})
//
// 	// Loop through the groups we've got
// 	for _, codegroup := range driver.config.CodeGroups {
// 		groups = append(groups, suit.RadioGroupOption{ // Add a new radio buton
// 			Title:       codegroup.Name,
// 			Value:       codegroup.Name,
// 			DisplayIcon: "folder-open",
// 		},
// 		)
// 	}
//
// 	// Loop through all Orvibo devices.
// 	for _, allone := range driver.device {
// 		// If it's an AllOne
// 		if allone.Device.DeviceType == orvibo.ALLONE {
// 			// Add a Radio button with our AllOne's name and MAC Address
// 			allones = append(allones, suit.RadioGroupOption{
// 				Title:       allone.Device.Name,
// 				DisplayIcon: "play",
// 				Value:       allone.Device.MACAddress,
// 			},
// 			)
//
// 		}
// 	}
//
// 	title := "New IR Code" // Up here for readability
//
// 	screen := suit.ConfigurationScreen{
// 		Title: title,
// 		Sections: []suit.Section{ // New array of sections
// 			suit.Section{ // New section
// 				Contents: []suit.Typed{
// 					suit.StaticText{ // Some introductory text
// 						Title: "About this screen",
// 						Value: "Please enter a name and a description for this code. You must also pick an AllOne. When you're ready, click 'Start Learning' and press a button on your remote",
// 					},
// 					suit.InputHidden{ // Not actually used by my code, but you can use InputHidden to pass stuff back to c.Configure()
// 						Name:  "id",
// 						Value: "",
// 					},
// 					suit.InputText{ // Textbox
// 						Name:        "name",
// 						Before:      "Name for this code",
// 						Placeholder: "TV On", // Placeholder is the faded text that appears inside a textbox, giving you a hint as to what to type in
// 						Value:       "",
// 					},
// 					suit.InputText{
// 						Name:        "description",
// 						Before:      "Code Description",
// 						Placeholder: "Living Room TV On",
// 						Value:       "",
// 					},
// 					suit.RadioGroup{
// 						Title:   "Select an AllOne to blast from",
// 						Name:    "allone",
// 						Options: allones, // We created our radio group before, and now we put it in here
// 					},
// 					suit.RadioGroup{
// 						Title:   "Select a group to add this code to",
// 						Name:    "group",
// 						Options: groups, // Same with our code groups
// 					},
// 				},
// 			},
// 		},
// 		Actions: []suit.Typed{
// 			suit.ReplyAction{ // This is not a CloseAction, because we want to go back to the list of IR codes, not back to the main menu. Hence why we use a ReplyAction with "list"
// 				Label:        "Cancel",
// 				Name:         "list",
// 				DisplayClass: "default",
// 			},
// 			suit.ReplyAction{
// 				Label:        "Start Learning",
// 				Name:         "save",
// 				DisplayClass: "success",
// 				DisplayIcon:  "star",
// 			},
// 		},
// 	}
//
// 	return &screen, nil
// }

// You know the drill. I don't think it even needs to accept an *OrviboDriverConfig, because you could just call driver.config
// func (c *configService) newgroup(config *OrviboDriverConfig) (*suit.ConfigurationScreen, error) {
//
// 	title := "New Code Group"
// 	// New screen
// 	screen := suit.ConfigurationScreen{
// 		Title: title,
// 		Sections: []suit.Section{
// 			suit.Section{
// 				Contents: []suit.Typed{
// 					suit.StaticText{
// 						Title: "About this screen",
// 						Value: "On this page you can create a new group to put your codes in. For example, you might create a group called 'Living Room' to store codes relating to your home theater in your living room",
// 					},
// 					suit.InputHidden{
// 						Name:  "id",
// 						Value: "",
// 					},
// 					suit.InputText{
// 						Name:        "name",
// 						Before:      "Name for this group",
// 						Placeholder: "Home Theater",
// 						Value:       "",
// 					},
// 					suit.InputText{
// 						Name:        "description",
// 						Before:      "Description of this group",
// 						Placeholder: "Codes related to the home theater",
// 						Value:       "",
// 					},
// 				},
// 			},
// 		},
// 		Actions: []suit.Typed{
// 			suit.ReplyAction{
// 				Label:        "Cancel",
// 				Name:         "list",
// 				DisplayClass: "default",
// 			},
// 			suit.ReplyAction{
// 				Label:        "Save Group",
// 				Name:         "savegroup",
// 				DisplayClass: "success",
// 				DisplayIcon:  "star",
// 			},
// 		},
// 	}
//
// 	return &screen, nil
// }

// Aye-aye, captain.
// Not actually needed (?)
func i(i int) *int {
	return &i
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
