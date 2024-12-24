package auth

import "github.com/urfave/cli/v2"

// NewCommand creates the auth command group
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "auth",
		Usage: "Authentication commands",
		Subcommands: []*cli.Command{
			newRegisterCommand(),
			newLoginCommand(),
		},
		Action: func(c *cli.Context) error {
			return cli.ShowSubcommandHelp(c)
		},
	}
}
