package auth

import (
	"github.com/neox5/openk/internal/app/client"
	"github.com/urfave/cli/v2"
)

func newLoginCommand() *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "Login with existing credentials",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "Username for login",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "Password for login",
			},
		},
		Action: loginAction,
	}
}

func loginAction(c *cli.Context) error {
	username, err := readUsername(c.String("username"))
	if err != nil {
		return err
	}

	password, err := readPassword(c.String("password"), false)
	if err != nil {
		return err
	}

	opts := client.LoginOptions{
		Username: username,
		Password: password,
	}

	return client.Login(c.Context, opts)
}
