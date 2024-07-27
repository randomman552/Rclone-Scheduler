package main

import (
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
			Name: "rclone.host",
			EnvVars: []string{
				"RCLONE_HOST",
			},
			Value: "localhost",
		},
		&cli.StringFlag{
			Name: "rclone.port",
			EnvVars: []string{
				"RCLONE_PORT",
			},
			Value: "5572",
		},
		&cli.StringFlag{
			Name: "rclone.protocol",
			EnvVars: []string{
				"RCLONE_PROTOCOL",
			},
			Value: "https",
		},
		&cli.StringFlag{
			Name: "backup.schedule",
			EnvVars: []string{
				"BACKUP_SCHEDULE",
			},
			Value: "0 0 * * 0",
		},
		&cli.StringFlag{
			Name: "backup.source",
			EnvVars: []string{
				"BACKUP_SOURCE",
			},
			Value: "/data",
		},
		&cli.StringFlag{
			Name: "backup.remote",
			EnvVars: []string{
				"BACKUP_REMOTE",
			},
			Value: "remomte",
		},
		&cli.StringFlag{
			Name: "backup.destination",
			EnvVars: []string{
				"BACKUP_DEST",
			},
			Value: "/backup",
		},
		&cli.StringFlag{
			Name: "gotify.url",
			EnvVars: []string{
				"GOTIFY_URL",
			},
		},
		&cli.StringFlag{
			Name: "gotify.token",
			EnvVars: []string{
				"GOTIFY_TOKEN",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	// Create scheduler
	s, err := createScheduler(c)
	if err != nil {
		return err
	}

	// Start scheduler
	s.Start()

	// Sleep until terminated
	select {}
}
