package nrlogrus

import (
	"fmt"
	"os"

	"github.com/newrelic/go-agent/v3/integrations/logcontext"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

// ContextFormatter is a `logrus.Formatter` that will format logs for sending
// to New Relic.
type ContextFormatter struct {
	appName   string
	formatter logrus.Formatter
}

// NewFormatter creates a new `logrus.Formatter` that will format logs for
// sending to New Relic.
func NewFormatter(appName string, formatter logrus.Formatter) ContextFormatter {
	return ContextFormatter{
		appName:   appName,
		formatter: formatter,
	}
}

// NewFormatterFromEnvironment creates a new `logrus.Formatter` that will
// format logs for sending to New Relic.  The application name is read from
// the `NEW_RELIC_APP_NAME` environment variable.
func NewFormatterFromEnvironment(formatter logrus.Formatter) ContextFormatter {
	return NewFormatter(os.Getenv("NEW_RELIC_APP_NAME"), formatter)
}

// Format renders a single log entry.
func (f ContextFormatter) Format(e *logrus.Entry) ([]byte, error) {
	md := newrelic.LinkingMetadata{EntityName: f.appName, EntityType: "SERVICE"}

	if ctx := e.Context; nil != ctx {
		if trx := newrelic.FromContext(ctx); nil != trx {
			md = trx.GetLinkingMetadata()
		}
	}

	logcontext.AddLinkingMetadata(e.Data, md)

	e.Data[logcontext.KeyTimestamp] = uint64(e.Time.UnixNano()) / uint64(1000*1000)
	e.Data[logcontext.KeyMessage] = e.Message
	e.Data[logcontext.KeyLevel] = e.Level

	if e.HasCaller() {
		e.Data[logcontext.KeyFile] = e.Caller.File
		e.Data[logcontext.KeyLine] = e.Caller.Line
		e.Data[logcontext.KeyMethod] = e.Caller.Function
	}

	b, err := f.formatter.Format(e)
	if err != nil {
		return nil, fmt.Errorf("error formatting log entry: %w", err)
	}

	return b, nil
}
