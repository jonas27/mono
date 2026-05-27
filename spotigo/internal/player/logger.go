package player

import (
	"fmt"
	"log/slog"

	librespot "github.com/devgianlu/go-librespot"
)

type cliLogger struct {
	log *slog.Logger
}

func newLogger() librespot.Logger {
	return &cliLogger{log: slog.Default()}
}

func (l *cliLogger) Tracef(string, ...any)                     {}
func (l *cliLogger) Debugf(string, ...any)                     {}
func (l *cliLogger) Infof(f string, a ...any)                  { l.log.Info(fmt.Sprintf(f, a...)) }
func (l *cliLogger) Warnf(f string, a ...any)                  { l.log.Warn(fmt.Sprintf(f, a...)) }
func (l *cliLogger) Errorf(f string, a ...any)                 { l.log.Error(fmt.Sprintf(f, a...)) }
func (l *cliLogger) Trace(...any)                              {}
func (l *cliLogger) Debug(...any)                              {}
func (l *cliLogger) Info(a ...any)                             { l.log.Info(fmt.Sprint(a...)) }
func (l *cliLogger) Warn(a ...any)                             { l.log.Warn(fmt.Sprint(a...)) }
func (l *cliLogger) Error(a ...any)                            { l.log.Error(fmt.Sprint(a...)) }
func (l *cliLogger) WithField(k string, v any) librespot.Logger { return &cliLogger{log: l.log.With(k, v)} }
func (l *cliLogger) WithError(err error) librespot.Logger      { return &cliLogger{log: l.log.With("error", err)} }
