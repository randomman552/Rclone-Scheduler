package main

import (
	"log"

	"github.com/go-co-op/gocron/v2"
	"github.com/urfave/cli/v2"
)

func createScheduler(c *cli.Context) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return s, err
	}

	// Create scheduled backup task
	backupSchedule := getBackupSchedule(c)
	_, err = s.NewJob(
		gocron.CronJob(backupSchedule, true),
		gocron.NewTask(
			func() {
				client := getRCloneClient(c)
				srcPath := getBackupSourcePath(c)
				destPath := getBackupDestinationPath(c)

				res := client.StartSync(srcPath, destPath)
				log.Printf("Started backup job with id '%d'", res.JobId)
			},
		))

	if err != nil {
		return s, err
	}

	return s, nil
}
