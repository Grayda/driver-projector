package main

import (
	"fmt"
	"log" // Similar thing, I suppose?
	// For outputting stuff to the screen
	"time" // Used as part of "setInterval" and for pausing code to allow for data to come back

	"github.com/Grayda/go-dell"           // The magic part that lets us control sockets
	"github.com/ninjasphere/go-ninja/api" // Ninja Sphere API
	"github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/support"
)

// package.json is required, otherwise the app just exits and doesn't show any output
var info = ninja.LoadModuleInfo("./package.json")
var serial string           // Declared but not used?
var driver *ProjectorDriver // So we can access this in our configuration.go file

// Are we ready to rock? This is sphere-orvibo only code by the way. You don't need to do this in your own driver?
var ready = false
var started = false // Stops us from running theloop twice

// OrviboDriver holds info about our driver, including our configuration
type ProjectorDriver struct {
	support.DriverSupport
	config *ProjectorDriverConfig // This is how we save and load IR codes and such. Call this by using driver.config
	conn   *ninja.Connection
	device map[string]*dell.Projector // A list of devices we've found. This is in addition to the list go-orvibo maintains
}

// OrviboDriverConfig holds config info. The learningIR* stuff should be in its own struct, but I haven't got that far yet.
type ProjectorDriverConfig struct {
	Initialised bool // Has our driver run once before?
}

// No config provided? Set up some defaults
func defaultConfig() *ProjectorDriverConfig {

	return &ProjectorDriverConfig{
		Initialised: false,
	}
}

// NewDriver does what it says on the tin: makes a new driver for us to run. This is called through main.go
func NewProjectorDriver() (*ProjectorDriver, error) {

	// Make a new OrviboDriver. Ampersand means to make a new copy, not reference the parent one (so A = new B instead of A = new B, C = A)
	driver = &ProjectorDriver{}
	// Empty map of OrviboDevices
	driver.device = make(map[string]*dell.Projector)
	// Initialize our driver. Throw back an error if necessary. Remember, := is basically a short way of saying "var blah string = 'abcd'"
	err := driver.Init(info)

	if err != nil {
		log.Fatalf("Failed to initialize Projector driver: %s", err)
	}

	// Now we export the driver so the Sphere can find it
	err = driver.Export(driver)

	if err != nil {
		log.Fatalf("Failed to export Orvibo driver: %s", err)
	}

	// NewDriver returns two things, OrviboDriver, and an error if present
	return driver, nil
}

// Start is where the fun and magic happens! The driver is fired up and starts finding sockets
func (d *ProjectorDriver) Start(config *ProjectorDriverConfig) error {
	log.Printf("Driver Starting with config %v", config)

	d.config = config // Load our config

	if !d.config.Initialised { // No config loaded? Make one
		d.config = defaultConfig()
	}

	// This tells the API that we're going to expose a UI, and to run GetActions() in configuration.go
	d.Conn.MustExportService(&configService{d}, "$driver/"+info.ID+"/configure", &model.ServiceAnnouncement{
		Schema: "/protocol/configuration",
	})

	// If we've not started the driver
	if started == false {
		// Start a loop that handles everything this driver does (finding sockets, blasting IR etc.)
		// We put it in its own loop to keep the code neat
		theloop(d, config)
	}

	return d.SendEvent("config", config)
}

func theloop(d *ProjectorDriver, config *ProjectorDriverConfig) error {
	_, err := dell.Init()
	if err != nil {
		fmt.Println("Error preparing commands. Error is:", err)
	}

	for { // Loop forever
		select { // This lets us do non-blocking channel reads. If we have a message, process it. If not, check for UDP data and loop
		case msg := <-dell.Events:
			switch msg.Name {
			case "ready":
				fmt.Println("Ready to start listening for commands..")
				_, err = dell.Listen()
				if err != nil {
					fmt.Println(err)
				}
			case "projectorfound":
				_, err = dell.AddProjector(msg.ProjectorInfo)
				_, _ = newDevice(driver, driver.Conn, msg.ProjectorInfo)
				fmt.Println("Projector with UUID of " + msg.ProjectorInfo.UUID + " was found at " + msg.ProjectorInfo.IP + ". Make: " + msg.ProjectorInfo.Make + ". Model:" + msg.ProjectorInfo.Model + ". Revision:" + msg.ProjectorInfo.Revision)
			case "listening":
				fmt.Println("Listening for projectors via DDDP")
			case "projectoradded":
				fmt.Println("================== Adding Device")

			}
		}
	}
	return nil
}

// Stop does nothing. Though if it did, we could pass "quit" to theloop and clean up timers and such. The only way drivers are stopped now, are by force (reboot, Ctrl+C now)
func (d *ProjectorDriver) Stop() error {
	return fmt.Errorf("This driver does not support being stopped. YOU HAVE NO POWER HERE.")

}

// Analogous to Javascript's setInterval. Runs a function after a certain duration and keeps running it until "true" is passed to it
func setInterval(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}
