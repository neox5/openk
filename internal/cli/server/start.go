package server

import (
	"github.com/neox5/openk/internal/app"
	"github.com/urfave/cli/v2"
)

func newStartCommand() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start the openK server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to config file",
				EnvVars: []string{"OPENK_CONFIG"},
			},
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "Log level (debug, info, warn, error)",
				Value:   "info",
				EnvVars: []string{"OPENK_LOG_LEVEL"},
			},
		},
		Action: runServer,
	}
}

func runServer(c *cli.Context) error {
	return app.StartServer(c.Context)
}
