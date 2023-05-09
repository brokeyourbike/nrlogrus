package nrlogrus

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/newrelic/go-agent/v3/integrations/logcontext"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

const appName = "my app"

func newBenchmarkLogger(out io.Writer) *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(NewFormatter(appName, &logrus.JSONFormatter{}))
	l.SetReportCaller(true)
	l.SetOutput(out)
	return l
}

func newTestLogger(out io.Writer) (*logrus.Logger, *test.Hook) {
	l, hook := test.NewNullLogger()
	l.SetFormatter(NewFormatter(appName, &logrus.JSONFormatter{}))
	l.SetReportCaller(true)
	l.SetOutput(out)
	return l, hook
}

func TestNewFormatter(t *testing.T) {
	formatter := NewFormatter(appName, &logrus.JSONFormatter{})
	assert.Equal(t, appName, formatter.appName)
}

func TestNewFormatterFromEnvironment(t *testing.T) {
	os.Setenv("NEW_RELIC_APP_NAME", "app-name-from-env")

	formatter := NewFormatterFromEnvironment(&logrus.JSONFormatter{})
	assert.Equal(t, "app-name-from-env", formatter.appName)
}

func TestFormat(t *testing.T) {
	out := bytes.NewBuffer([]byte{})
	log, hook := newTestLogger(out)

	log.Info("Hello World!")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, appName, hook.LastEntry().Data[logcontext.KeyEntityName])
	assert.Equal(t, "SERVICE", hook.LastEntry().Data[logcontext.KeyEntityType])
}

func BenchmarkRawJSONFormatter(b *testing.B) {
	log := newBenchmarkLogger(bytes.NewBuffer([]byte("")))
	log.Formatter = new(logrus.JSONFormatter)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.WithContext(ctx).Info("Hello World!")
	}
}

func BenchmarkRawTextFormatter(b *testing.B) {
	log := newBenchmarkLogger(bytes.NewBuffer([]byte("")))
	log.Formatter = new(logrus.TextFormatter)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.WithContext(ctx).Info("Hello World!")
	}
}

func BenchmarkWithoutTransaction(b *testing.B) {
	log := newBenchmarkLogger(bytes.NewBuffer([]byte("")))
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.WithContext(ctx).Info("Hello World!")
	}
}

func BenchmarkWithTransaction(b *testing.B) {
	app := &newrelic.Application{}
	txn := app.StartTransaction("TestLogDistributedTracingDisabled")
	log := newBenchmarkLogger(bytes.NewBuffer([]byte("")))
	ctx := newrelic.NewContext(context.Background(), txn)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.WithContext(ctx).Info("Hello World!")
	}
}
