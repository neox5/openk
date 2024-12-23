package cli

import (
	"context"
	"flag"
	"fmt"
)

func Execute(ctx context.Context) error {
	flag.Parse()

	if len(flag.Args()) < 1 {
		return fmt.Errorf("no command specified")
	}

	switch flag.Arg(0) {
	case "version":
		return version(ctx)
	case "server":
		return serverStart(ctx)
	default:
		return fmt.Errorf("unknown command: %s", flag.Arg(0))
	}
}
