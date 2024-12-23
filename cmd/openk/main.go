package main

import (
	"context"
	"fmt"
	"os"

	"github.com/neox5/openk/internal/app"
	"github.com/neox5/openk/internal/cli"
)

func main() {
	// Create base application context
	ctx := app.NewContext(context.Background())

	if err := cli.Execute(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
