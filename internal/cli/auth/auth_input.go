package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/neox5/openk/internal/opene"
	"golang.org/x/term"
)

func readUsername(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}

	var username string
	for {
		fmt.Print("Enter username: ")
		fmt.Scanln(&username)
		username = strings.TrimSpace(username)
		if username == "" {
			fmt.Println("Username cannot be empty. Please try again.")
			continue
		}
		return username, nil
	}
}

func readPassword(flagValue string, confirm bool) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}

	for {
		fmt.Print("Enter password: ")
		pass, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", opene.NewInternalError("auth", "read_password", "failed to read password").
				Wrap(opene.AsError(err, "auth", opene.CodeInternal))
		}
		if len(pass) == 0 {
			fmt.Println("Password cannot be empty. Please try again.")
			continue
		}

		if confirm {
			fmt.Print("Confirm password: ")
			confirm, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				return "", opene.NewInternalError("auth", "read_password", "failed to read password confirmation").
					Wrap(opene.AsError(err, "auth", opene.CodeInternal))
			}

			if string(pass) != string(confirm) {
				fmt.Println("Passwords do not match. Please try again.")
				continue
			}
		}
		return string(pass), nil
	}
}
