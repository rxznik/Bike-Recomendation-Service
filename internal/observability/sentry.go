package observability

import (
	"errors"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

const (
	devEnvType  = "dev"
	prodEnvType = "prod"
)

type SentryScopeOptions struct {
	Level   *sentry.Level
	Tags    map[string]string
	Extra   map[string]any
	Context *sentry.Context
	Request *http.Request
}

func InitGlobalSentry(dsn string, envType string, serviceName string) error {
	var isDebug bool
	var tracesSampleRate float64

	switch envType {
	case devEnvType:
		isDebug = true
		tracesSampleRate = 1
	case prodEnvType:
		isDebug = false
		tracesSampleRate = 0.33
	default:
		msg := "unknown env type"
		zap.L().Fatal(msg, zap.String("env", envType))
		return errors.New(msg)
	}

	return sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            isDebug,
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: tracesSampleRate,
		ServerName:       serviceName,
	})
}

func FlushAndRecoverSentry(timeout ...time.Duration) {
	var t time.Duration
	if len(timeout) == 0 {
		t = 5 * time.Second
	} else {
		t = timeout[0]
	}
	sentry.Flush(t)
	sentry.Recover()
}

func CaptureMessageSentry(message string) {
	sentry.CaptureMessage(message)
}

func CaptureExceptionSentry(err error) {
	sentry.CaptureException(err)
}

func CaptureEventSentry(event *sentry.Event) {
	sentry.CaptureEvent(event)
}

func NewLocalHubSentry() *sentry.Hub {
	return sentry.CurrentHub().Clone()
}

func ConfigureScopeSentry(opts *SentryScopeOptions) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		if opts == nil {
			return
		}

		if opts.Context != nil {
			scope.SetContext("context", *opts.Context)
		}

		if opts.Level != nil {
			scope.SetLevel(*opts.Level)
		}

		if opts.Tags != nil {
			scope.SetTags(opts.Tags)
		}

		if opts.Extra != nil {
			scope.SetExtras(opts.Extra)
		}

		if opts.Request != nil {
			scope.SetRequest(opts.Request)
		}
	})
}
