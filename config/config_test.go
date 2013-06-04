package config_test

import (
	"github.com/sellweek/TOGY/config"
	"testing"
	"time"
)

var conf *config.Config

var tz = time.FixedZone("", -7200)

func init() {
	var err error
	conf, err = config.Get("./config.json")
	if err != nil {
		panic(err)
	}
}

func TestNormalBroadcast(t *testing.T) {
	tm := time.Date(2012, 10, 4, 9, 24, 0, 0, tz)
	if !conf.BroadcastingTime(tm) {
		t.Error("Does not broadcast on a normal day.")
	}
}

func TestNormalNotBroadcast(t *testing.T) {
	tm := time.Date(2012, 10, 4, 7, 00, 0, 0, tz)
	if conf.BroadcastingTime(tm) {
		t.Error("Does not broadcast on a normal day.")
	}
}

func TestWeekend(t *testing.T) {
	c := conf
	for i := 8; i < 2000; i += 7 {
		c.Weekends = false
		tm := time.Date(2012, 9, i, 13, 0, 0, 0, tz)
		if conf.BroadcastingTime(tm) {
			t.Error("Broadcasts during the weekend.")
		}
		c.Weekends = true
		if !conf.BroadcastingTime(tm) {
			t.Error("Does not broadcast during the weekend. Date: " + tm.String())
		}
	}
}

func TestOverrideDay(t *testing.T) {
	tm := time.Date(2012, 10, 7, 0, 0, 0, 0, tz)
	if !conf.IsOverridenDay(tm) {
		t.Error("Did not recognize overriden date")
	}
}

func TestOverriddenNotBroadcast(t *testing.T) {
	tm := time.Date(2012, 10, 7, 0, 0, 0, 0, tz)
	if conf.BroadcastingTime(tm) {
		t.Error("Broadcasted out of set time on overridden date.")
	}
}

func TestOverriddenBroadcast(t *testing.T) {
	tm := time.Date(2012, 10, 7, 9, 0, 0, 0, tz)
	if !conf.IsOverridenDay(tm) {
		t.Error("Did not broadcast on overridden date.")
	}
}
