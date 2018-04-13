package ato

import (
	"fmt"
	"github.com/reef-pi/reef-pi/controller/utils"
	"log"
)

type Usage struct {
	Pump int            `json:"pump"`
	Time utils.TeleTime `json:"time"`
}

func (u1 Usage) Rollup(ux utils.Metric) (utils.Metric, bool) {
	u2 := ux.(Usage)
	log.Printf("%v,%v", u1.Time, u2.Time)
	u := Usage{Time: u1.Time, Pump: u1.Pump}
	if u1.Time.Hour() == u2.Time.Hour() {
		u.Pump += u2.Pump
		return u, false
	}
	return u2, true
}

func (u1 Usage) Before(ux utils.Metric) bool {
	u2 := ux.(Usage)
	return u1.Time.Before(u2.Time)
}

func (c *Controller) NotifyIfNeeded(a ATO, u Usage) {
	if !a.Notify.Enable {
		return
	}
	if u.Pump >= a.Notify.Max {
		subject := "[Reef-Pi ALERT] ATO pump usage out of range"
		format := "Current usage of ATO pump (%d) is for sensor '%s' above acceptable value (%d)"
		body := fmt.Sprintf(format, u.Pump, a.Name, a.Notify.Max)
		c.telemetry.Alert(subject, "Elevated ATO usage. "+body)
	}
}
