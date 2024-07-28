package main

// Notifier interface
type Notifer interface {
	IsEnabled() bool
	NotifyReady(NotifyReadyContext)
	NotifyBackupStarted(NotifyBackupStartedContext)
	NotifyBackupFinished(NotifyBackupFinishedContext)
}

// Notification context for application ready
type NotifyReadyContext struct {
	Schedule         string
	NextBackupTime   string
	DurationToBackup string
}

// Notification context for backup started notifications
type NotifyBackupStartedContext struct {
	StartedAt string
	JobId     int
}

// Notification constant for backup finished notifications
type NotifyBackupFinishedContext struct {
	JobId     int
	Success   bool
	Duration  string
	Bytes     string
	Checks    uint
	Deletes   uint
	Transfers uint
	Errors    uint
	Renames   uint
	Speed     string
}
