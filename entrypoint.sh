#!/bin/sh

# Create crontab
touch rclone
echo "${BACKUP_SCHEDULE} /home/rclone/run.py" > /etc/crontabs/rclone
echo Running with cron schedule "${BACKUP_SCHEDULE}"

# Start cron
crond -f -l 8 > /dev/stdout