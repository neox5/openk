package cli

import (
	"context"
	"fmt"

	"github.com/neox5/openk/internal/buildinfo"
)

func version(ctx context.Context, short bool) error {
	if short {
		fmt.Println(buildinfo.Get().ShortVersion())
	} else {
		fmt.Println(buildinfo.Get())
	}
	return nil
}
