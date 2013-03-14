package manager

import (
	"github.com/sellweek/TOGY/config"
	"time"
)

const (
	startBroadcast = iota
	stopBroadcast
	block
	unblock
)

//Manager is a struct used to manage
//the broadcast can be started and stopped
//using it.
type Manager struct {
	//Chan used to send commands to broadcast manager.
	broadcastChan chan int
	//Chan used to send errors back from broadcast manager.
	broadcastErr chan error
	//Chan used to send termination signal to schedule manager.
	scheduleChan chan bool
	//Chan used to signal 
	reloadSignal chan bool
	//Chan used to send signals that the config is updated
	//and managers should restart.
	config *config.Config
}

//Run loads a local config file specified in cp
//and starts and keeps running the broadcast,
//handling updates and scheduling.
func Run(cp string) (err error) {
	for {
		var c *config.Config
		c, err = config.Get(cp)
		if err != nil {
			return
		}
		c.Log.Println("Reloading manager")
		mgr := New(c)
		mgr.Run()
		c.Log.Println("Manager running")
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
	ut := time.Tick(m.config.UpdateInterval)
	go updateManager(m, ut)
}

//New returns a new Manager with all the chans initialized.
func New(c *config.Config) (m *Manager) {
	m = new(Manager)
	m.broadcastChan = make(chan int)
	m.broadcastErr = make(chan error)
	m.scheduleChan = make(chan bool)
	m.reloadSignal = make(chan bool)
	m.config = c
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
func (m *Manager) sendAndWaitForError(msg int) error {
	m.broadcastChan <- msg
	c := time.After(time.Second)
	<-c
	select {
	case err := <-m.broadcastErr:
		return err
	default:
		return nil
	}
	return nil
}
