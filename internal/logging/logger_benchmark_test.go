package logging_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

// Test data
var (
	simpleString = "test message"
	typicalID    = "01234567-89ab-cdef-0123-456789abcdef"
	longString   = "This is a much longer message that might contain some detailed information about what happened in the system"
	duration     = 123456 * time.Microsecond
	timestamp    = time.Date(2024, 12, 22, 15, 0o4, 0o5, 0, time.UTC)
)

type complexStruct struct {
	ID        string
	Count     int
	Valid     bool
	CreatedAt time.Time
}

var complexData = complexStruct{
	ID:        typicalID,
	Count:     42,
	Valid:     true,
	CreatedAt: timestamp,
}

func setupSlog() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}

func setupSlogZerolog() *slog.Logger {
	zl := zerolog.New(io.Discard)
	opts := slogzerolog.Option{Level: slog.LevelInfo, Logger: &zl}
	return slog.New(opts.NewZerologHandler())
}

func setupZerolog() zerolog.Logger {
	return zerolog.New(io.Discard)
}

// Simple log line benchmarks
func BenchmarkSimpleLog_SlogInefficient(b *testing.B) {
	logger := setupSlog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(simpleString, "key", "value")
	}
}

func BenchmarkSimpleLog_SlogZerologInefficient(b *testing.B) {
	logger := setupSlogZerolog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(simpleString, "key", "value")
	}
}

func BenchmarkSimpleLog_SlogEfficient(b *testing.B) {
	logger := setupSlog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogAttrs(nil, slog.LevelInfo, simpleString,
			slog.String("key", "value"))
	}
}

func BenchmarkSimpleLog_SlogZerologEfficient(b *testing.B) {
	logger := setupSlogZerolog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogAttrs(nil, slog.LevelInfo, simpleString,
			slog.String("key", "value"))
	}
}

func BenchmarkSimpleLog_Zerolog(b *testing.B) {
	logger := setupZerolog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("key", "value").Msg(simpleString)
	}
}

// Typical log line benchmarks
func BenchmarkTypicalLog_SlogInefficient(b *testing.B) {
	logger := setupSlog()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "operation completed",
			"request_id", typicalID,
			"duration", duration,
			"status", 200)
	}
}

func BenchmarkTypicalLog_SlogZerologInefficient(b *testing.B) {
	logger := setupSlogZerolog()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "operation completed",
			"request_id", typicalID,
			"duration", duration,
			"status", 200)
	}
}

func BenchmarkTypicalLog_SlogEfficient(b *testing.B) {
	logger := setupSlog()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogAttrs(ctx, slog.LevelInfo, "operation completed",
			slog.String("request_id", typicalID),
			slog.Duration("duration", duration),
			slog.Int("status", 200))
	}
}

func BenchmarkTypicalLog_SlogZerologEfficient(b *testing.B) {
	logger := setupSlogZerolog()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogAttrs(ctx, slog.LevelInfo, "operation completed",
			slog.String("request_id", typicalID),
			slog.Duration("duration", duration),
			slog.Int("status", 200))
	}
}

func BenchmarkTypicalLog_Zerolog(b *testing.B) {
	logger := setupZerolog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().
			Str("request_id", typicalID).
			Dur("duration", duration).
			Int("status", 200).
			Msg("operation completed")
	}
}

// Complex log line benchmarks
func BenchmarkComplexLog_SlogInefficient(b *testing.B) {
	logger := setupSlog()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, longString,
			"request_id", typicalID,
			"duration", duration,
			"data", complexData,
			"timestamp", timestamp,
			"retry_count", 3,
			"errors", []string{"error1", "error2"},
			"metadata", map[string]string{
				"region": "us-west-2",
				"zone":   "a",
			})
	}
}

func BenchmarkComplexLog_SlogZerologInefficient(b *testing.B) {
	logger := setupSlogZerolog()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, longString,
			"request_id", typicalID,
			"duration", duration,
			"data", complexData,
			"timestamp", timestamp,
			"retry_count", 3,
			"errors", []string{"error1", "error2"},
			"metadata", map[string]string{
				"region": "us-west-2",
				"zone":   "a",
			})
	}
}

func BenchmarkComplexLog_SlogEfficient(b *testing.B) {
	logger := setupSlog()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogAttrs(ctx, slog.LevelInfo, longString,
			slog.String("request_id", typicalID),
			slog.Duration("duration", duration),
			slog.Any("data", complexData),
			slog.Time("timestamp", timestamp),
			slog.Int("retry_count", 3),
			slog.Any("errors", []string{"error1", "error2"}),
			slog.Any("metadata", map[string]string{
				"region": "us-west-2",
				"zone":   "a",
			}))
	}
}

func BenchmarkComplexLog_SlogZerologEfficient(b *testing.B) {
	logger := setupSlogZerolog()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogAttrs(ctx, slog.LevelInfo, longString,
			slog.String("request_id", typicalID),
			slog.Duration("duration", duration),
			slog.Any("data", complexData),
			slog.Time("timestamp", timestamp),
			slog.Int("retry_count", 3),
			slog.Any("errors", []string{"error1", "error2"}),
			slog.Any("metadata", map[string]string{
				"region": "us-west-2",
				"zone":   "a",
			}))
	}
}

func BenchmarkComplexLog_Zerolog(b *testing.B) {
	logger := setupZerolog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().
			Str("request_id", typicalID).
			Dur("duration", duration).
			Interface("data", complexData).
			Time("timestamp", timestamp).
			Int("retry_count", 3).
			Strs("errors", []string{"error1", "error2"}).
			Str("metadata.region", "us-west-2").
			Str("metadata.zone", "a").
			Msg(longString)
	}
}
