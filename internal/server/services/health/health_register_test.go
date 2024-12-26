package health_test

import (
	"log/slog"
	"testing"

	"github.com/neox5/openk/internal/server/services/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestRegisterHealthServers(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		t.Run("registers with logger", func(t *testing.T) {
			server := grpc.NewServer()
			logger := slog.Default()

			v1Server, err := health.RegisterHealthServers(server, logger)
			require.NoError(t, err)
			assert.NotNil(t, v1Server)
		})

		t.Run("registers with nil logger", func(t *testing.T) {
			server := grpc.NewServer()

			v1Server, err := health.RegisterHealthServers(server, nil)
			require.NoError(t, err)
			assert.NotNil(t, v1Server)
		})
	})
}
