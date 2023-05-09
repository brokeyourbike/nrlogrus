# nrlogrus

[![Go Reference](https://pkg.go.dev/badge/github.com/brokeyourbike/nrlogrus.svg)](https://pkg.go.dev/github.com/brokeyourbike/nrlogrus)
[![Go Report Card](https://goreportcard.com/badge/github.com/brokeyourbike/nrlogrus)](https://goreportcard.com/report/github.com/brokeyourbike/nrlogrus)
[![Maintainability](https://api.codeclimate.com/v1/badges/215df4233533b9971565/maintainability)](https://codeclimate.com/github/brokeyourbike/nrlogrus/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/215df4233533b9971565/test_coverage)](https://codeclimate.com/github/brokeyourbike/nrlogrus/test_coverage)

NewRelic hook for Logrus, with logcontext support.

## Why

Official hooks are not covering cases when NewRelic client is not yet created, or failed, but we still want logs for this entity to be connected to the APM. Turns out for this we need to only send the `entity.name` and `entity.type`.

## Install

```bash
go get github.com/brokeyourbike/nrlogrus
```

## Use

You can pass the application name to the function:

```golang
log.SetFormatter(nrlogrus.NewFormatter("my-app", &logrus.JSONFormatter{}))
```

or read it from environment variable `NEW_RELIC_APP_NAME`:

```go
log.SetFormatter(nrlogrus.NewFormatterFromEnvironment(&logrus.JSONFormatter{}))
```

## Thanks

- [logcontext-v2/nrlogrus](https://pkg.go.dev/github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus)
- [logcontext/nrlogrusplugin](https://pkg.go.dev/github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin)

## Authors
- [Ivan Stasiuk](https://github.com/brokeyourbike) | [Twitter](https://twitter.com/brokeyourbike) | [LinkedIn](https://www.linkedin.com/in/brokeyourbike) | [stasi.uk](https://stasi.uk)

## License
[Apache-2.0 License](https://github.com/glocurrency/nrlogrus/blob/main/LICENSE)