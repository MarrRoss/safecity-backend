package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewZeroLogger(level zerolog.Level) zerolog.Logger {
	//infoSampler := &zerolog.BurstSampler{
	//	Burst:       3,
	//	Period:      5 * time.Second,
	//	NextSampler: &zerolog.BasicSampler{N: 5},
	//}
	//warnSampler := &zerolog.BurstSampler{
	//	Burst:  3,
	//	Period: 1 * time.Second,
	//	// Log every 5th message after exceeding the burst rate of 3 messages per
	//	// second
	//	NextSampler: &zerolog.BasicSampler{N: 5},
	//}

	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
		// os.Stdout,
	).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
	//Sample(zerolog.LevelSampler{
	//	WarnSampler: warnSampler,
	//	InfoSampler: infoSampler,
	//})
	return logger
}
