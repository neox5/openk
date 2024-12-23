package cli

import (
	"context"

	"github.com/neox5/openk/internal/app"
)

func serverStart(ctx context.Context) error {
	return app.StartServer(ctx)
}
