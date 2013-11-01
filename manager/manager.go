package manager

import (
	"github.com/sellweek/TOGY/config"
	"os"
	"time"
)

const (
	startBroadcast bMsg = iota
	stopBroadcast
	block
	unblock
)

//bMsg is the type of messages used
//to control the broadcast manager.
type bMsg int

//Manager is a struct used to manage
//the broadcast can be started and stopped
//using it.
type Manager struct {
	//Chan used to send commands to broadcast manager.
	broadcastChan chan bMsg
	//Chan used to send errors back from broadcast manager.
	broadcastErr chan error
	//Chan used to send termination signal to schedule manager.
	scheduleChan chan bool
	//Chan used to signal
	reloadSignal chan bool
	//Chan used to send signals that the config is updated
	//and managers should restart.
	config               *config.Config
	currentPresentations []string
}

//RunBroadcast loads a local config file specified in cp,
//starts up and keeps running the broadcast,
//handling all the updates and scheduling.
func RunBroadcast(cp string) (err error) {
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
		c.Notice("Manager running")
	}
	return
}

//m.Run starts broadcast, handling
//scheduling. It shuts down when the
//config is updated.
func (m *Manager) Run() {
	m.Start()
	<-m.reloadSignal
	m.scheduleChan <- true
	m.stopBroadcast()
	close(m.broadcastChan)
}

//Start starts broadcast, screen and update manager.
func (m *Manager) Start() {
	go broadcastManager(m)
	st := time.Tick(time.Second * 10)
	go scheduleManager(m, st)
	//ut := time.Tick(m.config.UpdateInterval)
	//go updateManager(m, ut)
}

//New returns a new Manager with all the chans initialized
//and containing a list of broadcasts available in broadcast
//folder.
func New(c *config.Config) (m *Manager, err error) {
	m = new(Manager)
	m.broadcastChan = make(chan bMsg)
	m.broadcastErr = make(chan error)
	m.scheduleChan = make(chan bool)
	m.reloadSignal = make(chan bool)
	m.config = c
	m.currentPresentations, err = getBroadcastDirs(c)
	return
}

func getBroadcastDirs(c *config.Config) (ids []string, err error) {
	dir, err := os.Open(c.BroadcastDir)
	if err != nil {
		return
	}

	ids, err = dir.Readdirnames(0)
	return
}

//Starts the handler application and turns the screen on.
func (m *Manager) startBroadcast() error {
	return m.sendAndWaitForError(startBroadcast)
}

//Stops the handler application and turns the screen off.
func (m *Manager) stopBroadcast() error {
	return m.sendAndWaitForError(stopBroadcast)
}

//Blocks the broadcast manager from
//starting or stopping the handler application,
//turning the screen on and off.
func (m *Manager) block() {
	m.broadcastChan <- block
}

//Allows broadcast manager to start or stop
//the handler application and turn the screen on and off.
func (m *Manager) unblock() {
	m.broadcastChan <- unblock
}

//Sends a message to the broadcast manager,
//returning an erro, if it occurs there.
func (m *Manager) sendAndWaitForError(msg bMsg) error {
	m.broadcastChan <- msg
	err := <-m.broadcastErr
	return err
}
