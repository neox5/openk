package auth

import (
	"github.com/neox5/openk/internal/app/client"
	"github.com/urfave/cli/v2"
)

func newRegisterCommand() *cli.Command {
	return &cli.Command{
		Name:  "register",
		Usage: "Register a new user",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "Username for registration",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "Password for registration",
			},
		},
		Action: registerAction,
	}
}

func registerAction(c *cli.Context) error {
	username, err := readUsername(c.String("username"))
	if err != nil {
		return err
	}

	password, err := readPassword(c.String("password"), true)
	if err != nil {
		return err
	}

	opts := client.RegisterOptions{
		Username: username,
		Password: password,
	}
	
	return client.Register(c.Context, opts)
}
