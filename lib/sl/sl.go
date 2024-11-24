package sl

import (
	"fmt"
	"log/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Infof(log *slog.Logger, format string, args ...any) {
	log.Info(fmt.Sprintf(format, args))
}
