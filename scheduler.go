package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/urfave/cli/v2"
)

// Notification context for backup started notifications
type BackupStartedContext struct {
	StartedAt string
	JobId     int
}

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
				// Start backup
				client := getRCloneClient(c)
				srcPath := getBackupSourcePath(c)
				destPath := getBackupDestinationPath(c)

				res := client.StartSync(srcPath, destPath)
				log.Printf("Started backup job with id '%d'", res.JobId)

				// Notification context
				context := BackupStartedContext{
					StartedAt: time.Now().Format(time.RFC822),
					JobId:     res.JobId,
				}

				// Send notifications
				gotify := NewGotifyNotifier(c)
				if gotify.IsEnabled() {
					gotify.NotifyBackupStarted(context)
				}
			},
		))

	if err != nil {
		return s, err
	}

	return s, nil
}
