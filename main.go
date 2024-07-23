package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "Rclone Scheduler"
	app.Usage = "Scheduler for RClone Daemon"
	app.Action = run

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "rclone.host",
			Required: true,
			EnvVars: []string{
				"RCLONE_HOST",
			},
			Value: "localhost",
		},
		&cli.StringFlag{
			Name:     "rclone.port",
			Required: true,
			EnvVars: []string{
				"RCLONE_PORT",
			},
			Value: "5572",
		},
		&cli.StringFlag{
			Name:     "rclone.protocol",
			Required: true,
			EnvVars: []string{
				"RCLONE_PROTOCOL",
			},
			Value: "https",
		},
		&cli.StringFlag{
			Name: "schedule",
			EnvVars: []string{
				"BACKUP_SCHEDULE",
			},
			Value: "0 0 * * 1",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	// Get parameters
	protocol := c.String("rclone.protocol")
	host := c.String("rclone.host")
	port := c.String("rclone.port")
	backupSchedule := c.String("schedule")

	_ = fmt.Sprintf("%s://%s:%s", protocol, host, port)
	log.Default().Printf("Backing up with schedule '%s'", backupSchedule)

	// Create scheduler
	s, err := createScheduler(backupSchedule)
	if err != nil {
		return err
	}

	// Start scheduler
	s.Start()

	// Sleep until terminated
	select {}
}
