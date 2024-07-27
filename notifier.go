package main

import "time"

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
	NextBackupTime   time.Time
	DurationToBackup time.Duration
}

// Notification context for backup started notifications
type NotifyBackupStartedContext struct {
	StartedAt string
	JobId     int
}

// Notification constant for backup finished notifications
type NotifyBackupFinishedContext struct {
	Duration string
	JobId    int
}
