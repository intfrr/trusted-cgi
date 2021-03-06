package application

import (
	"context"
	"github.com/robfig/cron"
	"log"
	"time"
)

func (project *Project) RunCron(ctx context.Context) {
	now := time.Now()
	last := project.lastScheduler
	for _, app := range project.CloneApps() {
		app.RunScheduled(ctx, last, now, project.GlobalEnvironment())
	}
	project.lastScheduler = now
}

func (app *App) RunScheduled(ctx context.Context, last, now time.Time, globalEnv map[string]string) {
	for _, plan := range app.Manifest.Cron {
		sched, err := cron.Parse(plan.Cron)
		if err != nil {
			log.Println(app.UID, plan.Cron, "-", err)
			continue
		}
		if !sched.Next(last).After(now) {
			log.Println("execution", app.UID, plan.Action)
			_, err = app.InvokeAction(ctx, plan.Action, time.Duration(plan.TimeLimit), globalEnv)
			if err != nil {
				log.Println(app.UID, plan.Cron, plan.Action, err)
			}
		}
	}
}
