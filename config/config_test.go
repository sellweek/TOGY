package config_test

import (
	"TOGY/config"
	"testing"
	"time"
)

var conf config.Config
var tz, _ = time.LoadLocation("UTC")

func init() {
	var err error
	conf, err = config.Get("./config.json")
	if err != nil {
		panic(err)
	}
}

func TestNormalBroadcast(t *testing.T) {
	tm := time.Date(2012, 10, 4, 9, 24, 0, 0, tz)
	if !conf.BroadcastTime(tm) {
		t.Error("Does not broadcast on a normal day.")
	}
}

func TestNormalNotBroadcast(t *testing.T) {
	tm := time.Date(2012, 10, 4, 7, 00, 0, 0, tz)
	if conf.BroadcastTime(tm) {
		t.Error("Does not broadcast on a normal day.")
	}
}

func TestWeekend(t *testing.T) {
	tm := time.Date(2012, 9, 8, 13, 0, 0, 0, tz)
	if conf.BroadcastTime(tm) {
		t.Error("Broadcasts during the weekend.")
	}
}

func TestOverrideDay(t *testing.T) {
	tm := time.Date(2012, 10, 7, 0, 0, 0, 0, tz)
	if !conf.IsOverridenDay(tm) {
		t.Error("Did not recognize overriden date")
	}
}

func TestOverridenNotBroadcast(t *testing.T) {
	tm := time.Date(2012, 10, 7, 0, 0, 0, 0, tz)
	if conf.BroadcastTime(tm) {
		t.Error("Broadcasted out of set time on overriden date.")
	}
}

func TestOverridenBroadcast(t *testing.T) {
	tm := time.Date(2012, 10, 7, 9, 0, 0, 0, tz)
	if !conf.IsOverridenDay(tm) {
		t.Error("Did not broadcast on overriden date.")
	}
}
