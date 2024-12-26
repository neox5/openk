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
				Name:    "host",
				Usage:   "Server host address",
				Value:   "0.0.0.0",
				EnvVars: []string{"OPENK_HOST"},
			},
			&cli.IntFlag{
				Name:    "grpc-port",
				Usage:   "gRPC server port",
				Value:   9090,
				EnvVars: []string{"OPENK_GRPC_PORT"},
			},
			&cli.IntFlag{
				Name:    "http-port",
				Usage:   "HTTP server port",
				Value:   8080,
				EnvVars: []string{"OPENK_HTTP_PORT"},
			},
		},
		Action: func(c *cli.Context) error {
			return app.StartServer()
		},
	}
}
