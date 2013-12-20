//Package manager contains all the "moving parts" of TOGY client.
//It handles rotating, updating and partially also scheduling of the broadcasts.
//It spawns 3 goroutines:
//
//rotator, which handles rotation of different presentations.
//
//broadcastManager, which handles scheduling.
//It starts and stops the rotator and turns the screen on and off.
//
//updateManager, which handles periodically checks server for any changes and handles them.
package manager

import (
	"github.com/sellweek/TOGY/config"
	"os"
	"time"
)

//Manager is the struct containing chans used for communication between
//goroutines, the list of active presentations and a pointer to the current
//config.
type Manager struct {
	kbChan              chan bool
	broadcastKilled     chan bool
	reload              chan bool
	activePresentations []string
	config              *config.Config
}

//RunContinuosly loads a local config file specified in cp,
//starts up and keeps running the broadcast,
//handling all the updates and scheduling.
func RunContinuosly(cp string) (err error) {
	for {
		var (
			c   *config.Config
			mgr *Manager
		)
		c, err = config.Get(cp)
		if err != nil {
			return
		}
		c.Notice("Reloading manager")
		mgr, err = New(c)
		if err != nil {
			return
		}

		mgr.Run()
	}
	return
}

//Run starts the Manager and waits for a signal to reload config.
//At that point, it stops the goroutines and returns.
func (m *Manager) Run() {
	m.Start()
	<-m.reload
	m.killBroadcast()
}

//New returns a fully populated Manager, prepared to be run.
func New(c *config.Config) (m *Manager, err error) {
	m = new(Manager)
	m.config = c
	m.kbChan = make(chan bool)
	m.reload = make(chan bool)
	m.broadcastKilled = make(chan bool)
	m.activePresentations, err = getBroadcastDirs(c)
	return
}

//Start spawns goroutines used by Manager.
func (m *Manager) Start() {
	m.startBroadcast()
	ut := time.Tick(m.config.UpdateInterval)
	go updateManager(m, ut)
}

//killBroadcast stops all the goroutines and turns off the screen.
func (m *Manager) killBroadcast() {
	m.config.Debug("Killing broadcast manager")
	m.kbChan <- true
	m.config.Debug("Waiting for broadcast manager response")
	<-m.broadcastKilled
	m.config.Debug("Broadcast manager killed")
}

//startBroadcast starts broadcastManager.
func (m *Manager) startBroadcast() {
	st := time.Tick(time.Second * 10)
	go broadcastManager(m, st)
}

//getBroadcastDirs returns a slice containing names of folders in
//c.BroadcastDir
func getBroadcastDirs(c *config.Config) (ids []string, err error) {
	dir, err := os.Open(c.BroadcastDir)
	if err != nil {
		return
	}

	ids, err = dir.Readdirnames(0)
	return
}
