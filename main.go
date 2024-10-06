package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/hellofresh/health-go/v5"
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
		&cli.BoolFlag{
			Name: "backup.now",
			EnvVars: []string{
				"BACKUP_NOW",
			},
			Value: false,
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
	go s.Start()

	// Create healthchecker
	h, _ := health.New(health.WithComponent(health.Component{
		Name:    "RClone Scheduler",
		Version: "v1.0",
	}),
		health.WithChecks(health.Config{
			Name: "RClone HTTP",
			Check: func(ctx context.Context) error {
				rclone := getRCloneClient(c)
				versionInfo := rclone.Version()

				if versionInfo == nil {
					return errors.New("RClone did not return successfully")
				}
				return nil
			},
		}))

	http.Handle("/health", h.Handler())
	http.ListenAndServe(":3000", nil)

	// Sleep until terminated
	select {}
}
