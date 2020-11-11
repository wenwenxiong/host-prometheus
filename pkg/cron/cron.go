package cron

import(
	"github.com/go-co-op/gocron"
	"time"
)

type Cron struct {
	cron *gocron.Scheduler
}

func NewCron() (*Cron, error){
	var c Cron
	c.cron = gocron.NewScheduler(time.UTC)

	return &c, nil
}

func (c *Cron) Cron() (*gocron.Scheduler){
	return c.cron
}