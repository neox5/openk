package cli

import (
	"fmt"
	"os"

	"github.com/neox5/openk/internal/buildinfo"
	"github.com/neox5/openk/internal/cli/auth"
	"github.com/neox5/openk/internal/cli/server"
	"github.com/urfave/cli/v2"
)

func Execute() error {
	app := &cli.App{
		Name:  "openK",
		Usage: "openK secret management system",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show version information",
			},
		},
		Commands: []*cli.Command{
			server.NewCommand(),
			auth.NewCommand(),
		},
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				fmt.Println(buildinfo.Get())
				return nil
			}
			return cli.ShowAppHelp(c)
		},
	}

	return app.Run(os.Args)
}
