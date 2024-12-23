package cli

import (
	"context"
	"flag"
	"fmt"
)

func Execute(ctx context.Context) error {
	// Define version flags
	showVersion := flag.Bool("version", false, "Display version information")
	shortVersion := flag.Bool("v", false, "Display version number only")
	flag.Parse()

	// Handle version flags first
	if *showVersion || *shortVersion {
		return version(ctx, *shortVersion)
	}

	// Handle commands
	if len(flag.Args()) < 1 {
		return fmt.Errorf("no command specified")
	}

	switch flag.Arg(0) {
	case "version":
		return version(ctx, false)
	default:
		return fmt.Errorf("unknown command: %s", flag.Arg(0))
	}
}
