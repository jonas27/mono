package player

import (
	"fmt"
	"os"

	librespot "github.com/devgianlu/go-librespot"
)

// cliLogger implements librespot.Logger, surfacing warnings and errors only.
type cliLogger struct{}

func newLogger() librespot.Logger { return cliLogger{} }

func (cliLogger) Tracef(string, ...any)                    {}
func (cliLogger) Debugf(string, ...any)                    {}
func (cliLogger) Infof(f string, a ...any)                 { fmt.Fprintf(os.Stderr, "info: "+f+"\n", a...) }
func (cliLogger) Warnf(f string, a ...any)                 { fmt.Fprintf(os.Stderr, "warn: "+f+"\n", a...) }
func (cliLogger) Errorf(f string, a ...any)                { fmt.Fprintf(os.Stderr, "err:  "+f+"\n", a...) }
func (cliLogger) Trace(...any)                             {}
func (cliLogger) Debug(...any)                             {}
func (cliLogger) Info(a ...any)                            { fmt.Fprintln(os.Stderr, "info:", fmt.Sprint(a...)) }
func (cliLogger) Warn(a ...any)                            { fmt.Fprintln(os.Stderr, "warn:", fmt.Sprint(a...)) }
func (cliLogger) Error(a ...any)                           { fmt.Fprintln(os.Stderr, "err: ", fmt.Sprint(a...)) }
func (l cliLogger) WithField(string, any) librespot.Logger { return l }
func (l cliLogger) WithError(error) librespot.Logger       { return l }
