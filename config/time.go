package config

import (
	"github.com/sellweek/TOGY/util"
	"time"
)

//IsOverriddenDay tests if submitted time is during an overriden day.
func (c Config) IsOverriddenDay(t time.Time) bool {
	_, ok := c.OverrideDays[util.NormalizeDate(t).Unix()]
	return ok
}

//IsBroadcastingTime tests if it should broadcast on time with timeconfig.
func (tc TimeConfig) IsBroadcastingTime(t time.Time) bool {
	afterOn := util.NormalizeTime(t).After(tc.TurnOn)
	beforeOff := tc.TurnOff.After(util.NormalizeTime(t))
	return afterOn && beforeOff
}

//BroadcastingTime tests if according to the config there should be a broadcast on specified time.
func (c Config) BroadcastingTime(t time.Time) bool {
	if c.OverrideOn {
		return true
	}
	if c.OverrideOff {
		return false
	}
	if c.IsOverriddenDay(t) {
		return c.OverrideDays[util.NormalizeDate(t).Unix()].IsBroadcastingTime(t)
	}
	if c.Weekends {
		return c.StandardTimeSettings.IsBroadcastingTime(t)
	}

	return !isWeekend(t) && c.StandardTimeSettings.IsBroadcastingTime(t)
}

//Broadcast returns whether there should be broadcast at the current time.
func (c Config) Broadcast() bool {
	now := time.Now()
	return c.BroadcastingTime(now)
}

func isWeekend(t time.Time) bool {
	return 0 == t.Weekday() || 6 == t.Weekday()
}
