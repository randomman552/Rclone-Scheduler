# Rclone-Scheduler
[![Build Status](https://drone.ggrainger.uk/api/badges/randomman552/Rclone-Scheduler/status.svg)](https://drone.ggrainger.uk/randomman552/Rclone-Scheduler)

Docker image to schedule a regular rclone backup using the rclone API

## Configuration
Configuration is done through environment variables

### Rclone connection variables
| Variable        | Description                                  | Default     |
|----------------:|:--------------------------------------------:|:------------|
| RCLONE_HOST     | Host in which the rclone daemon is running   | `localhost` |
| RCLONE_PORT     | Port in which the rclone daemon is listening | `5572`      |
| RCLONE_PROTOCOL | The protocol to use (`http`, or `https`)     | `https`     |

### Backup config variables
Configuration for backup operation.\
Bear in mind all paths are used by the Rclone daemon and not this script

| Variable           | Description                                        | Default     |
|-------------------:|:--------------------------------------------------:|:------------|
| BACKUP_SCHEDULE    | The cron schedule to run the backup on             | `0 0 * * 0` |
| BACKUP_SOURCE      | Source to get the data from                        | `/data`     |
| BACKUP_REMOTE      | The remote to use as a destination when backing up | `remote`    |
| BACKUP_DEST        | The destination of the backup on the remote        | `/backup`   |
| BACKUP_FILTER_FROM | The path to the file to use as a backup filter     |             |

### Restoring a backup
***TODO***\
Currently this setup does not support the restoration of files from a remote.