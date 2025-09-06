package observability

import (
	//"github.com/alexdrl/zerowater"
	"github.com/rs/zerolog"
)

type Observability struct {
	Logger zerolog.Logger
}

func New(logger zerolog.Logger) *Observability {
	return &Observability{
		Logger: logger,
	}
}

func (o Observability) GetLogger() *zerolog.Logger {
	return &o.Logger
}

//func (o Observability) WatermillAdapter() *zerowater.ZerologLoggerAdapter {
//	return zerowater.NewZerologLoggerAdapter(
//		o.Logger.With().Str("component", "windmill").Logger(),
//	)
//}
