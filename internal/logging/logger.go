package logging

import (
	"os"
	"strings"

	stdlog "log"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx/fxevent"
)

type (
	ZeroLogger struct {
		Logger *zerolog.Logger
	}
)

var (
	logger = zerolog.New(os.Stdout).Level(zerolog.TraceLevel).With().Timestamp().Logger()
)

func Init() error {
	cfg := newConfig()

	if cfg.Dev {
		logger = zerolog.New(zerolog.NewConsoleWriter()).Level(zerolog.TraceLevel).With().Timestamp().Logger()
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}

	if level < zerolog.NoLevel {
		logger.Info().Msgf("Setting log level to: %s", level.String())
		zerolog.SetGlobalLevel(level)
		logger.Level(level)
	}

	return nil
}

func newZeroLogger() *zerolog.Logger {
	log.Logger = logger

	return &logger
}

func newFxLogger(logger *zerolog.Logger) fxevent.Logger {
	return &ZeroLogger{logger}
}

// Assert that our `ZeroLogger` implements `fxevent.Logger`.
var _ fxevent.Logger = (*ZeroLogger)(nil)

//nolint:gocyclo,gocognit
func (l *ZeroLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Logger.
			Info().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Logger.
				Err(e.Err).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStart hook failed")
		} else {
			l.Logger.
				Info().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		l.Logger.
			Info().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Logger.
				Err(e.Err).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStop hook failed")
		} else {
			l.Logger.
				Info().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStop hook executed")
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.Logger.
				Err(e.Err).
				Str("type", e.TypeName).
				Str("module", e.ModuleName).
				Msg("error encountered while applying options")
		} else {
			l.Logger.
				Info().
				Str("type", e.TypeName).
				Str("module", e.ModuleName).
				Msg("supplied")
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.
				Info().
				Str("constructor", e.ConstructorName).
				Str("module", e.ModuleName).
				Str("type", rtype).
				Bool("private", e.Private).
				Msg("provided")
		}
		if e.Err != nil {
			l.Logger.
				Err(e.Err).
				Str("module", e.ModuleName).
				Msg("error encountered while applying options")
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.
				Info().
				Str("module", e.ModuleName).
				Str("type", rtype).
				Msg("replaced")
		}
		if e.Err != nil {
			l.Logger.
				Err(e.Err).
				Str("module", e.ModuleName).
				Msg("error encountered while replacing")
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.
				Info().
				Str("decorator", e.DecoratorName).
				Str("module", e.ModuleName).
				Str("type", rtype).
				Msg("decorated")
		}
		if e.Err != nil {
			l.Logger.
				Err(e.Err).
				Str("module", e.ModuleName).
				Msg("error encountered while applying options")
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		l.Logger.
			Info().
			Str("function", e.FunctionName).
			Str("module", e.ModuleName).
			Msg("invoking")
	case *fxevent.Invoked:
		if e.Err != nil {
			l.Logger.
				Err(e.Err).
				Str("stack", e.Trace).
				Str("function", e.FunctionName).
				Str("module", e.ModuleName).
				Msg("invoke failed")
		}
	case *fxevent.Stopping:
		l.Logger.
			Info().
			Str("signal", strings.ToUpper(e.Signal.String())).
			Msg("received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.
				Err(e.Err).
				Msg("stop failed")
		}
	case *fxevent.RollingBack:
		l.Logger.
			Err(e.StartErr).
			Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Logger.
				Err(e.Err).
				Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.
				Err(e.Err).
				Msg("start failed")
		} else {
			l.Logger.
				Info().
				Msg("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Logger.
				Err(e.Err).
				Msg("custom logger initialization failed")
		} else {
			l.Logger.
				Info().
				Str("function", e.ConstructorName).
				Msg("initialized custom fxevent.Logger")
		}
	}
}
