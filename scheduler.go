package main

import (
	"log"
	"strconv"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/urfave/cli/v2"
)

func createScheduler(c *cli.Context) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return s, err
	}

	// Storage engine for storing values between jobs
	storageEngine := NewMemoryStorageEngine()

	// Create task to start a backup
	backupSchedule := getBackupSchedule(c)
	backupJob, err := s.NewJob(
		gocron.CronJob(backupSchedule, false),
		gocron.NewTask(
			StartBackupTask,
			c,
			storageEngine,
		))

	if err != nil {
		log.Printf("Failed to create backup schedule: %s", err)
	}

	if err != nil {
		log.Printf("Failed to get next backup time: %s", err)
	}

	// Job to check sync status (every 5 seconds)
	_, err = s.NewJob(
		gocron.DurationJob(time.Duration(5*1000*1000*1000)),
		gocron.NewTask(
			CheckBackupStatusTask,
			c,
			storageEngine,
		))

	if err != nil {
		log.Printf("Failed to create check schedule: %s", err)
	}

	// Job to run once at startup
	_, err = s.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()),
		gocron.NewTask(
			AppliationReadyTask,
			c,
			backupJob,
		),
	)

	if err != nil {
		log.Printf("Failed to create startup job: %s", err)
	}

	return s, nil
}

func AppliationReadyTask(c *cli.Context, backupJob gocron.Job) {
	backupSchedule := getBackupSchedule(c)
	nextRun, _ := backupJob.NextRun()
	nextRunStr := nextRun.Format(time.RFC822)
	oneSecond := time.Duration(1 * 1000 * 1000 * 1000)
	timeBeforeNextJob := time.Until(nextRun)
	timeBeforeNextJobStr := timeBeforeNextJob.Round(oneSecond).String()

	log.Printf("Backing up with schedule '%s'", backupSchedule)
	log.Printf("First run in '%s' at '%s'", timeBeforeNextJobStr, nextRunStr)

	// Send gotify notification
	gotify := NewGotifyNotifier(c)
	if gotify.IsEnabled() {
		notifyContext := NotifyReadyContext{
			Schedule:         backupSchedule,
			NextBackupTime:   nextRunStr,
			DurationToBackup: timeBeforeNextJobStr,
		}

		gotify.NotifyReady(notifyContext)
	}
}

// The function used to start a backup job
func StartBackupTask(c *cli.Context, storageEngine MemoryStorageEngine) {
	storedJobId := storageEngine.GetValue("currentJobId")
	if storedJobId != nil {
		log.Printf("Backup job is already running. Skipping...")
		return
	}

	// Start backup
	client := getRCloneClient(c)
	srcPath := getBackupSourcePath(c)
	destPath := getBackupDestinationPath(c)

	res := client.StartSync(srcPath, destPath)
	if res != nil {
		log.Printf("Started backup job with id '%d'", res.JobId)
		storageEngine.SetValue("currentJobId", res.JobId)

		// Notification context
		context := NotifyBackupStartedContext{
			StartedAt: time.Now().Format(time.RFC822),
			JobId:     res.JobId,
		}

		// Send notifications
		gotify := NewGotifyNotifier(c)
		if gotify.IsEnabled() {
			gotify.NotifyBackupStarted(context)
		}
	}
}

// Task function used to check the status of the currently running backup
func CheckBackupStatusTask(c *cli.Context, storageEngine MemoryStorageEngine) {
	storedJobId := storageEngine.GetValue("currentJobId")

	if storedJobId != nil {
		jobId := storedJobId.(int)
		rclone := getRCloneClient(c)

		// Check the job status with rclone
		jobStatus := rclone.GetSyncStatus(jobId)
		if jobStatus.Finished {
			log.Printf("Finished backup job with id '%d'", jobId)

			// Clear from storage engine, as the job is now done
			storageEngine.SetValue("currentJobId", nil)

			// Notification context
			context := NotifyBackupFinishedContext{
				Duration: strconv.FormatFloat(jobStatus.Duration, 'f', 0, 64),
				JobId:    jobId,
			}

			// Send notifications
			gotify := NewGotifyNotifier(c)
			if gotify.IsEnabled() {
				gotify.NotifyBackupFinished(context)
			}
		}
	}
}
