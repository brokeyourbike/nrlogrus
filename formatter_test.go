package nrlogrus

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testAppName   = "app-name"
	testTime      = time.Date(2014, time.November, 28, 1, 1, 0, 0, time.UTC)
	matchAnything = struct{}{}
)

func newTestLogger(out io.Writer) *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(NewFormatter(testAppName, &logrus.JSONFormatter{}))
	l.SetReportCaller(true)
	l.SetOutput(out)
	return l
}

func validateOutput(t *testing.T, out *bytes.Buffer, expected map[string]interface{}) {
	var actual map[string]interface{}

	err := json.Unmarshal(out.Bytes(), &actual)
	require.NoError(t, err)

	for k, v := range expected {
		found, ok := actual[k]
		assert.Truef(t, ok, "key %s not found:\nactual=%s", k, actual)

		if v != matchAnything && found != v {
			t.Errorf("value for key %s is incorrect:\nactual=%s\nexpected=%s", k, found, v)
		}
	}
	for k, v := range actual {
		_, ok := expected[k]
		assert.Truef(t, ok, "unexpected key found:\nkey=%s\nvalue=%s", k, v)
	}
}

func BenchmarkRawJSONFormatter(b *testing.B) {
	log := newTestLogger(bytes.NewBuffer([]byte("")))
	log.Formatter = new(logrus.JSONFormatter)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.WithContext(ctx).Info("Hello World!")
	}
}

func BenchmarkRawTextFormatter(b *testing.B) {
	log := newTestLogger(bytes.NewBuffer([]byte("")))
	log.Formatter = new(logrus.TextFormatter)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.WithContext(ctx).Info("Hello World!")
	}
}

func BenchmarkWithoutTransaction(b *testing.B) {
	log := newTestLogger(bytes.NewBuffer([]byte("")))
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
	log := newTestLogger(bytes.NewBuffer([]byte("")))
	ctx := newrelic.NewContext(context.Background(), txn)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.WithContext(ctx).Info("Hello World!")
	}
}
