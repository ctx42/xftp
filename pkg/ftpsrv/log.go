package ftpsrv

import (
	"context"

	"github.com/rs/zerolog"
)

// LogWarn logs error as a warning. Properly handles nil and joined errors.
// Both context and the meta arguments may be set to nil.
func LogWarn(ctx context.Context, log zerolog.Logger, err error, meta any) {
	recError(ctx, log, zerolog.WarnLevel, err, meta)
}

// LogError logs error. Properly handles nil and joined errors. Both context
// and the meta arguments may be set to nil.
func LogError(ctx context.Context, log zerolog.Logger, err error, meta any) {
	recError(ctx, log, zerolog.ErrorLevel, err, meta)
}

func recError(
	ctx context.Context,
	log zerolog.Logger,
	level zerolog.Level,
	err error,
	meta any,
) {

	if err == nil {
		return
	}

	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		for i, e := range joined.Unwrap() {
			em := []any{"error_index", i}
			log.WithLevel(level).
				Ctx(ctx).
				Fields(meta).
				Fields(em).
				Msg(e.Error())
		}
		return
	}

	evt := log.
		WithLevel(level).
		Ctx(ctx).
		Fields(meta)
	evt.Msg(err.Error())
}
