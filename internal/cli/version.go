package cli

import (
	"context"
	"fmt"

	"github.com/neox5/openk/internal/buildinfo"
)

func version(ctx context.Context) error {
	fmt.Println(buildinfo.Get())
	return nil
}
