package manager

import (
	"github.com/sellweek/TOGY/config"
	"os"
	"time"
)

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

func (m *Manager) Run() {
	m.Start()
	<-m.reload
	m.killBroadcast()
}

func New(c *config.Config) (m *Manager, err error) {
	m = new(Manager)
	m.config = c
	m.kbChan = make(chan bool)
	m.reload = make(chan bool)
	m.broadcastKilled = make(chan bool)
	m.activePresentations, err = getBroadcastDirs(c)
	return
}

func (m *Manager) Start() {
	m.startBroadcast()
	ut := time.Tick(m.config.UpdateInterval)
	go updateManager(m, ut)
}

func (m *Manager) killBroadcast() {
	m.config.Debug("Killing broadcast manager")
	m.kbChan <- true
	m.config.Debug("Waiting for broadcast manager response")
	<-m.broadcastKilled
	m.config.Debug("Broadcast manager killed")
}

func (m *Manager) startBroadcast() {
	st := time.Tick(time.Second * 10)
	go broadcastManager(m, st)
}

func getBroadcastDirs(c *config.Config) (ids []string, err error) {
	dir, err := os.Open(c.BroadcastDir)
	if err != nil {
		return
	}

	ids, err = dir.Readdirnames(0)
	return
}
