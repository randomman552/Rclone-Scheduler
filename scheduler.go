package main

import (
	"log"

	"github.com/go-co-op/gocron/v2"
)

func createScheduler(backupSchedule string) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return s, err
	}

	_, err = s.NewJob(
		gocron.CronJob("* * * * * *", true),
		gocron.NewTask(
			func() {
				log.Default().Printf("Hello World!")
			},
		))

	if err != nil {
		return s, err
	}

	return s, nil
}
