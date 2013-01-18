package manager

import (
	"TOGY/config"
	"time"
)

const (
	startBroadcast = iota
	stopBroadcast
	block
	unblock
)

type Manager struct {
	broadcastChan chan int
	broadcastErr  chan error
	scheduleChan  chan bool
	reloadSignal  chan bool
	config        *config.Config
}

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
	}
	return
}

func (m *Manager) Run() {
	m.Start()
	<-m.reloadSignal
	m.scheduleChan <- true
	m.stopBroadcast()
	close(m.broadcastChan)
}

func (m *Manager) Start() {
	go broadcastManager(m)
	st := time.Tick(time.Second * 10)
	go scheduleManager(m, st)
	ut := time.Tick(m.config.UpdateInterval)
	go updateManager(m, ut)
}

func New(c *config.Config) (m *Manager) {
	m = new(Manager)
	m.broadcastChan = make(chan int)
	m.broadcastErr = make(chan error)
	m.scheduleChan = make(chan bool)
	m.reloadSignal = make(chan bool)
	m.config = c
	return
}

func (m *Manager) startBroadcast() error {
	return m.sendAndWaitForError(startBroadcast)
}

func (m *Manager) stopBroadcast() error {
	return m.sendAndWaitForError(stopBroadcast)
}

func (m *Manager) block() {
	m.broadcastChan <- block
}

func (m *Manager) unblock() {
	m.broadcastChan <- unblock
}

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
