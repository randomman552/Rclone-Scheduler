package main

import (
	"log"
	"math"
	"time"

	"github.com/dustin/go-humanize"
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

// Task run on startup
func AppliationReadyTask(c *cli.Context, backupJob gocron.Job) {
	backupSchedule := getBackupSchedule(c)
	backupNow := c.Bool("backup.now")
	nextRun, _ := backupJob.NextRun()
	nextRunStr := nextRun.Format(time.RFC822)
	timeBeforeNextJobStr := humanize.Time(nextRun)

	log.Printf("Backing up with schedule '%s'", backupSchedule)
	log.Printf("First backup will start %s (%s)", timeBeforeNextJobStr, nextRunStr)

	if backupNow {
		log.Printf("Starting backup immedietly")
	}

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

	// Start the backup job straight away if requested
	if backupNow {
		backupJob.RunNow()
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
			if jobStatus.Success {
				log.Printf("Finished backup job with id '%d' successfully", jobId)
			} else {
				log.Printf("Finished backup job with id '%d' unsuccessfully", jobId)
			}

			// Clear from storage engine, as the job is now done
			storageEngine.SetValue("currentJobId", nil)
			jobStats := rclone.GetSyncStats(jobId)

			// Send notifications
			gotify := NewGotifyNotifier(c)
			if gotify.IsEnabled() {
				roundDuration := time.Duration(time.Second)
				jobDuration := time.Duration(jobStatus.Duration * float64(time.Second)).Round(roundDuration)

				// Prevent divide by 0 error
				averageTransferSpeed := jobStats.Speed
				if uint(jobStats.TransferTime) != 0 {
					averageTransferSpeed = float64(jobStats.Bytes / uint(jobStats.TransferTime))
				}

				context := NotifyBackupFinishedContext{
					JobId:                jobId,
					Success:              jobStatus.Success,
					Duration:             jobDuration.String(),
					DataTransferred:      humanize.IBytes(uint64(jobStats.Bytes)),
					AverageSpeed:         humanize.IBytes(uint64(math.Round(jobStats.Speed))) + "/S",
					AverageTransferSpeed: humanize.IBytes(uint64(averageTransferSpeed)) + "/S",
					Checks:               humanize.Comma(int64(jobStats.Checks)),
					Deletes:              humanize.Comma(int64(jobStats.Deletes)),
					Transfers:            humanize.Comma(int64(jobStats.Transfers)),
					Errors:               humanize.Comma(int64(jobStats.Errors)),
					Renames:              humanize.Comma(int64(jobStats.Renames)),
				}

				gotify.NotifyBackupFinished(context)
			}
		}
	}
}
