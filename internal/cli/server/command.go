package server

import "github.com/urfave/cli/v2"

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Server management commands",
		Subcommands: []*cli.Command{
			newStartCommand(),
		},
		Action: func(c *cli.Context) error {
			return cli.ShowSubcommandHelp(c)
		},
	}
}
