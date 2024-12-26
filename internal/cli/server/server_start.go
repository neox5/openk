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
			// Server-specific flags would go here
		},
		Action: func(c *cli.Context) error {
			return app.StartServer()
		},
	}
}
