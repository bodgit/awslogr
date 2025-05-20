// Package awslogr provides a means to pass a [logr.Logger] to the AWS SDK as
// a [logging.Logger].
//
// The optional [logging.ContextLogger] interface is implemented which is
// especially useful if using the [otellogr] bridge as this will then include
// trace & span ID attributes in the logs when using the [otelaws] middleware.
//
// An OpenTelemetry example:
//
//	l := logr.New(otellogr.NewLogSink("my/pkg/name"))
//
//	nl, err := awslogr.New(l)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	cfg, err := config.LoadDefaultConfig(context.Background(),
//		config.WithClientLogMode(aws.LogRetries|aws.logRequest),
//		config.WithLogger(nl))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	otelaws.AppendMiddlewares(&cfg.APIOptions)
//
//	// Use cfg to now create an AWS service client
//
// [otellogr]: https://go.opentelemetry.io/contrib/bridges/otellogr
// [otelaws]: http://go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws
package awslogr

import (
	"context"
	"fmt"

	"github.com/aws/smithy-go/logging"
	"github.com/go-logr/logr"
)

type logger struct {
	l logr.Logger
	f func(logging.Classification) int
	c context.Context //nolint:containedctx
	k string
}

func (l *logger) Logf(classification logging.Classification, format string, v ...interface{}) {
	nl := l.l.V(l.f(classification))

	if !nl.Enabled() {
		return
	}

	if l.c != nil && l.k != "" {
		nl.Info(fmt.Sprintf(format, v...), l.k, l.c)
	} else {
		nl.Info(fmt.Sprintf(format, v...))
	}
}

func (l *logger) WithContext(ctx context.Context) logging.Logger {
	return &logger{
		l: l.l,
		f: l.f,
		c: ctx,
		k: l.k,
	}
}

var (
	_ logging.Logger        = new(logger)
	_ logging.ContextLogger = new(logger)
)

// WithContextKey sets the attribute key used for logging a [context.Context]
// attribute. Using the empty string will suppress the [context.Context] from
// being included.
func WithContextKey(key string) func(*logger) error {
	return func(l *logger) error {
		l.k = key

		return nil
	}
}

// WithClassificationLevel allows setting a custom mapping function from
// [logging.Classification] to logr verbosity levels.
func WithClassificationLevel(f func(logging.Classification) int) func(*logger) error {
	return func(l *logger) error {
		l.f = f

		return nil
	}
}

// New returns a new [logging.Logger] that by default logs a [context.Context]
// attribute with the key "context". Verbosity level 1 is used for any
// [logging.Debug] messages, otherwise level 0 is used.
func New(l logr.Logger, options ...func(*logger) error) (logging.Logger, error) {
	nl := &logger{
		l: l,
		f: func(classification logging.Classification) int {
			if classification == logging.Debug {
				return 1
			}

			return 0
		},
		k: "context",
	}

	for _, option := range options {
		if err := option(nl); err != nil {
			return nil, err
		}
	}

	return nl, nil
}
